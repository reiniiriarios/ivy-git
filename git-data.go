package main

import (
	"ivy-git/git"
)

// This file has commands for the frontend to use to get
// data from the git package.

type DataResponse struct {
	Response string
	Message  string
	Data     any
}

func dataResponse(err error, data any) DataResponse {
	if err != nil {
		return DataResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}
	return DataResponse{
		Response: "success",
		Data:     data,
	}
}

// Get current branch for currently selected repo.
func (a *App) GetCurrentBranch() DataResponse {
	branch, err := a.Git.GetCurrentBranch()
	return dataResponse(err, git.Branch{
		Name: branch,
	})
}

// Get list of all branches for currently selected repo.
func (a *App) GetBranches() DataResponse {
	branches, err := a.Git.GetBranches()
	return dataResponse(err, branches)
}

// Switch branch on currently selected repo.
func (a *App) SwitchBranch(branch string, remote string) DataResponse {
	err := a.Git.SwitchBranch(branch, remote)
	if err != nil {
		return dataResponse(err, false)
	}

	b := git.Branch{
		Name: branch,
	}
	upstream, err := a.Git.GetBranchUpstream(branch)
	if err == nil {
		b.Upstream = upstream
	}

	return dataResponse(err, b)
}

// If branch exists locally.
func (a *App) BranchExists(name string) bool {
	return a.Git.BranchExists(name)
}

// Pull branch.
func (a *App) PullRemoteBranch(remote string, branch string, rebase bool) DataResponse {
	err := a.Git.PullRemoteBranch(remote, branch, rebase)
	return dataResponse(err, true)
}

// Get list of changed files.
func (a *App) GitListChanges() DataResponse {
	changesX, changesY, err := a.Git.GitListChanges()
	return dataResponse(err, struct {
		ChangesX []git.Change
		ChangesY []git.Change
	}{
		ChangesX: changesX,
		ChangesY: changesY,
	})
}

func (a *App) GetRemotes() DataResponse {
	remotes, err := a.Git.GetRemotes()
	if err != nil {
		return dataResponse(err, false)
	}
	current_remote, err := a.Git.GetRemoteForCurrentBranch()
	return dataResponse(err, struct {
		Remotes       []git.Remote
		CurrentRemote string
	}{
		Remotes:       remotes,
		CurrentRemote: current_remote,
	})
}

func (a *App) FetchRemote(remote string) DataResponse {
	err := a.Git.FetchRemote(remote)
	return dataResponse(err, true)
}

func (a *App) PushRemote(remote string) DataResponse {
	branch, err := a.Git.GetCurrentBranch()
	if err != nil {
		return dataResponse(err, true)
	}
	err = a.Git.PushRemoteBranch(remote, branch, false)
	return dataResponse(err, true)
}

func (a *App) PullRemote(remote string) DataResponse {
	branch, err := a.Git.GetCurrentBranch()
	if err != nil {
		return dataResponse(err, true)
	}
	// todo: set rebase flag depending on user settings
	err = a.Git.PullRemoteBranch(remote, branch, true)
	return dataResponse(err, true)
}

func (a *App) PushBranch(branch string) DataResponse {
	err := a.Git.PushBranch(branch)
	return dataResponse(err, true)
}

func (a *App) PullBranch(branch string) DataResponse {
	// todo: set rebase flag depending on user settings
	err := a.Git.PullBranch(branch, true)
	return dataResponse(err, true)
}

// Hard reset the current branch to its default remote.
func (a *App) ResetBranchToRemote(branch string) DataResponse {
	err := a.Git.ResetBranchToRemote(branch)
	return dataResponse(err, true)
}

// Get additional commit details not listed in the table.
func (a *App) GetCommitDetails(hash string) DataResponse {
	commit, err := a.Git.GetCommitDetails(hash)
	return dataResponse(err, commit)
}

// Get commit diff summary from diff-tree --numstat and --name-status.
func (a *App) GetCommitDiffSummary(hash string) DataResponse {
	files, err := a.Git.GetCommitDiffSummary(hash)
	return dataResponse(err, files)
}

// Get list of commits and all associated details for display.
func (a *App) GetCommitList(limit uint64, offset uint64) DataResponse {
	HEAD, commits, graph, err := a.Git.GetCommitsAndGraph(limit, offset)
	return dataResponse(err, struct {
		HEAD    git.Ref
		Commits []git.Commit
		Graph   git.Graph
	}{
		HEAD:    HEAD,
		Commits: commits,
		Graph:   graph,
	})
}

