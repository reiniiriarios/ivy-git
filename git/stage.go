package git

import (
	"errors"
)

func (g *Git) isStagedEmpty() bool {
	d, err := g.RunCwd("diff", "--name-only", "--cached")
	if err != nil {
		println(err.Error())
		return true
	}
	d = parseOneLine(d)
	return d == ""
}

func (g *Git) StageFiles(file ...string) error {
	cmd := append([]string{"add"}, file...)
	_, err := g.RunCwd(cmd...)
	return err
}

func (g *Git) UnstageFile(file string) error {
	if file == "" {
		return errors.New("no file specified")
	}
	_, err := g.RunCwd("reset", "--", file)
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

func (g *Git) StagePartial(diff Diff, newFile bool) error {
	patch := diff.createPatch(newFile)
	_, err := g.RunCwdStdin([]string{"apply", "--cached", "--unidiff-zero", "--whitespace=nowarn", "-"}, patch)
	return err
}

func (g *Git) UnstagePartial(diff Diff) error {
	patch := diff.createDiscardPatch()
	_, err := g.RunCwdStdin([]string{"apply", "--unidiff-zero", "--whitespace=nowarn", "-"}, patch)
	return err
}
