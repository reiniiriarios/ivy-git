package main

import (
	"fmt"
	"ivy-git/files"
	"ivy-git/git"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// This file has commands for the frontend to use to get data from the git package.
// Minimal processing to translate reasonable data to a standard DataResponse, along
// with occasional combination of methods when those methods don't make sense to
// combine in the context of the git package.

type DataResponse struct {
	Response       string
	Message        string
	FetchCommitMsg bool
	Data           any
}

func dataResponse(err error, data any) DataResponse {
	if err != nil {
		fmt.Printf("%g", err)
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

func (a *App) fileTooLarge(file string) bool {
	if a.RepoSaveData.CurrentRepo == "" {
		runtime.LogError(a.ctx, fmt.Sprintf("trying to find size of file %s without current repo selected", file))
		return false
	}
	f := filepath.Join(a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo].Directory, file)
	return files.FileTooLarge(f)
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
	changesX, changesY, changesC, err := a.Git.GitListChanges()
	return dataResponse(err, struct {
		ChangesX map[string]*git.Change
		ChangesY map[string]*git.Change
		ChangesC map[string]*git.Change
	}{
		ChangesX: changesX,
		ChangesY: changesY,
		ChangesC: changesC,
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
	if err == nil {
		err = a.Git.PushRemoteBranch(remote, branch, false, false)
	}
	return dataResponse(err, true)
}

func (a *App) PullRemote(remote string) DataResponse {
	branch, err := a.Git.GetCurrentBranch()
	if err == nil {
		// todo: set rebase flag depending on user settings
		err = a.Git.PullRemoteBranch(remote, branch, true)
	}
	return dataResponse(err, true)
}

func (a *App) PushBranch(branch string, force bool) DataResponse {
	must_force, err := a.Git.PushBranch(branch, force)
	if must_force {
		return DataResponse{
			Response: "must-force",
		}
	}
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
	if !a.Git.HasCommits(a.Git.Repo.Directory) {
		return DataResponse{
			Response: "no-commits",
		}
	}

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
	must_force, err := a.Git.DeleteBranch(branch, force, remote)
	if must_force {
		return DataResponse{
			Response: "must-force",
		}
	}
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

func (a *App) AddTag(hash string, name string, message string, push bool) DataResponse {
	err := a.Git.AddTag(hash, name, message, push)
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

func (a *App) RemoveFile(file string) DataResponse {
	err := a.Git.RemoveFile(file)
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

func (a *App) DiscardChanges(file string) DataResponse {
	err := a.Git.DiscardChanges(file)
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

func (a *App) MakeCommit(subject string, body string, amend bool) DataResponse {
	err := a.Git.MakeCommit(subject, body, amend)
	return dataResponse(err, true)
}

func (a *App) MakeStash(subject string) DataResponse {
	err := a.Git.MakeStash(subject)
	return dataResponse(err, true)
}

func (a *App) CheckoutCommit(hash string) DataResponse {
	err := a.Git.CheckoutCommit(hash)
	return dataResponse(err, true)
}

func (a *App) HardReset(hash string) DataResponse {
	err := a.Git.ResetToCommit(hash, true)
	return dataResponse(err, true)
}

func (a *App) SoftReset(hash string) DataResponse {
	err := a.Git.ResetToCommit(hash, false)
	return dataResponse(err, true)
}

func (a *App) RevertCommit(hash string) DataResponse {
	err := a.Git.RevertCommit(hash)
	return dataResponse(err, true)
}

func (a *App) DropCommit(hash string) DataResponse {
	err := a.Git.DropCommit(hash)
	return dataResponse(err, true)
}

func (a *App) CherryPick(hash string, record bool, no_commit bool) DataResponse {
	err := a.Git.CherryPick(hash, record, no_commit)
	return dataResponse(err, true)
}

func (a *App) UpdateMain(branch string) DataResponse {
	r := a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo]
	r.Main = branch
	a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo] = r
	go a.saveRepoData()
	a.Git.Repo = r
	return dataResponse(nil, true)
}

func (a *App) GetInProgressCommitMessageEdit() DataResponse {
	message := git.CommitMessage{}
	msg, err := a.Git.GetInProgressCommitMessage(false)
	if err == nil {
		message = a.Git.ParseCommitMessage(msg)
	}
	return dataResponse(nil, message)
}

func (a *App) GetInProgressCommitMessageMerge() DataResponse {
	message := git.CommitMessage{}
	msg, err := a.Git.GetInProgressCommitMessage(true)
	if err == nil {
		message = a.Git.ParseCommitMessage(msg)
	}
	return dataResponse(err, message)
}

func (a *App) GetLastCommitMessage() DataResponse {
	message, err := a.Git.GetLastCommitMessage()
	return dataResponse(err, message)
}

func (a *App) GetWorkingFileParsedDiff(file string, status string, staged bool, ignore_size bool) DataResponse {
	if !ignore_size && a.fileTooLarge(file) {
		return DataResponse{
			Response: "too-large",
		}
	} else {
		diff, err := a.Git.GetWorkingFileParsedDiff(file, status, staged)
		return dataResponse(err, diff)
	}
}

func (a *App) GetCommitFileParsedDiff(hash string, file string, oldfile string, ignore_size bool) DataResponse {
	if !ignore_size && a.fileTooLarge(file) {
		return DataResponse{
			Response: "too-large",
		}
	} else {
		diff, err := a.Git.GetCommitFileParsedDiff(hash, file, oldfile)
		return dataResponse(err, diff)
	}
}

func (a *App) GetConflictParsedDiff(file string, ignore_size bool) DataResponse {
	if !ignore_size && a.fileTooLarge(file) {
		return DataResponse{
			Response: "too-large",
		}
	} else {
		diff, err := a.Git.GetConflictParsedDiff(file)
		return dataResponse(err, diff)
	}
}

func (a *App) StagePartialFile(diff git.Diff, filename string, status string) DataResponse {
	err := a.Git.StagePartial(diff, filename, status)
	return dataResponse(err, false)
}

func (a *App) UnstagePartialFile(diff git.Diff, filename string, status string) DataResponse {
	err := a.Git.UnstagePartial(diff, filename, status)
	return dataResponse(err, false)
}

func (a *App) AddRemote(name string, fetch_url string, push_url string) DataResponse {
	err := a.Git.AddRemote(name, fetch_url, push_url)
	if err == nil {
		err = a.Git.FetchRemote(name)
	}
	return dataResponse(err, false)
}

func (a *App) DeleteRemote(name string) DataResponse {
	err := a.Git.RemoveRemote(name)
	return dataResponse(err, false)
}

func (a *App) ResolveDiffConflicts(diff git.Diff) DataResponse {
	err := a.Git.ResolveDiffConflicts(diff)
	if err == nil {
		err = a.Git.StageFiles(diff.File)
	}
	return dataResponse(err, false)
}

func (a *App) GetRepoState() git.RepoState {
	if a.RepoSaveData.CurrentRepo == "" || a.Git.Repo.Directory == "" {
		return git.RepoStateNil
	}
	return a.Git.GetRepoState()
}

func (a *App) PopStash(hash string, index bool) DataResponse {
	err := a.Git.PopStash(hash, index)
	return dataResponse(err, false)
}

func (a *App) ApplyStash(stash string, index bool) DataResponse {
	err := a.Git.ApplyStash(stash, index)
	return dataResponse(err, false)
}

func (a *App) DropStash(stash string) DataResponse {
	err := a.Git.DropStash(stash)
	return dataResponse(err, false)
}

func (a *App) CreateBranchFromStash(stash string, branch_name string) DataResponse {
	err := a.Git.CreateBranchFromStash(stash, branch_name)
	return dataResponse(err, false)
}

func (a *App) GetHighlightedFile(file string) DataResponse {
	path := filepath.Join(a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo].Directory, file)
	data, err := files.HighlightFile(path)
	return dataResponse(err, data)
}

func (a *App) GetHighlightedFileRange(file string, ranges [][2]int) DataResponse {
	path := filepath.Join(a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo].Directory, file)
	data, err := files.HighlightFileSelection(path, ranges)
	return dataResponse(err, data)
}

func (a *App) GetGitConfigLocal() DataResponse {
	cfg, err := a.Git.GetConfigLocal()
	return dataResponse(err, cfg)
}

func (a *App) GetGitConfigGlobal() DataResponse {
	cfg, err := a.Git.GetConfigGlobal()
	return dataResponse(err, cfg)
}

func (a *App) GetGitConfigSystem() DataResponse {
	cfg, err := a.Git.GetConfigSystem()
	return dataResponse(err, cfg)
}

func (a *App) UpdateGitConfigUserName(list string, value string) DataResponse {
	err := a.Git.UpdateUserName(list, value)
	return dataResponse(err, true)
}

func (a *App) UpdateGitConfigUserEmail(list string, value string) DataResponse {
	err := a.Git.UpdateUserEmail(list, value)
	return dataResponse(err, true)
}

func (a *App) UpdateGitConfigUserSigningKey(list string, value string) DataResponse {
	err := a.Git.UpdateUserSigningKey(list, value)
	return dataResponse(err, true)
}

func (a *App) RebaseContinue() DataResponse {
	err := a.Git.RebaseContinue()
	return dataResponse(err, true)
}

func (a *App) RebaseAbort() DataResponse {
	err := a.Git.RebaseAbort()
	return dataResponse(err, true)
}

func (a *App) RebaseSkip() DataResponse {
	err := a.Git.RebaseSkip()
	return dataResponse(err, true)
}

func (a *App) CherryPickContinue() DataResponse {
	err := a.Git.CherryPickContinue()
	return dataResponse(err, true)
}

func (a *App) CherryPickAbort() DataResponse {
	err := a.Git.CherryPickAbort()
	return dataResponse(err, true)
}

func (a *App) CherryPickSkip() DataResponse {
	err := a.Git.CherryPickSkip()
	return dataResponse(err, true)
}
