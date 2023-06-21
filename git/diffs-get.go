package git

import (
	"errors"
	"ivy-git/files"
)

// Get a simple list of untracked files, no other data.
func (g *Git) GetUntrackedFiles() (string, error) {
	l, err := g.run("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return "", err
	}
	return l, nil
}

func (g *Git) GetUncommittedDiff() (string, error) {
	diff, err := g.run("--no-pager", "diff", "HEAD^", "--")
	if err != nil {
		if errorCode(err) == BadRevision {
			return "", nil
		}
		return "", err
	}
	return diff, nil
}

func (g *Git) GetDiffRemoteCurrent() (string, error) {
	branch, err := g.GetCurrentBranch()
	if err != nil {
		return "", err
	}
	if branch == "" {
		return "", nil
	}
	remote, err := g.getBranchRemote(branch, true)
	if err != nil {
		return "", err
	}
	if remote == "" {
		return "", nil
	}
	diff, err := g.GetDiffRemote(remote, branch)
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) GetDiffRemote(remote string, branch string) (string, error) {
	if branch == "" {
		return "", errors.New("no branch name specified")
	}
	if remote == "" {
		return "", errors.New("no remote name specified")
	}

	diff, err := g.run("--no-pager", "diff", remote+"/"+branch)
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) GetDiffStaged() (string, error) {
	diff, err := g.run("--no-pager", "diff", "--staged")
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) findMergeBase(hash1 string, hash2 string) (string, error) {
	if hash1 == "" || hash2 == "" {
		return "", errors.New("no commit hash specified")
	}

	b, err := g.run("merge-base", hash1, hash2)
	if err != nil {
		return "", err
	}
	b = parseOneLine(b)
	return b, nil
}

func (g *Git) GetWorkingFileParsedDiff(file string, status string, staged bool) (Diff, error) {
	raw, flags, err := g.getWorkingFileDiff(file, status, staged)
	if err != nil {
		return Diff{}, err
	}
	diff := Diff{
		File:  file,
		raw:   raw,
		Flags: flags,
	}
	err = diff.parse()
	if err != nil {
		return Diff{}, err
	}
	err = g.highlightDiff(&diff)
	if err != nil {
		return Diff{}, err
	}
	return diff, nil
}

func (g *Git) getWorkingFileDiff(file string, status string, staged bool) (string, []string, error) {
	cmd := []string{"diff", "-w", "--no-ext-diff", "--patch", "-z", "--no-color"}

	var d string
	flags := []string{}
	var err error

	// UNTRACKED FILES
	if status == FileUntracked {
		cmd = append(cmd, "--no-index", "--", "/dev/null", file)
		d, err = g.runWithOpts(cmd, gitRunOpts{always_return_stderr: true})
		if errorCode(err) == ReplaceLineEndings {
			flags = append(flags, "replace-line-endings")
		}
		// --no-index emulates exit codes from `diff`, will return 1 when changes found
		// https://github.com/git/git/blob/1f66975deb8402131fbf7c14330d0c7cdebaeaa2/diff-no-index.c#L300
		if err != nil && errorCode(err) != ExitStatus1 {
			return "", flags, err
		}
	} else {
		if staged {
			// STAGED TRACKED FILES
			cmd = append(cmd, "--cached", "--", file)
		} else {
			// UNSTAGED TRACKED FILES
			cmd = append(cmd, "--", file)
		}
		d, err = g.runWithOpts(cmd, gitRunOpts{always_return_stderr: true})
		if errorCode(err) == ReplaceLineEndings {
			flags = append(flags, "replace-line-endings")
		} else if err != nil {
			return "", flags, err
		}
	}

	return d, flags, nil
}

func (g *Git) GetConflictParsedDiff(file string) (Diff, error) {
	x, y, err := g.getFileStatus(file)
	if err != nil {
		return Diff{}, err
	}

	var raw string
	if (x == FileAdded || x == FileUntracked) && (y == FileAdded || y == FileUntracked) {
		// If both sides are new, the conflict should be in both Ours and Theirs diffs.
		raw, err = g.getDiffOurs(file)
	} else {
		raw, err = g.getDiffBase(file)
	}
	if err != nil {
		return Diff{}, err
	}

	diff := Diff{
		File: file,
		raw:  raw,
	}
	err = diff.parseConflicts()
	if err != nil {
		return Diff{}, err
	}
	err = g.highlightDiff(&diff)
	if err != nil {
		return Diff{}, err
	}
	return diff, nil
}

func (g *Git) getDiffBase(file string) (string, error) {
	d, err := g.run("diff", "--base", "--", file)
	if err != nil {
		return "", err
	}
	return d, nil
}

func (g *Git) getDiffOurs(file string) (string, error) {
	d, err := g.run("diff", "--ours", "--", file)
	if err != nil {
		return "", err
	}
	return d, nil
}

func (g *Git) GetCommitFileParsedDiff(hash string, file string, oldfile string) (Diff, error) {
	raw, err := g.getCommitFileDiff(hash, file, oldfile)
	if err != nil {
		return Diff{}, err
	}
	diff := Diff{
		File: file,
		raw:  raw,
	}
	err = diff.parse()
	if err != nil {
		return Diff{}, err
	}
	err = g.highlightDiff(&diff)
	if err != nil {
		return Diff{}, err
	}
	return diff, nil
}

func (g *Git) getCommitFileDiff(hash string, file string, oldfile string) (string, error) {
	cmd := []string{"log", hash, "-w", "-m", "-1", "--first-parent", "--patch", "-z", "--no-color", "--", file}
	if oldfile != "" {
		cmd = append(cmd, oldfile)
	}
	d, err := g.run(cmd...)
	if err != nil {
		return "", err
	}
	return d, nil
}

func (g *Git) highlightDiff(diff *Diff) error {
	highlight, err := files.HighlightDiff(diff.File, diff.raw)
	if err != nil {
		return err
	}
	diff.Lang = highlight.Lang
	diff.Highlight = highlight.Highlight

	return nil
}
