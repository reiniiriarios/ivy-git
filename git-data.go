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
func (a *App) SwitchBranch(branch string) BranchResponse {
	err := a.Git.SwitchBranch(branch)
	if err != nil {
		return BranchResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	b := git.Branch{
		Name: branch,
	}

	upstream, err := a.Git.GetBranchUpstream(branch)
	if err == nil {
		b.Upstream = upstream
	}

	return BranchResponse{
		Response: "success",
		Branch:   b,
	}
}

// If branch exists locally.
func (a *App) BranchExists(name string) bool {
	return a.Git.BranchExists(name)
}

// Pull branch.
func (a *App) PullRemoteBranch(remote string, branch string, rebase bool) GenericResponse {
	err := a.Git.PullRemoteBranch(remote, branch, rebase)
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

func (a *App) PushRemote(remote string) GenericResponse {
	branch, err := a.Git.GetCurrentBranch()
	if err != nil {
		return GenericResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}
	err = a.Git.PushRemoteBranch(remote, branch)
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

func (a *App) PullRemote(remote string) GenericResponse {
	branch, err := a.Git.GetCurrentBranch()
	if err != nil {
		return GenericResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}
	// todo: set rebase flag depending on user settings
	err = a.Git.PullRemoteBranch(remote, branch, true)
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

func (a *App) PushBranch(branch string) GenericResponse {
	err := a.Git.PushBranch(branch)
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

func (a *App) PullBranch(branch string) GenericResponse {
	// todo: set rebase flag depending on user settings
	err := a.Git.PullBranch(branch, true)
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

func (a *App) ResetBranchToRemote(branch string) GenericResponse {
	err := a.Git.ResetBranchToRemote(branch)
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
func (a *App) GetCommitList(limit uint64, offset uint64) CommitsResponse {
	HEAD, commits, graph, err := a.Git.GetCommitsAndGraph(limit, offset)
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

type CommitsSignResponse struct {
	Response string
	Message  string
	Commits  git.CommitsSigned
}

// FRONTEND: Get list of commit hashes and their signature status.
func (a *App) GetCommitsSignStatus(limit uint64, offset uint64) CommitsSignResponse {
	commits, err := a.Git.GetCommitsSignStatus(limit, offset)
	if err != nil {
		return CommitsSignResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}
	return CommitsSignResponse{
		Response: "success",
		Commits:  commits,
	}
}

type CommitSignResponse struct {
	Response  string
	Message   string
	Signature git.CommitSignature
}

// FRONTEND: Get commit signature data.
func (a *App) GetCommitSignature(hash string) CommitSignResponse {
	commit, err := a.Git.GetCommitSignature(hash)
	if err != nil {
		return CommitSignResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}
	return CommitSignResponse{
		Response:  "success",
		Signature: commit,
	}
}

// FRONTEND: Delete a branch.
func (a *App) DeleteBranch(branch string, force bool, remote bool) GenericResponse {
	err := a.Git.DeleteBranch(branch, force, remote)
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

// FRONTEND: Delete a remote branch.
func (a *App) DeleteRemoteBranch(branch string, remote string, force bool) GenericResponse {
	err := a.Git.DeleteRemoteBranch(branch, remote, force)
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

// FRONTEND: Rename a branch.
func (a *App) RenameBranch(branch string, new_name string) GenericResponse {
	err := a.Git.RenameBranch(branch, new_name)
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

// FRONTEND: Rebase current branch on branch.
func (a *App) RebaseOnBranch(branch string) GenericResponse {
	err := a.Git.RebaseOnBranch(branch)
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
