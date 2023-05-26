package git

import "errors"

func (g *Git) ResetToCommit(hash string, hard bool) error {
	if hash == "" {
		return errors.New("no commit hash specified")
	}
	var err error
	if hard {
		_, err = g.RunCwd("reset", "--hard", hash)
	} else {
		_, err = g.RunCwd("reset", hash)
	}
	return err
}

func (g *Git) DropCommit(hash string) error {
	if hash == "" {
		return errors.New("no commit hash specified")
	}
	_, err := g.RunCwd("rebase", "--onto", hash+"^", hash)
	return err
}
