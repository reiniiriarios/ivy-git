package git

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Branch struct {
	Name     string
	Upstream string
	Remote   string
}

func (g *Git) NumBranches() uint64 {
	b, err := g.run("branch")
	if err != nil {
		return 0
	}
	lines := parseLines(b)
	return uint64(len(lines))
}

// Get current branch for currently selected repo.
func (g *Git) GetCurrentBranch() (string, error) {
	g.GetRepoState()
	if g.Repo.State == RepoStateRebase || g.Repo.State == RepoStateRebaseMerge || g.Repo.State == RepoStateRebaseInteractive {
		return g.getCurrentBranchDuringRebaseMerge()
	}
	return g.getCurrentBranchFromSymbolicRef()
}

func (g *Git) getCurrentBranchDuringRebaseMerge() (string, error) {
	// GitDirRebaseApply
	branch := ""
	dirs := []string{GitDirRebaseMerge, GitDirRebaseApply}
	var err error = nil
	for _, dir := range dirs {
		branch, err = g.getCurrentBranchDuringRebase(dir)
		if branch != "" {
			break
		}
	}
	return branch, err
}

func (g *Git) getCurrentBranchDuringRebase(dir string) (string, error) {
	path := filepath.Join(g.Repo.Directory, ".git", dir, GitFileRebaseMergeHeadName)
	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	fi, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fi.Close()
	scanner := bufio.NewScanner(fi)
	if scanner.Scan() {
		ref := scanner.Text()
		if strings.HasPrefix(ref, "refs/heads/") {
			return ref[11:], nil
		}
	}
	return "", err
}

func (g *Git) getCurrentBranchFromSymbolicRef() (string, error) {
	// Git 2.22+  `git branch --show-current`
	// Git 1.8+   `git symbolic-ref --short HEAD`
	// Git 1.7+   `git rev-parse --abbrev-ref HEAD` -- fails when no commits
	// earlier    `git symbolic-ref HEAD` -- works, but returns full ref
	ref, err := g.run("symbolic-ref", "HEAD")
	if err != nil {
		// If the HEAD is detached, this error will proc.
		if errorCode(err) == RefNotSymbolic {
			ref, err = g.run("rev-parse", "--abbrev-ref", "HEAD")
		}
		// If any of these errors, assume there's no branch selected and no commits available to parse.
		if errorCode(err) == NoCommitsYet || errorCode(err) == BadRevision || errorCode(err) == UnknownRevisionOrPath || errorCode(err) == ExitStatus1 {
			return "", nil
		}
		if err != nil {
			// Otherwise something odd may be happening, display the error.
			return "", err
		}
	}
	// refs/heads/main => main
	ref = parseOneLine(ref)
	parts := strings.Split(ref, "/")
	branch := parts[len(parts)-1]

	return branch, nil
}

// Get list of all branches for currently selected repo.
func (g *Git) GetBranches() ([]Branch, error) {
	branch_list := []Branch{}

	branches, err := g.run("for-each-ref", "--format", "%(refname:lstrip=2)"+GIT_LOG_SEP+"%(upstream:short)", "refs/heads/**")
	if err != nil {
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

	// If there are no branches, we're probably in a new repo. Get the default branch.
	if len(branch_list) == 0 {
		ref, err := g.run("symbolic-ref", "HEAD")
		// Ignore errors.
		if err == nil {
			// refs/heads/main => main
			ref = parseOneLine(ref)
			if len(ref) > 1 {
				parts := strings.Split(ref, "/")
				branch_list = append(branch_list, Branch{
					Name: parts[len(parts)-1],
				})
			}
		}
	}

	return branch_list, nil
}

// Get list of all remote branches for currently selected repo.
func (g *Git) GetRemoteBranches() ([]Branch, error) {
	branch_list := []Branch{}

	branches, err := g.run("for-each-ref", "--format", "%(refname:lstrip=2)", "refs/remotes/**")
	if err != nil {
		return branch_list, err
	}

	bs := parseLines(branches)
	for _, branch := range bs {
		// origin/branch-name
		parts := strings.Split(branch, "/")
		if len(parts) == 2 {
			// don't need origin/HEAD
			if parts[1] != "HEAD" {
				branch_list = append(branch_list, Branch{
					Name:   parts[1],
					Remote: parts[0],
				})
			}
		}
	}

	return branch_list, nil
}

func (g *Git) GetBranchUpstream(branch string) (string, error) {
	if branch == "" {
		return "", errors.New("no branch name specified")
	}

	b, err := g.run("for-each-ref", "--format", "%(upstream:short)", branch)
	if err != nil {
		return "", err
	}
	b = parseOneLine(b)
	return b, nil
}

// If branch exists locally.
func (g *Git) BranchExists(name string) bool {
	if name == "" {
		return false
	}
	_, err := g.run("rev-parse", "--verify", name)
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

	rl, err := g.run("rev-list", "--left-right", "--count", branch+"..."+remote+"/"+branch)
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

func (g *Git) NumMainBranchCommits() uint64 {
	main := g.NameOfMainBranch()
	num, err := g.NumCommitsOnBranch(main)
	if err != nil {
		return 0
	}
	return uint64(num)
}

func (g *Git) NumCommitsOnBranch(branch string) (uint64, error) {
	if branch == "" {
		return 0, errors.New("no branch name specified")
	}

	n, err := g.run("rev-list", "--count", branch)
	if err != nil {
		// Ignore errors here, there may not be a branch selected.
		return 0, nil
	}
	n = parseOneLine(n)
	num, _ := strconv.ParseInt(n, 10, 32)
	return uint64(num), nil
}

func (g *Git) getBranchRemote(branch string, fallback bool) (string, error) {
	if branch == "" {
		return "", errors.New("no branch name specified")
	}

	r, err := g.run("config", "branch."+branch+".remote")
	if err != nil && errorCode(err) != ExitStatus1 {
		return "", err
	}
	r = parseOneLine(r)

	if r == "" && fallback {
		r = g.getRemoteFallback()
	}

	if r == "" {
		// This is a last ditch. Ignore errors.
		r, _ = g.findRemoteBranch(branch)
		if r != "" {
			// Turn 'origin/main' into 'origin'.
			if idx := strings.IndexByte(r, '/'); idx >= 0 {
				r = r[:idx]
			}
		}
	}

	return r, nil
}

// Fallback to remote:
//
//	for current branch
//	for main branch
//	whatever is first in the list of remotes
//	origin
func (g *Git) getRemoteFallback() string {
	r, err := g.GetRemoteForCurrentBranch()
	if err == nil && r != "" {
		return r
	}

	main := g.NameOfMainBranch()
	r, err = g.getBranchRemote(main, false)
	if err == nil && r != "" {
		return r
	}

	rs, err := g.getRemoteNames()
	if err == nil && len(rs) > 0 {
		return rs[0]
	}

	return "origin"
}

func (g *Git) fetchBranchRemote(branch string, remote string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	if remote == "" {
		return errors.New("no remote name specified")
	}

	_, err := g.run("fetch", remote, branch)
	return err
}

// If a branch exists on a specific remote.
func (g *Git) branchExistsOnRemote(branch string, remote string) bool {
	if branch == "" || remote == "" {
		return false
	}

	ls, _ := g.run("ls-remote", "--heads", remote, branch)
	ls = parseOneLine(ls)
	return ls != ""
}
