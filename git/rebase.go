package git

func (g *Git) RebaseContinue() error {
	_, err := g.runWithOpts([]string{"rebase", "--continue"}, gitRunOpts{
		env: []env{{
			key:   "GIT_EDITOR",
			value: "true",
		}},
	})
	return err
}

func (g *Git) RebaseAbort() error {
	_, err := g.run("rebase", "--abort")
	return err
}

func (g *Git) RebaseSkip() error {
	_, err := g.run("rebase", "--skip")
	return err
}
