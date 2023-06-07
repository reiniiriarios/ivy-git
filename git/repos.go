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
}

func (g *Git) IsDir(directory string) bool {
	_, err := os.Stat(directory)
	return !os.IsNotExist(err)
}

func (g *Git) IsGitRepo(directory string) bool {
	if !g.IsDir(directory) {
		return false
	}
	r, err := g.Run(directory, "rev-parse")
	if err != nil {
		return false
	}

	return r == ""
}

func (g *Git) CloneRepo(url string, name string, directory string) error {
	// Directory error checking happens before this method is called.
	_, err := g.Run(directory, "clone", url, name)
	return err
}

func (g *Git) HasCommits(directory string) bool {
	_, err := g.Run(directory, "rev-list", "--count", "HEAD", "--")
	if err == nil {
		return true
	}
	if errorCode(err) != NoCommitsYet && errorCode(err) != BadRevision {
		println(err.Error())
	}
	return false
}

// Check common names for main branch.
func (g *Git) NameOfMainBranchForRepo(repo_dir string) string {
	r, err := g.Run("-C", repo_dir, "for-each-ref", "--format=%(refname:short)", "refs/heads/main", "refs/heads/master", "refs/heads/trunk")
	if err != nil {
		// Screw it, return something.
		return "main"
	}
	r = parseOneLine(r)
	if !strings.Contains(r, "\n") {
		return r
	}
	// More than one result.
	if strings.Contains(r, "master") {
		return "master"
	}
	// Default to main.
	return "main"
}

// Name of main branch for current repo.
func (g *Git) NameOfMainBranch() string {
	if g.Repo == (Repo{}) {
		return ""
	}
	return g.NameOfMainBranchForRepo(g.Repo.Directory)
}

func (g *Git) LsFiles() ([]string, error) {
	f, err := g.RunCwd("ls-files")
	if err != nil {
		return []string{}, err
	}
	files := parseLines(f)
	for i := range files {
		files[i] = filepath.Join(g.Repo.Directory, files[i])
	}
	return files, nil
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
	GitFileBisectLog              = "BISECT_LOG"
	GitDirRebaseMerge             = "rebase-merge"
	GitFileRebaseMergeInteractive = "interactive" // in GitDirRebaseMerge
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
	if g.gitDirHasFile(filepath.Join(GitDirRebaseMerge, GitFileRebaseMergeInteractive)) {
		return RepoStateRebaseInteractive
	}
	if g.gitDirHasFile(GitDirRebaseMerge) {
		return RepoStateRebaseMerge
	}
	if g.gitDirHasFile(filepath.Join(GitDirRebaseApply, GitFileRebaseApplyRebasing)) {
		return RepoStateRebase
	}
	if g.gitDirHasFile(filepath.Join(GitDirRebaseApply, GitFileRebaseApplyApplying)) {
		return RepoStateApply
	}
	if g.gitDirHasFile(GitDirRebaseApply) {
		return RepoStateApplyOrRebase
	}
	if g.gitDirHasFile(GitFileMergeHead) {
		return RepoStateMerge
	}
	if g.gitDirHasFile(GitFileRevertHead) {
		if g.gitDirHasFile(filepath.Join(GitDirSequencer, GitFileSequencerTodo)) {
			return RepoStateRevertSequence
		}
		return RepoStateRevert
	}
	if g.gitDirHasFile(GitFileCherryPickHead) {
		if g.gitDirHasFile(filepath.Join(GitDirSequencer, GitFileSequencerTodo)) {
			return RepoStateCherryPickSequence
		}
		return RepoStateCherryPick
	}
	if g.gitDirHasFile(GitFileBisectLog) {
		return RepoStateBisect
	}

	return RepoStateNone
}

func (g *Git) gitDirHasFile(file string) bool {
	_, err := os.Stat(filepath.Join(g.Repo.Directory, ".git", file))
	return !os.IsNotExist(err)
}
