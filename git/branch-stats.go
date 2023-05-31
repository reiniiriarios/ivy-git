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

// Get current branch for currently selected repo.
func (g *Git) GetCurrentBranch() (string, error) {
	branch, err := g.RunCwd("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		println(err.Error())
		return "", err
	}
	branch = strings.ReplaceAll(strings.ReplaceAll(branch, "\r", ""), "\n", "")

	return branch, nil
}

// Get list of all branches for currently selected repo.
func (g *Git) GetBranches() ([]Branch, error) {
	branch_list := []Branch{}

	branches, err := g.RunCwd("for-each-ref", "--format", "%(refname:lstrip=2)"+GIT_LOG_SEP+"%(upstream:short)", "refs/heads/**")
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

	b, err := g.RunCwd("for-each-ref", "--format", "%(upstream:short)", branch)
	if err != nil {
		println(err.Error())
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

func (g *Git) getBranchRemote(branch string, fallback bool) (string, error) {
	if branch == "" {
		return "", errors.New("no branch name specified")
	}

	r, err := g.RunCwd("config", "branch."+branch+".remote")
	if err != nil {
		return "", err
	}
	r = parseOneLine(r)

	if r == "" && fallback {
		r = g.getRemoteFallback()
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

	_, err := g.RunCwd("fetch", remote, branch)
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
