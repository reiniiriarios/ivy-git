package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Branch struct {
	Name string
}

type BranchResponse struct {
	Response string
	Message  string
	Branch   Branch
}

type BranchesResponse struct {
	Response string
	Message  string
	Branches map[string]Branch
}

// FRONTEND: Get current branch for currently selected repo.
func (a *App) GetCurrentBranch() BranchResponse {
	branch, err := a.GitCwd("rev-parse", "--abbrev-ref", "HEAD")
	branch = strings.ReplaceAll(strings.ReplaceAll(branch, "\r", ""), "\n", "")
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return BranchResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	return BranchResponse{
		Response: "success",
		Branch: Branch{
			Name: branch,
		},
	}
}

// FRONTEND: Get list of all branches for currently selected repo.
func (a *App) GetBranches() BranchesResponse {
	branches, err := a.GitCwd("branch", "--list", "--format", "'%(refname:short)'")
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return BranchesResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	branch_list := make(map[string]Branch)
	bs := a.getLines(branches)
	for _, branch := range bs {
		if strings.Trim(branch, " ") != "" {
			branch_list[branch] = Branch{
				Name: branch,
			}
		}
	}

	return BranchesResponse{
		Response: "success",
		Branches: branch_list,
	}
}

// FRONTEND: Switch branch on currently selected repo.
func (a *App) SwitchBranch(branch string) GenericResponse {
	_, err := a.GitCwd("checkout", branch)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return GenericResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	return GenericResponse{
		Response: "success",
	}
}

// If branch exists locally.
func (a *App) BranchExists(name string) bool {
	_, err := a.GitCwd("rev-parse", "--verify", name)
	return err == nil
}

// Get commits ahead and behind branch is from specific remote.
func (a *App) getAheadBehind(branch string, remote string) (uint32, uint32, error) {
	rl, err := a.GitCwd("rev-list", "--left-right", "--count", branch+"..."+remote+"/"+branch)
	if err != nil {
		return 0, 0, err
	}
	ab := strings.Fields(rl)
	if len(ab) != 2 {
		return 0, 0, errors.New("error parsing rev-list --left-right")
	}
	ahead, _ := strconv.ParseInt(ab[0], 10, 32)
	behind, _ := strconv.ParseInt(ab[1], 10, 32)
	return uint32(ahead), uint32(behind), nil
}

func (a *App) PushBranch(remote string, branch string) GenericResponse {
	_, err := a.GitCwd("push", remote, branch+":"+branch)
	if err != nil {
		return GenericResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}
	return GenericResponse{
		Response: "success",
	}
}

func (a *App) PullBranch(remote string, branch string, rebase bool) GenericResponse {
	var err error
	if rebase {
		_, err = a.GitCwd("pull", remote, branch+":"+branch, "--rebase")
	} else {
		_, err = a.GitCwd("pull", remote, branch+":"+branch)
	}
	if err != nil {
		return GenericResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}
	return GenericResponse{
		Response: "success",
	}
}
