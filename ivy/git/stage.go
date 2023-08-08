package git

import (
	"errors"
)

func (g *Git) isStagedEmpty() bool {
	d, err := g.run("diff", "--name-only", "--cached")
	if err != nil {
		return true
	}
	d = parseOneLine(d)
	return d == ""
}

func (g *Git) StageFiles(file ...string) error {
	cmd := append([]string{"add"}, file...)
	_, err := g.run(cmd...)
	return err
}

func (g *Git) UnstageFile(file string) error {
	if file == "" {
		return errors.New("no file specified")
	}
	_, err := g.run("reset", "--", file)
	return err
}

func (g *Git) RemoveFile(file string) error {
	if file == "" {
		return errors.New("no file specified")
	}
	_, err := g.run("rm", "--", file)
	return err
}

func (g *Git) StageAll() error {
	_, err := g.run("add", "--all")
	return err
}

func (g *Git) UnstageAll() error {
	_, err := g.run("reset")
	return err
}

func (g *Git) StagePartial(diff Diff, filename string, status string) error {
	var from string
	new_file := fileIsNew(status)
	if new_file {
		from = ""
	} else {
		from = filename
	}
	patch := diff.createPatch(from, filename, new_file)
	if patch == "" {
		return nil
	}
	_, err := g.runWithOpts([]string{"apply", "--cached", "--unidiff-zero", "--ignore-whitespace", "--whitespace=nowarn", "-"}, gitRunOpts{stdin: patch})
	return err
}

func (g *Git) UnstagePartial(diff Diff, filename string, status string) error {
	patch := diff.createDiscardPatch(filename)
	if patch == "" {
		return nil
	}
	_, err := g.runWithOpts([]string{"apply", "--cached", "--unidiff-zero", "--ignore-whitespace", "--whitespace=nowarn", "-"}, gitRunOpts{stdin: patch})
	return err
}
