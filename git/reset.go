package git

import "errors"

func (g *Git) ResetToCommit(hash string, hard bool) error {
	if hash == "" {
		return errors.New("no commit hash specified")
	}
	var err error
	if hard {
		_, err = g.run("reset", "--hard", hash)
	} else {
		_, err = g.run("reset", hash)
	}
	return err
}

func (g *Git) DropCommit(hash string) error {
	if hash == "" {
		return errors.New("no commit hash specified")
	}
	_, err := g.run("rebase", "--onto", hash+"^", hash)
	return err
}
