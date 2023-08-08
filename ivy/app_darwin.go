//go:build darwin

package ivy

import (
	"bytes"
	"errors"
	"os/exec"
)

// for css
const DEFAULT_BG_OPACITY = 67

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
