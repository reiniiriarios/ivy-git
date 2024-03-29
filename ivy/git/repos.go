package git

import (
	"os"
	"path/filepath"
	"strings"
)

// This file deals with repo general git things.

type Repo struct {
	Name      string
	Directory string
	Main      string
	State     RepoState
}

func (g *Git) IsDir(directory string) bool {
	_, err := os.Stat(directory)
	return !os.IsNotExist(err)
}

func (g *Git) IsGitRepo(directory string) bool {
	if !g.IsDir(directory) {
		return false
	}
	r, err := g.runWithOpts([]string{"rev-parse"}, gitRunOpts{directory: directory})
	if err != nil {
		return false
	}

	return r == ""
}

func (g *Git) CloneRepo(url string, name string, directory string) error {
	// Directory error checking happens before this method is called.
	_, err := g.runWithOpts([]string{"clone", url, name}, gitRunOpts{directory: directory})
	return err
}

// Initialize git in a specific directory.
func (g *Git) GitInit(directory string) error {
	_, err := g.runWithOpts([]string{"init"}, gitRunOpts{directory: directory})
	return err
}

func (g *Git) HasCommits(directory string) bool {
	_, err := g.runWithOpts([]string{"rev-list", "--count", "HEAD", "--"}, gitRunOpts{directory: directory})
	if err == nil {
		return true
	}
	if errorCode(err) != NoCommitsYet && errorCode(err) != BadRevision {
		println("Error (HasCommits)", err.Error())
	}
	return false
}

// Check common names for main branch.
func (g *Git) NameOfMainBranchForRepo(repo_dir string) string {
	r, err := g.runWithOpts([]string{
		"for-each-ref",
		"--format=%(refname:short)",
		"refs/heads/main",
		"refs/heads/master",
		"refs/heads/trunk",
		"refs/heads/develop",
	}, gitRunOpts{directory: repo_dir})
	if err != nil {
		// If the previous is failing, just return something.
		return "main"
	}
	r = parseOneLine(r)

	// If we get no results, see if there are any heads. Maybe the repo is new.
	if strings.TrimSpace(r) == "" {
		all_refs, err := g.runWithOpts([]string{"for-each-ref", "--format=%(refname:short)", "refs/heads"}, gitRunOpts{directory: repo_dir})
		// If there are no heads at all, get the current branch from the HEAD symbolic-ref.
		// Ignore errors.
		if err == nil && strings.TrimSpace(all_refs) == "" {
			ref, err := g.runWithOpts([]string{"symbolic-ref", "HEAD"}, gitRunOpts{directory: repo_dir})
			// Ignore errors.
			if err == nil {
				// refs/heads/main => main
				ref = parseOneLine(ref)
				if len(ref) > 1 {
					parts := strings.Split(ref, "/")
					return parts[len(parts)-1]
				}
			}
		}
	}

	// If no newlines, there was only one result. Return it.
	if !strings.Contains(r, "\n") {
		return r
	}

	// Check the following in order.
	lines := parseLines(r)
	if stringInSlice("master", lines) {
		return "master"
	}
	if stringInSlice("main", lines) {
		return "main"
	}
	if stringInSlice("trunk", lines) {
		return "trunk"
	}
	if stringInSlice("develop", lines) {
		return "develop"
	}

	// Default to main.
	return "main"
}

// Go really just needs this on its own.
func stringInSlice(needle string, haystack []string) bool {
	for _, search := range haystack {
		if search == needle {
			return true
		}
	}
	return false
}

// Name of main branch for current repo.
func (g *Git) NameOfMainBranch() string {
	if g.Repo == (Repo{}) {
		return ""
	}
	return g.NameOfMainBranchForRepo(g.Repo.Directory)
}

type RepoState string

