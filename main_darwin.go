//go:build darwin

package main

import (
	"bytes"
	"errors"
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/options"
)

// for css
const DEFAULT_BG_OPACITY = 67

func appBgColor() *options.RGBA {
	return &options.RGBA{R: 7, G: 18, B: 34, A: 64}
}

func openDir(dir string) error {
	cmd := exec.Command("open", dir)
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
