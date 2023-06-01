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
		println(err.Error())
		return false
	}

	return r == ""
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
