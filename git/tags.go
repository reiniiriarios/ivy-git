package git

func (g *Git) PushTag(name string) error {
	remote, err := g.getRemoteForCurrentBranch()
	if err != nil {
		return err
	}
	_, err = g.RunCwd("push", remote, name)
	return err
}

func (g *Git) DeleteTag(name string) error {
	remote, err := g.getRemoteForCurrentBranch()
	if err != nil {
		return err
	}
	_, err = g.RunCwd("tag", "-d", name)
	if err != nil {
		return err
	}
	_, err = g.RunCwd("push", remote, ":refs/tags/"+name)
	return err
}

func (g *Git) AddTag(hash string, name string, annotated bool, message string, push bool) error {
	var err error
	if annotated {
		if message == "" {
			_, err = g.RunCwd("tag", "-a", name, hash)
		} else {
			_, err = g.RunCwd("tag", "-a", name, hash, "-m", message)
		}
	} else {
		_, err = g.RunCwd("tag", name, hash)
	}
	if err != nil {
		return err
	}

	if push {
		err = g.PushTag(name)
		return err
	}

	return nil
}
