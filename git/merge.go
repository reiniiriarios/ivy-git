package git

import (
	"errors"
)

func (g *Git) MergeCommit(target_branch string, no_commit bool, no_ff bool) error {
	if target_branch == "" {
		return errors.New("no branch name specified")
	}

	var err error
	cmd := []string{"merge"}
	if no_commit {
		cmd = append(cmd, "--no-commit")
	}
	if no_ff {
		cmd = append(cmd, "--no-ff")
	}
	cmd = append(cmd, target_branch)
	res, err := g.RunCwd(cmd...)
	if err != nil {
		return err
	}
	res = parseOneLine(res)
	if res == "Already up to date." {
		return errors.New(res)
	}
	// todo: conflicts

	return nil
}

func (g *Git) MergeSquash(target_branch string) error {
	if target_branch == "" {
		return errors.New("no branch name specified")
	}

	_, err := g.RunCwd("merge", "--squash", target_branch)
	if err != nil {
		return err
	}
	// todo: conflicts, display commit message

	return nil
}

func (g *Git) MergeRebase(target_branch string) error {
	if target_branch == "" {
		return errors.New("no branch name specified")
	}

	_, err := g.RunCwd("rebase", "--merge", target_branch)
	if err != nil {
		return err
	}
	// todo: conflicts

	return nil
}

func (g *Git) MergeFastForward(target_branch string) error {
	if target_branch == "" {
		return errors.New("no branch name specified")
	}

	_, err := g.RunCwd("merge", "--ff-only", target_branch)
	if err != nil {
		return err
	}
	// with ff-only, there will be no conflicts

	return nil
}
