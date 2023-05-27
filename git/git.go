package git

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
)

type Git struct {
	AppCtx context.Context
	Repo
}

// Run a git command in a specific directory.
func (g *Git) Run(directory string, command ...string) (string, error) {
	return g.runCmd(directory, command, false)
}

// Run a git command in the directory of the currently selected repo.
func (g *Git) RunCwd(command ...string) (string, error) {
	if g.Repo == (Repo{}) {
		return "", errors.New("no current git directory available")
	}
	return g.runCmd(g.Repo.Directory, command, false)
}

// Run a git command in the directory of the currently selected repo. Ignore errors.
func (g *Git) RunCwdNoError(command ...string) (string, error) {
	if g.Repo == (Repo{}) {
		return "", errors.New("no current git directory available")
	}
	return g.runCmd(g.Repo.Directory, command, true)
}

func (g *Git) runCmd(directory string, command []string, ignore_error bool) (string, error) {
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
	if !ignore_error && err != nil {
		return outb.String(), g.ParseGitError(errb.String(), err)
	}

	return outb.String(), nil
}
