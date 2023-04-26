package main

import (
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

func (a *App) GetCurrentBranch(repo_id string) BranchResponse {
	repo, exists := a.RepoSaveData.Repos[repo_id]
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

// func (a *App) GetBranches() map[string]Branch {
// 	return a.RepoSaveData.Repos
// }

// func (a *App) SwitchBranch(repo string) {

// }