const (
	RepoStateNil                = ""
	RepoStateNone               = "RepoStateNone"
	RepoStateRebaseInteractive  = "RepoStateRebaseInteractive"
	RepoStateRebaseMerge        = "RepoStateRebaseMerge"
	RepoStateMerge              = "RepoStateMerge"
	RepoStateRebase             = "RepoStateRebase"
	RepoStateApply              = "RepoStateApply"
	RepoStateApplyOrRebase      = "RepoStateApplyOrRebase"
	RepoStateRevert             = "RepoStateRevert"
	RepoStateRevertSequence     = "RepoStateRevertSequence"
	RepoStateCherryPick         = "RepoStateCherryPick"
	RepoStateCherryPickSequence = "RepoStateCherryPickSequence"
	RepoStateBisect             = "RepoStateBisect"

	GitFileHead                   = "HEAD"
	GitFileOrigHead               = "ORIG_HEAD"
	GitFileFetchHead              = "FETCH_HEAD"
	GitFileMergeHead              = "MERGE_HEAD"
	GitFileRevertHead             = "REVERT_HEAD"
	GitFileCherryPickHead         = "CHERRY_PICK_HEAD"
	GitFileRebaseHead             = "REBASE_HEAD"
	GitFileBisectLog              = "BISECT_LOG"
	GitDirRebaseMerge             = "rebase-merge"
	GitFileRebaseMergeInteractive = "interactive" // in GitDirRebaseMerge
	GitFileRebaseMergeHeadName    = "head-name"   // in GitDirRebaseMerge or GitDirRebaseApply
	GitDirRebaseApply             = "rebase-apply"
	GitFileRebaseApplyRebasing    = "rebasing" // in GitDirRebaseApply
	GitFileRebaseApplyApplying    = "applying" // in GitDirRebaseApply
	GitDirSequencer               = "sequencer"
	GitFileSequencerHead          = "head"    // in GitDirSequencer
	GitFileSequencerOptions       = "options" // in GitDirSequencer
	GitFileSequencerTodo          = "todo"    // in GitDirSequencer
	GitFileStash                  = "stash"
	GitDirRefs                    = "refs"
	GitDirRefsHeads               = "heads"       // in GitDirRefs
	GitDirRefsTags                = "tags"        // in GitDirRefs
	GitDirRefsRemotes             = "remotes"     // in GitDirRefs
	GitDirRefsNotes               = "notes"       // in GitDirRefs
	GitFileRenamedRef             = "RENAMED-REF" // in GitDirRefs
	GitFileRefsHeadsMaster        = "master"      // in GitDirRefsHeads
)

func (g *Git) GetRepoState() RepoState {
	var state RepoState = RepoStateNone

	if g.gitDirHasFile(filepath.Join(GitDirRebaseMerge, GitFileRebaseMergeInteractive)) {
		state = RepoStateRebaseInteractive
	} else if g.gitDirHasFile(GitDirRebaseMerge) {
		state = RepoStateRebaseMerge
	} else if g.gitDirHasFile(filepath.Join(GitDirRebaseApply, GitFileRebaseApplyRebasing)) {
		state = RepoStateRebase
	} else if g.gitDirHasFile(filepath.Join(GitDirRebaseApply, GitFileRebaseApplyApplying)) {
		state = RepoStateApply
	} else if g.gitDirHasFile(GitDirRebaseApply) {
		state = RepoStateApplyOrRebase
	} else if g.gitDirHasFile(GitFileMergeHead) {
		state = RepoStateMerge
	} else if g.gitDirHasFile(GitFileRevertHead) {
		if g.gitDirHasFile(filepath.Join(GitDirSequencer, GitFileSequencerTodo)) {
			state = RepoStateRevertSequence
		} else {
			state = RepoStateRevert
		}
	} else if g.gitDirHasFile(GitFileCherryPickHead) {
		if g.gitDirHasFile(filepath.Join(GitDirSequencer, GitFileSequencerTodo)) {
			state = RepoStateCherryPickSequence
		} else {
			state = RepoStateCherryPick
		}
	} else if g.gitDirHasFile(GitFileBisectLog) {
		state = RepoStateBisect
	}

	g.Repo.State = state

	return state
}

func (g *Git) gitDirHasFile(file string) bool {
	_, err := os.Stat(filepath.Join(g.Repo.Directory, ".git", file))
	return !os.IsNotExist(err)
}
