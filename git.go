package main

import (
	"bytes"
	"errors"
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) Git(directory string, command ...string) (string, error) {
	// Run every git command in a specific directory.
	command = append([]string{"-C", directory}, command...)
	cmd := exec.Command("git", command[0:]...)

	// Git outputs much to stderr that isn't error, so treat it the same.
	// e.g.  git switch test
	//       STDERR Switched to branch 'test'
	//       STDOUT Your branch is up to date with 'origin/main'.
	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = cmd.Stdout

	err := cmd.Run()
	if err != nil {
		runtime.LogError(a.ctx, "err2")
		return outb.String(), err
	}

	return outb.String(), nil
}

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
