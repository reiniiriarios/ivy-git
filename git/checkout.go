package git

import "strings"

func (g *Git) CheckoutCommit(hash string) error {
	_, err := g.RunCwd("checkout", hash)
	return err
}

func (g *Git) DiscardChanges(files ...string) error {
	return g.checkoutIndex(files)
}

// Forcefully update the working dir with info from the index for a set of files.
//
// Equivalent to `git checkout -- files`, but passes paths via stdin to avoid too long arguments.
// Will not yield errors for paths that don't exist in the index.
func (g *Git) checkoutIndex(paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	paths_joined := strings.Join(paths, "\x00")

	_, err := g.RunCwdStdin([]string{"checkout-index", "-f", "-u", "-q", "--stdin", "-z"}, paths_joined)

	if err != nil && errorCode(err) != ExitStatus1 {
		return err
	}

	return nil
}
