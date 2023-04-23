package main

import (
	"bytes"
	"errors"
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) Git(command ...string) (string, error) {
	cmd := exec.Command("git", command[0:]...)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err2 := cmd.Run()
	if err2 != nil {
		return outb.String(), err2
	}
	if errb.String() != "" {
		return outb.String(), errors.New(errb.String())
	}

	return outb.String(), nil
}

func (a *App) IsGitRepo(directory string) bool {
	r, err := a.Git("-C", directory, "rev-parse")
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return false
	}

	if r == "" {
		return true
	}
	return false
}
