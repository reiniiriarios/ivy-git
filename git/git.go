package git

import (
	"bytes"
	"errors"
	"os/exec"
)

type Git struct {
	Repo
}

// Run a git command in a specific directory.
func (g *Git) Run(directory string, command ...string) (string, error) {
	// Run every git command in a specific directory.
	command = append([]string{"-C", directory}, command...)
	cmd := exec.Command("git", command[0:]...)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	// Git outputs much to stderr that isn't error.
	// e.g.  git switch test
	//       STDERR Switched to branch 'test'
	//       STDOUT Your branch is up to date with 'origin/main'.
	// The error may be 'exit status 1', but if there is an
	// error, errb should contain the relevant information.
	err := cmd.Run()
	if err != nil {
		return outb.String(), errors.New(errb.String())
	}

	return outb.String(), nil
}

// Run a git command in the directory of the currently selected repo.
func (g *Git) RunCwd(command ...string) (string, error) {
	if g.Repo == (Repo{}) {
		return "", errors.New("no current git directory available")
	}
	return g.Run(g.Repo.Directory, command...)
}
