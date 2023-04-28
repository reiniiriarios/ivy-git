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

func (a *App) GetCurrentBranch() BranchResponse {
	repo, exists := a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo]
	if !exists {
		return BranchResponse{
			Response: "error",
			Message:  "Repo not found.",
		}
	}

	branch, err := a.Git(repo.Directory, "rev-parse", "--abbrev-ref", "HEAD")
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

func (a *App) GetBranches() BranchesResponse {
	repo, exists := a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo]
	if !exists {
		return BranchesResponse{
			Response: "error",
			Message:  "Repo not found.",
		}
	}

	branches, err := a.Git(repo.Directory, "branch", "--list", "--format", "'%(refname:short)'")
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return BranchesResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	branch_list := make(map[string]Branch)
	bs := strings.Split(strings.ReplaceAll(branches, "\r\n", "\n"), "\n")
	for _, branch := range bs {
		branch = strings.Trim(branch, "'")
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

func (a *App) SwitchBranch(branch string) GenericResponse {
	repo, exists := a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo]
	if !exists {
		return GenericResponse{
			Response: "error",
			Message:  "Repo not found.",
		}
	}

	_, err := a.Git(repo.Directory, "switch", branch)
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
