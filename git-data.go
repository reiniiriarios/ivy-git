package main

import (
	"ivy-git/git"
)

// This file has commands for the frontend to use to get
// data from the git package.

type GenericResponse struct {
	Response string
	Message  string
}

type BranchResponse struct {
	Response string
	Message  string
	Branch   git.Branch
}

// Get current branch for currently selected repo.
func (a *App) GetCurrentBranch() BranchResponse {
	branch, err := a.Git.GetCurrentBranch()
	if err != nil {
		return BranchResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	return BranchResponse{
		Response: "success",
		Branch: git.Branch{
			Name: branch,
		},
	}
}

type BranchesResponse struct {
	Response string
	Message  string
	Branches map[string]git.Branch
}

// Get list of all branches for currently selected repo.
func (a *App) GetBranches() BranchesResponse {
	branches, err := a.Git.GetBranches()
	if err != nil {
		return BranchesResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	return BranchesResponse{
		Response: "success",
		Branches: branches,
	}
}

// Switch branch on currently selected repo.
func (a *App) SwitchBranch(branch string) GenericResponse {
	err := a.Git.SwitchBranch(branch)
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

// If branch exists locally.
func (a *App) BranchExists(name string) bool {
	return a.Git.BranchExists(name)
}

// Pull branch.
func (a *App) PullBranch(remote string, branch string, rebase bool) GenericResponse {
	err := a.Git.PullBranch(remote, branch, rebase)
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

type ChangesResponse struct {
	Response string
	Message  string
	ChangesX []git.Change
	ChangesY []git.Change
}

// Get list of changed files.
func (a *App) GitListChanges() ChangesResponse {
	changesX, changesY, err := a.Git.GitListChanges()
	if err != nil {
		return ChangesResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	return ChangesResponse{
		Response: "success",
		ChangesX: changesX,
		ChangesY: changesY,
	}
}

type RemotesResponse struct {
	Response string
	Message  string
	Remotes  []git.Remote
}

func (a *App) GetRemotes() RemotesResponse {
	remotes, err := a.Git.GetRemotes()
	if err != nil {
		return RemotesResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	return RemotesResponse{
		Response: "success",
		Remotes:  remotes,
	}
}

func (a *App) FetchRemote(remote string) GenericResponse {
	err := a.Git.FetchRemote(remote)
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

type CommitResponse struct {
	Response string
	Message  string
	Commit   git.CommitAddl
}

// FRONTEND: Get additional commit details not listed in the table.
func (a *App) GetCommitDetails(hash string) CommitResponse {
	commit, err := a.Git.GetCommitDetails(hash)
	if err != nil {
		return CommitResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}
	return CommitResponse{
		Response: "success",
		Commit:   commit,
	}
}

type CommitDiffSummaryResponse struct {
	Response string
	Message  string
	Files    git.FileStatDir
}

// FRONTEND: Get commit diff summary from diff-tree --numstat and --name-status.
func (a *App) GetCommitDiffSummary(hash string) CommitDiffSummaryResponse {
	files, err := a.Git.GetCommitDiffSummary(hash)
	if err != nil {
		return CommitDiffSummaryResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	return CommitDiffSummaryResponse{
		Response: "success",
		Files:    files,
	}
}

type CommitsResponse struct {
	Response string
	Message  string
	HEAD     git.Ref
	Commits  []git.Commit
	Graph    git.Graph
}

// FRONTEND: Get list of commits and all associated details for display.
func (a *App) GetCommitList() CommitsResponse {
	HEAD, commits, graph, err := a.Git.GetCommitsAndGraph()
	if err != nil {
		return CommitsResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	return CommitsResponse{
		Response: "success",
		HEAD:     HEAD,
		Commits:  commits,
		Graph:    graph,
	}
}
