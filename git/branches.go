package git

import (
	"errors"
	"strconv"
	"strings"
)

type Branch struct {
	Name string
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
func (g *Git) GetBranches() (map[string]Branch, error) {
	branch_list := make(map[string]Branch)

	branches, err := g.RunCwd("branch", "--list", "--format", "'%(refname:short)'")
	if err != nil {
		println(err.Error())
		return branch_list, err
	}

	bs := parseLines(branches)
	for _, branch := range bs {
		if strings.Trim(branch, " ") != "" {
			branch_list[branch] = Branch{
				Name: branch,
			}
		}
	}

	return branch_list, nil
}

// Switch branch on currently selected repo.
func (g *Git) SwitchBranch(branch string) error {
	_, err := g.RunCwd("checkout", branch)
	return err
}

// If branch exists locally.
func (g *Git) BranchExists(name string) bool {
	_, err := g.RunCwd("rev-parse", "--verify", name)
	return err == nil
}

// Get commits ahead and behind branch is from specific remote.
func (g *Git) getAheadBehind(branch string, remote string) (uint32, uint32, error) {
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

func (g *Git) PushBranch(remote string, branch string) error {
	_, err := g.RunCwd("push", remote, branch+":"+branch)
	return err
}

func (g *Git) PullBranch(remote string, branch string, rebase bool) error {
	var err error
	if rebase {
		_, err = g.RunCwd("pull", remote, branch+":"+branch, "--rebase")
	} else {
		_, err = g.RunCwd("pull", remote, branch+":"+branch)
	}
	return err
}

func (g *Git) getNumCommitsOnBranch(branch string) (uint32, error) {
	n, err := g.RunCwd("rev-list", "--count", branch)
	if err != nil {
		return 0, err
	}
	num, _ := strconv.ParseInt(n, 10, 32)
	return uint32(num), nil
}
