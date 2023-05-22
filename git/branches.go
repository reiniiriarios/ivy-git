package git

import (
	"errors"
	"strconv"
	"strings"
)

type Branch struct {
	Name     string
	Upstream string
}

func (g *Git) NumBranches() uint64 {
	b, err := g.RunCwd("branch")
	if err != nil {
		println(err.Error())
		return 0
	}
	lines := parseLines(b)
	return uint64(len(lines))
}

// Check common names for main branch.
func (g *Git) NameOfMainBranch() string {
	r, err := g.RunCwd("for-each-ref", "--format=%(refname:short)", "refs/heads/main", "refs/heads/master", "refs/heads/trunk")
	if err != nil {
		// Screw it, return something.
		return "main"
	}
	r = parseOneLine(r)
	if !strings.Contains(r, "\n") {
		return r
	}
	// More than one result.
	if strings.Contains(r, "master") {
		return "master"
	}
	// Default to main.
	return "main"
}

// Get current branch for currently selected repo.
func (g *Git) GetCurrentBranch() (string, error) {
	branch, err := g.RunCwd("rev-parse", "--abbrev-ref", "HEAD")
	branch = strings.ReplaceAll(strings.ReplaceAll(branch, "\r", ""), "\n", "")
	if err != nil {
		println(err.Error())
		return "", err
	}

	return branch, nil
}

// Get list of all branches for currently selected repo.
func (g *Git) GetBranches() ([]Branch, error) {
	branch_list := []Branch{}

	branches, err := g.RunCwd("branch", "--list", "--format", "%(refname:short)"+GIT_LOG_SEP+"%(upstream:short)")
	if err != nil {
		println(err.Error())
		return branch_list, err
	}

	bs := parseLines(branches)
	for _, branch := range bs {
		parts := strings.Split(branch, GIT_LOG_SEP)
		if len(parts) == 2 {
			branch_list = append(branch_list, Branch{
				Name:     parts[0],
				Upstream: parts[1],
			})
		}
	}

	return branch_list, nil
}

func (g *Git) GetBranchUpstream(branch string) (string, error) {
	if branch == "" {
		return "", errors.New("no branch name specified")
	}

	b, err := g.RunCwd("branch", "--format", "%(upstream:short)", "--list", branch)
	if err != nil {
		println(err.Error())
		return "", err
	}
	b = parseOneLine(b)
	return b, nil
}

// Switch branch on currently selected repo.
func (g *Git) SwitchBranch(branch string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}

	_, err := g.RunCwd("checkout", branch)
	return err
}

// If branch exists locally.
func (g *Git) BranchExists(name string) bool {
	if name == "" {
		return false
	}
	_, err := g.RunCwd("rev-parse", "--verify", name)
	return err == nil
}

// Get commits ahead and behind branch is from specific remote.
func (g *Git) getAheadBehind(branch string, remote string) (uint32, uint32, error) {
	if branch == "" {
		return 0, 0, errors.New("no branch name specified")
	}
	if remote == "" {
		return 0, 0, errors.New("no remote name specified")
	}

	rl, err := g.RunCwd("rev-list", "--left-right", "--count", branch+"..."+remote+"/"+branch)
	if err != nil {
		return 0, 0, err
	}
	ab := strings.Fields(rl)
	if len(ab) != 2 {
		return 0, 0, errors.New("error parsing rev-list --left-right")
	}
	ahead, _ := strconv.ParseInt(ab[0], 10, 32)
	behind, _ := strconv.ParseInt(ab[1], 10, 32)
	return uint32(ahead), uint32(behind), nil
}

func (g *Git) PushBranch(branch string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}

	remote, err := g.getBranchRemote(branch)
	if err != nil {
		return err
	}
	err = g.PushRemoteBranch(remote, branch)
	return err
}

func (g *Git) PushRemoteBranch(remote string, branch string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	if remote == "" {
		return errors.New("no remote name specified")
	}

	_, err := g.RunCwd("push", remote, branch+":"+branch)
	return err
}

func (g *Git) PullBranch(branch string, rebase bool) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}

	remote, err := g.getBranchRemote(branch)
	if err != nil {
		return err
	}
	err = g.PullRemoteBranch(remote, branch, rebase)
	return err
}

