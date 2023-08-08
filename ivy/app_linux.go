//go:build linux

package ivy

import (
	"bytes"
	"errors"
	"os/exec"
)

// for css
const DEFAULT_BG_OPACITY = 100

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
