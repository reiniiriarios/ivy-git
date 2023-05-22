package git

import "path/filepath"

func (g *Git) StageFiles(file ...string) error {
	for i := range file {
		file[i] = filepath.Join(g.Repo.Directory, file[i])
	}
	cmd := append([]string{"add"}, file...)
	_, err := g.RunCwd(cmd...)
	return err
}

func (g *Git) UnstageFile(file string) error {
	_, err := g.RunCwd("reset", filepath.Join(g.Repo.Directory, file))
	return err
}

func (g *Git) StageAll() error {
	_, err := g.RunCwd("add", "--all")
	return err
}

func (g *Git) UnstageAll() error {
	_, err := g.RunCwd("reset")
	return err
}
