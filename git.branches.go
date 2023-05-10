package main

import (
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
