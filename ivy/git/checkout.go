package git

import (
	"os"
	"path/filepath"
	"strings"
)

func (g *Git) CheckoutCommit(hash string) error {
	_, err := g.run("checkout", hash)
	return err
}

func (g *Git) DiscardChanges(files ...string) error {
	for _, file := range files {
		_, y, err := g.getFileStatus(file)
		if err != nil {
			return err
		}
		if y == FileUntracked {
			err := os.Remove(filepath.Join(g.Repo.Directory, file))
			if err != nil {
				return err
			}
		} else {
			err := g.checkoutIndex(file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Forcefully update the working dir with info from the index for a set of files.
//
// Equivalent to `git checkout -- files`, but passes paths via stdin to avoid too long arguments.
// Will not yield errors for paths that don't exist in the index.
func (g *Git) checkoutIndex(paths ...string) error {
	if len(paths) == 0 {
		return nil
	}
	paths_joined := strings.Join(paths, "\x00")

	_, err := g.runWithOpts([]string{"checkout-index", "-f", "-u", "-q", "--stdin", "-z"}, gitRunOpts{stdin: paths_joined})

	if err != nil && errorCode(err) != ExitStatus1 {
		return err
	}

	return nil
}
