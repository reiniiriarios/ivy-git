package git

import (
	"errors"
)

// Get a simple list of untracked files, no other data.
func (g *Git) GetUntrackedFiles() (string, error) {
	l, err := g.RunCwd("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return "", err
	}
	return l, nil
}

func (g *Git) GetUncommittedDiff() (string, error) {
	diff, err := g.RunCwd("--no-pager", "diff", "HEAD^", "--")
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

	diff, err := g.RunCwd("--no-pager", "diff", remote+"/"+branch)
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) GetDiffStaged() (string, error) {
	diff, err := g.RunCwd("--no-pager", "diff", "--staged")
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) findMergeBase(hash1 string, hash2 string) (string, error) {
	if hash1 == "" || hash2 == "" {
		return "", errors.New("no commit hash specified")
	}

	b, err := g.RunCwd("merge-base", hash1, hash2)
	if err != nil {
		return "", err
	}
	b = parseOneLine(b)
	return b, nil
}

func (g *Git) GetWorkingFileParsedDiff(file string, status string, staged bool) (Diff, error) {
	raw, err := g.getWorkingFileDiff(file, status, staged)
	if err != nil {
		return Diff{}, err
	}
	diff := Diff{
		Raw: raw,
	}
	err = diff.parse()
	if err != nil {
		return Diff{}, err
	}
	return diff, nil
}

func (g *Git) getWorkingFileDiff(file string, status string, staged bool) (string, error) {
	cmd := []string{"diff", "-w", "--no-ext-diff", "--patch-with-raw", "-z", "--no-color"}

	var d string

	if status == FileUntracked {
		// --no-index emulates exit codes from `diff`, will return 1 when changes found
		// https://github.com/git/git/blob/1f66975deb8402131fbf7c14330d0c7cdebaeaa2/diff-no-index.c#L300
		cmd = append(cmd, "--no-index", "--", "/dev/null", file)
		d, _ = g.RunCwdNoError(cmd...)
	} else {
		if staged {
			cmd = append(cmd, "--cached", "--", file)
		} else {
			cmd = append(cmd, "--", file)
		}
		var err error
		d, err = g.RunCwd(cmd...)
		if err != nil {
			return "", err
		}
	}

	return d, nil
}

func (g *Git) GetConflictParsedDiff(file string) (Diff, error) {
	base, err := g.getDiffBase(file)
	if err != nil {
		return Diff{}, err
	}
	diff := Diff{
		Raw: base,
	}
	err = diff.parseConflicts()
	if err != nil {
		return Diff{}, err
	}
	return diff, nil
}

func (g *Git) getDiffBase(file string) (string, error) {
	d, err := g.RunCwd("diff", "--base", file)
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
		Raw: raw,
	}
	err = diff.parse()
	if err != nil {
		return Diff{}, err
	}
	return diff, nil
}

func (g *Git) getCommitFileDiff(hash string, file string, oldfile string) (string, error) {
	cmd := []string{"log", hash, "-w", "-m", "-1", "--first-parent", "--patch-with-raw", "-z", "--no-color", "--", file}
	if oldfile != "" {
		cmd = append(cmd, oldfile)
	}
	d, err := g.RunCwd(cmd...)
	if err != nil {
		return "", err
	}
	return d, nil
}
