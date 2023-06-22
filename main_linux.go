//go:build linux

package main

import (
	"bytes"
	"errors"
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/options"
)

func appBgColor() *options.RGBA {
	return &options.RGBA{R: 7, G: 18, B: 34, A: 255}
}

func openDir(dir string) error {
	// may need gnome-open or kd-open on some systems
	cmd := exec.Command("xdg-open", dir)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		if errb.String() != "" {
			return errors.New(errb.String())
		}
		return err
	}
	return nil
}