// Get list of commit hashes and their signature status.
func (a *App) GetCommitsSignStatus(limit uint64, offset uint64) DataResponse {
	commits, err := a.Git.GetCommitsSignStatus(limit, offset)
	return dataResponse(err, commits)
}

func (a *App) GetCommitSignature(hash string) DataResponse {
	commit, err := a.Git.GetCommitSignature(hash)
	return dataResponse(err, commit)
}

func (a *App) DeleteBranch(branch string, force bool, remote bool) DataResponse {
	err := a.Git.DeleteBranch(branch, force, remote)
	return dataResponse(err, true)
}

func (a *App) DeleteRemoteBranch(branch string, remote string, force bool) DataResponse {
	err := a.Git.DeleteRemoteBranch(branch, remote, force)
	return dataResponse(err, true)
}

func (a *App) RenameBranch(branch string, new_name string) DataResponse {
	err := a.Git.RenameBranch(branch, new_name)
	return dataResponse(err, true)
}

// Rebase current branch on branch.
func (a *App) RebaseOnBranch(branch string) DataResponse {
	err := a.Git.RebaseOnBranch(branch)
	return dataResponse(err, true)
}

func (a *App) DeleteTag(name string) DataResponse {
	err := a.Git.DeleteTag(name)
	return dataResponse(err, true)
}

func (a *App) PushTag(name string) DataResponse {
	err := a.Git.PushTag(name)
	return dataResponse(err, true)
}

func (a *App) AddTag(hash string, name string, annotated bool, message string, push bool) DataResponse {
	err := a.Git.AddTag(hash, name, annotated, message, push)
	return dataResponse(err, true)
}

func (a *App) MergeCommit(target_branch string, no_commit bool, no_ff bool) DataResponse {
	err := a.Git.MergeCommit(target_branch, no_commit, no_ff)
	return dataResponse(err, true)
}

func (a *App) MergeSquash(target_branch string) DataResponse {
	err := a.Git.MergeSquash(target_branch)
	return dataResponse(err, true)
}

func (a *App) MergeRebase(target_branch string) DataResponse {
	err := a.Git.MergeRebase(target_branch)
	return dataResponse(err, true)
}

func (a *App) MergeFastForward(target_branch string) DataResponse {
	err := a.Git.MergeFastForward(target_branch)
	return dataResponse(err, true)
}

func (a *App) CreateBranch(name string, at_hash string, checkout bool) DataResponse {
	err := a.Git.CreateBranch(name, at_hash, checkout)
	return dataResponse(err, true)
}

func (a *App) StageFile(file string) DataResponse {
	err := a.Git.StageFiles(file)
	return dataResponse(err, true)
}

func (a *App) UnstageFile(file string) DataResponse {
	err := a.Git.UnstageFile(file)
	return dataResponse(err, true)
}

func (a *App) StageAll() DataResponse {
	err := a.Git.StageAll()
	return dataResponse(err, true)
}

func (a *App) UnstageAll() DataResponse {
	err := a.Git.UnstageAll()
	return dataResponse(err, true)
}

func (a *App) NumBranches() DataResponse {
	n := a.Git.NumBranches()
	return dataResponse(nil, n)
}

func (a *App) NumTags() DataResponse {
	n := a.Git.NumTags()
	return dataResponse(nil, n)
}

func (a *App) NumMainBranchCommits() DataResponse {
	n, err := a.Git.NumCommitsOnBranch(a.Git.Repo.Main)
	return dataResponse(err, n)
}

func (a *App) MakeCommit(subject string, body string) DataResponse {
	err := a.Git.MakeCommit(subject, body)
	return dataResponse(err, true)
}

func (a *App) CheckoutCommit(hash string) DataResponse {
	err := a.Git.CheckoutCommit(hash)
	return dataResponse(err, true)
}

func (a *App) HardReset(hash string) DataResponse {
	err := a.Git.HardReset(hash)
	return dataResponse(err, true)
}

func (a *App) RevertCommit(hash string) DataResponse {
	err := a.Git.RevertCommit(hash)
	return dataResponse(err, true)
}
