//go:build windows

package git

import (
	"os/exec"
	"syscall"
)

func hideCmdPrompt(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
