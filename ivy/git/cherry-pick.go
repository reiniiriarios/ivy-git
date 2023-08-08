package git

func (g *Git) CherryPick(hash string, record bool, no_commit bool) error {
	cmd := []string{"cherry-pick", "--allow-empty"}
	if record {
		cmd = append(cmd, "-x")
	}
	if no_commit {
		cmd = append(cmd, "--no-commit")
	}
	cmd = append(cmd, hash)

	_, err := g.run(cmd...)
	return err
}

func (g *Git) CherryPickContinue() error {
	_, err := g.runWithOpts([]string{"cherry-pick", "--continue"}, gitRunOpts{
		env: []env{{
			key:   "GIT_EDITOR",
			value: "true",
		}},
	})
	return err
}

func (g *Git) CherryPickAbort() error {
	_, err := g.run("cherry-pick", "--abort")
	return err
}

func (g *Git) CherryPickSkip() error {
	_, err := g.run("cherry-pick", "--skip")
	return err
}
