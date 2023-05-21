package git

func (g *Git) PushTag(name string) error {
	_, err := g.RunCwd("push", name)
	return err
}
