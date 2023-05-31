package git

import (
	"bytes"
	"context"
	"errors"
	"io"
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

// Run a git command in the current directory and pipe stdin to it.
func (g *Git) RunCwdStdin(command []string, input string) (string, error) {
	if g.Repo == (Repo{}) {
		return "", errors.New("no current git directory available")
	}
	return g.runCmdStdin(g.Repo.Directory, command, input)
}

// Run a git command in the directory of the currently selected repo. Ignore errors.
func (g *Git) RunCwdNoError(command ...string) (string, error) {
	if g.Repo == (Repo{}) {
		return "", errors.New("no current git directory available")
	}
	return g.runCmd(g.Repo.Directory, command, true)
}

// Run a git command.
func (g *Git) runCmd(directory string, command []string, ignore_error bool) (string, error) {
	// Run command in a specific directory.
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

// Run a git command and pipe stdin to it.
func (g *Git) runCmdStdin(directory string, command []string, input string) (string, error) {
	// Run command in a specific directory.
	command = append([]string{"-C", directory}, command...)
	cmd := exec.Command("git", command[0:]...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err = cmd.Start()
	if err != nil {
		return "", g.ParseGitError(errb.String(), err)
	}

	// Because git does not continue before stdin is closed, this must be wrapped as so.
	// https://pkg.go.dev/os/exec#Cmd.StdinPipe
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input)
	}()

	err = cmd.Wait()

	if err != nil {
		return "", g.ParseGitError(errb.String(), err)
	}

	println(outb.String())
	println(errb.String())

	return outb.String(), nil
}

// Initialize git in a specific directory.
func (g *Git) GitInit(directory string) error {
	cmd := []string{"-C", directory, "init"}
	_, err := g.runCmd(g.Repo.Directory, cmd, false)
	return err
}
