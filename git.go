package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Run a git command in a specific directory.
func (a *App) Git(directory string, command ...string) (string, error) {
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
		runtime.LogError(a.ctx, err.Error())
		runtime.LogError(a.ctx, errb.String())
		return outb.String(), errors.New(errb.String())
	}

	return outb.String(), nil
}

// Run a git command in the directory of the currently selected repo.
func (a *App) GitCwd(command ...string) (string, error) {
	repo, exists := a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo]
	if !exists {
		return "", errors.New("no current git directory available")
	}
	return a.Git(repo.Directory, command...)
}

func (a *App) IsGitRepo(directory string) bool {
	r, err := a.Git(directory, "rev-parse")
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return false
	}

	return r == ""
}

// Parse a multiline string into an array.
func (a *App) getLines(s string) []string {
	var r []string
	l := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	for _, v := range l {
		v = strings.Trim(v, "'")
		r = append(r, v)
	}
	return r
}
