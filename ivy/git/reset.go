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
		_, err = g.run("reset", "--soft", hash)
	}
	// --mixed not handled
	return err
}

func (g *Git) ResetHead(hard bool) error {
	var err error
	if hard {
		_, err = g.run("reset", "--hard", "HEAD~1")
	} else {
		_, err = g.run("reset", "--soft", "HEAD~1")
	}
	// --mixed not handled
	return err
}

func (g *Git) DropCommit(hash string) error {
	if hash == "" {
		return errors.New("no commit hash specified")
	}
	_, err := g.run("rebase", "--onto", hash+"^", hash)
	return err
}
