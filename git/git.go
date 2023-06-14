package git

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
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
	env                  []env
}

type env struct {
	key   string
	value string
	empty bool
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

	// Set environment variables, if any.
	if len(opts.env) > 0 {
		env_save, err := setEnvVariables(opts.env)
		if err != nil {
			return "", err
		}
		// Environment variables will be reset at end.
		defer resetEnvVariables(env_save)
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

// Set environment variables, returning their current values.
func setEnvVariables(envs []env) ([]env, error) {
	// If we're setting environment variables, first save the current values.
	env_save := []env{}
	if len(envs) > 0 {
		for _, e := range envs {
			value, exists := os.LookupEnv(e.key)
			env_save = append(env_save, env{
				key:   e.key,
				value: value,
				empty: !exists,
			})
			err := os.Setenv(e.key, e.value)
			if err != nil {
				// In case some were set, reset them. Ignore errors here.
				resetEnvVariables(env_save)
				return nil, err
			}
		}
	}
	return env_save, nil
}

// Reset (or unset) environment variables to the given values
func resetEnvVariables(envs []env) error {
	if len(envs) > 0 {
		var err error
		for _, e := range envs {
			if e.empty {
				err = os.Unsetenv(e.key)
			} else {
				err = os.Setenv(e.key, e.value)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}
