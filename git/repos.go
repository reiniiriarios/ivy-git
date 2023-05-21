package git

import "os"

type Repo struct {
	Name      string
	Directory string
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
