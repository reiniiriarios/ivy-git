//go:build !windows

package git

import "os/exec"

func hideCmdPrompt(cmd *exec.Cmd) {
	// not used
}