func (g *Git) PullRemoteBranch(remote string, branch string, rebase bool) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	if remote == "" {
		return errors.New("no remote name specified")
	}

	var err error
	if rebase {
		_, err = g.RunCwd("pull", remote, branch+":"+branch, "--rebase")
	} else {
		_, err = g.RunCwd("pull", remote, branch+":"+branch)
	}
	return err
}

func (g *Git) NumMainBranchCommits() uint64 {
	main := g.NameOfMainBranch()
	num, err := g.NumCommitsOnBranch(main)
	if err != nil {
		println(err.Error())
		return 0
	}
	return uint64(num)
}

func (g *Git) NumCommitsOnBranch(branch string) (uint64, error) {
	if branch == "" {
		return 0, errors.New("no branch name specified")
	}

	n, err := g.RunCwd("rev-list", "--count", branch)
	if err != nil {
		println(err.Error())
		return 0, err
	}
	n = parseOneLine(n)
	num, _ := strconv.ParseInt(n, 10, 32)
	return uint64(num), nil
}

func (g *Git) getBranchRemote(branch string) (string, error) {
	if branch == "" {
		return "", errors.New("no branch name specified")
	}

	r, err := g.RunCwd("config", "branch."+branch+".remote")
	if err != nil {
		return "", err
	}
	r = parseOneLine(r)
	return r, nil
}

func (g *Git) fetchBranchRemote(branch string, remote string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	if remote == "" {
		return errors.New("no remote name specified")
	}

	_, err := g.RunCwd("fetch", remote, branch)
	return err
}

func (g *Git) ResetBranchToRemote(branch string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}

	remote, err := g.getBranchRemote(branch)
	if err != nil {
		return err
	}
	err = g.fetchBranchRemote(branch, remote)
	if err != nil {
		return err
	}

	current_branch, err := g.GetCurrentBranch()
	if err != nil {
		return err
	}
	err = g.SwitchBranch(branch)
	if err != nil {
		return err
	}

	_, err = g.RunCwd("reset", "--hard", remote+"/"+branch)

	g.SwitchBranch(current_branch)

	return err
}

// Delete a branch.
func (g *Git) DeleteBranch(branch string, force bool, delete_on_remotes bool) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}

	delete := "-d"
	if force {
		delete = "-D"
	}

	_, err := g.RunCwd("branch", delete, branch)
	if err != nil {
		return err
	}

	if delete_on_remotes {
		remotes, err := g.getRemoteNames()
		if err != nil {
			return err
		}
		for _, remote := range remotes {
			if g.branchExistsOnRemote(branch, remote) {
				_, err := g.RunCwd("push", delete, remote, branch)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Delete a remote branch.
func (g *Git) DeleteRemoteBranch(branch string, remote string, force bool) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	if remote == "" {
		return errors.New("no remote name specified")
	}

	delete := "-d"
	if force {
		delete = "-D"
	}
	_, err := g.RunCwd("push", delete, remote, branch)
	return err
}

// If a branch exists on a specific remote.
func (g *Git) branchExistsOnRemote(branch string, remote string) bool {
	if branch == "" || remote == "" {
		return false
	}

	ls, _ := g.RunCwd("ls-remote", "--heads", remote, branch)
	ls = parseOneLine(ls)
	return ls != ""
}

// Rename a branch locally and on all remotes.
func (g *Git) RenameBranch(branch string, new_name string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	if new_name == "" {
		return errors.New("no new branch name specified")
	}

	_, err := g.RunCwd("branch", "-m", branch, new_name)
	if err != nil {
		return err
	}

	remotes, err := g.getRemoteNames()
	if err != nil {
		return err
	}
	for _, remote := range remotes {
		if g.branchExistsOnRemote(branch, remote) {
			_, err := g.RunCwd("push", remote, ":"+branch, new_name)
			if err != nil {
				return err
			}
			_, err = g.RunCwd("push", "--set-upstream", remote, new_name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Git) RebaseOnBranch(branch string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	_, err := g.RunCwd("rebase", branch)
	return err
}

func (g *Git) CreateBranch(name string, at_hash string, checkout bool) error {
	if name == "" {
		return errors.New("no branch name specified")
	}
	if at_hash == "" {
		return errors.New("no commit hash specified")
	}
	var err error
	if checkout {
		_, err = g.RunCwd("checkout", "-b", name, at_hash)
	} else {
		_, err = g.RunCwd("branch", name, at_hash)
	}
	return err
}
