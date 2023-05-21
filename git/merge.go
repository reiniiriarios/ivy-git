package git

func (g *Git) MergeCommit(target_branch string, no_commit bool, no_ff bool) error {
	var err error
	cmd := []string{"merge"}
	if no_commit {
		cmd = append(cmd, "--no-commit")
	}
	if no_ff {
		cmd = append(cmd, "--no-ff")
	}
	cmd = append(cmd, target_branch)
	_, err = g.RunCwd(cmd...)
	if err != nil {
		return err
	}
	// todo: conflicts

	return nil
}

func (g *Git) MergeSquash(target_branch string) error {
	_, err := g.RunCwd("merge", "--squash", target_branch)
	if err != nil {
		return err
	}
	// todo: conflicts, display commit message

	return nil
}

func (g *Git) MergeRebase(target_branch string) error {
	_, err := g.RunCwd("rebase", "--merge", target_branch)
	if err != nil {
		return err
	}
	// todo: conflicts

	return nil
}

func (g *Git) MergeFastForward(target_branch string) error {
	_, err := g.RunCwd("merge", "--ff-only", target_branch)
	if err != nil {
		return err
	}
	// with ff-only, there will be no conflicts

	return nil
}
