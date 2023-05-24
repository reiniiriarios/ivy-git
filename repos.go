package main

import (
	"errors"
	"ivy-git/git"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// This file deals with repos added to the app and their settings.

type RepoResponse struct {
	Response string
	Message  string
	Id       string
	Repo     git.Repo
}

func (a *App) isCurrentRepo() bool {
	return a.Git.Repo != (git.Repo{})
}

// Get repo information.
func (a *App) GetRepos() map[string]git.Repo {
	return a.RepoSaveData.Repos
}

// Update currently selected repo.
func (a *App) UpdateSelectedRepo(repo string) DataResponse {
	if repo == "" {
		a.RepoSaveData.CurrentRepo = ""
		a.Git.Repo = git.Repo{}
		a.saveRepoData()
		return dataResponse(nil, false)
	}

	if _, exists := a.RepoSaveData.Repos[repo]; !exists {
		return dataResponse(errors.New("repo not found in list"), false)
	}

	if !a.Git.IsGitRepo(a.RepoSaveData.Repos[repo].Directory) {
		return dataResponse(errors.New("directory not found, or not identifiable as git repo"), false)
	}
	a.RepoSaveData.CurrentRepo = repo
	a.saveRepoData()
	a.Git.Repo = a.RepoSaveData.Repos[repo]
	return dataResponse(nil, false)
}

// Get the currently selected repo.
func (a *App) GetSelectedRepo() string {
	return a.RepoSaveData.CurrentRepo
}

// Add a new repo.
func (a *App) AddRepo() RepoResponse {
	d, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Git Repository",
	})

	if err != nil {
		return RepoResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	if d == "" {
		return RepoResponse{
			Response: "none",
		}
	}

	if !a.Git.IsGitRepo(d) {
		return RepoResponse{
			Response: "error",
			Message:  "The directory selected is not a git repository.",
		}
	}

	for _, r := range a.RepoSaveData.Repos {
		if r.Directory == d {
			return RepoResponse{
				Response: "error",
				Message:  "The directory selected is already added.",
			}
		}
	}

	id := uuid.New().String()
	newRepo := git.Repo{
		Name:      filepath.Base(d),
		Directory: d,
		Main:      a.Git.NameOfMainBranchForRepo(d),
	}

	if a.RepoSaveData.Repos == nil {
		a.RepoSaveData.Repos = make(map[string]git.Repo)
	}
	a.RepoSaveData.Repos[id] = newRepo
	a.saveRepoData()

	return RepoResponse{
		Response: "success",
		Id:       id,
		Repo:     newRepo,
	}
}

// Remove a repo from the list.
func (a *App) RemoveRepo(id string) map[string]git.Repo {
	delete(a.RepoSaveData.Repos, id)
	a.saveRepoData()
	return a.RepoSaveData.Repos
}
