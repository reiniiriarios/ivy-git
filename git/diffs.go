package git

import "strings"

func (g *Git) GetUncommittedDiff() (string, error) {
	diff, err := g.RunCwd("--no-pager", "diff", "HEAD^")
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) GetDiffRemoteCurrent() (string, error) {
	branch, err := g.GetCurrentBranch()
	if err != nil {
		return "", err
	}
	remote, err := g.getBranchRemote(branch)
	if err != nil {
		return "", err
	}
	diff, err := g.GetDiffRemote(remote, branch)
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) GetDiffRemote(remote string, branch string) (string, error) {
	diff, err := g.RunCwd("--no-pager", "diff", remote+"/"+branch)
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) GetDiffStaged() (string, error) {
	diff, err := g.RunCwd("--no-pager", "diff", "--staged")
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) findMergeBase(hash1 string, hash2 string) (string, error) {
	b, err := g.RunCwd("merge-base", hash1, hash2)
	if err != nil {
		return "", err
	}
	b = strings.Trim(strings.Trim(b, "\r"), "\n")
	return b, nil
}
