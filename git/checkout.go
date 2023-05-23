package git

func (g *Git) CheckoutCommit(hash string) error {
	_, err := g.RunCwd("checkout", hash)
	return err
}
