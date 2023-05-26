package git

import "errors"

func (g *Git) HardReset(hash string) error {
	if hash == "" {
		return errors.New("no commit hash specified")
	}
	_, err := g.RunCwd("reset", "--hard", hash)
	return err
}

func (g *Git) DropCommit(hash string) error {
	if hash == "" {
		return errors.New("no commit hash specified")
	}
	_, err := g.RunCwd("rebase", "--onto", hash+"^", hash)
	return err
}
