package git

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os/exec"
	"strings"
)

type Git struct {
	AppCtx context.Context
	Repo
}

type gitRunOpts struct {
	directory            string
	stdin                string
	ignore_errors        bool
	always_return_stderr bool
}

// Run a git command in the directory of the currently selected repo with the default options.
func (g *Git) run(command ...string) (string, error) {
	return g.runWithOpts(command, gitRunOpts{})
}

// Run a git command with the specified options.
func (g *Git) runWithOpts(command []string, opts gitRunOpts) (string, error) {
	if opts.directory == "" {
		if g.Repo == (Repo{}) || g.Repo.Directory == "" {
			return "", errors.New("no current git directory available")
		}
		opts.directory = g.Repo.Directory
	}

	// Run command in a specific directory.
	command = append([]string{"-C", opts.directory}, command...)
	cmd := exec.Command("git", command[0:]...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if opts.stdin != "" {
		err = cmd.Start()
		if err != nil {
			return "", g.ParseGitError(errb.String(), err)
		}

		// Because git does not continue before stdin is closed, this must be wrapped as so.
		// https://pkg.go.dev/os/exec#Cmd.StdinPipe
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, opts.stdin)
		}()

		err = cmd.Wait()

		if err != nil {
			return "", g.ParseGitError(errb.String(), err)
		}
	} else {
		err = cmd.Run()
	}

	// Git outputs much to stderr that isn't error.
	// e.g.  git switch test
	//       STDERR Switched to branch 'test'
	//       STDOUT Your branch is up to date with 'origin/main'.
	// Sometimes STDERR contains useful info we should return.
	// The error may be 'exit status 1', but if there is an
	// error, errb should contain the relevant information.
	if (opts.always_return_stderr && strings.TrimSpace(errb.String()) != "") || (!opts.ignore_errors && err != nil) {
		return outb.String(), g.ParseGitError(errb.String(), err)
	}

	return outb.String(), nil
}

// Initialize git in a specific directory.
func (g *Git) GitInit(directory string) error {
	_, err := g.runWithOpts([]string{"init"}, gitRunOpts{directory: directory})
	return err
}
