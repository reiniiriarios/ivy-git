package git

func (g *Git) LsFiles() ([]string, error) {
	f, err := g.run("ls-files")
	if err != nil {
		return []string{}, err
	}
	files := parseLines(f)
	return files, nil
}

func (g *Git) LsTreeBranch(branch string) ([]string, error) {
	f, err := g.run("ls-tree", "-r", "--name-only", branch)
	if err != nil {
		return []string{}, err
	}
	files := parseLines(f)
	return files, nil
}

func (g *Git) ShowFileOnBranch(branch string, file string) (string, error) {
	f, err := g.run("--no-pager", "show", branch+":"+file)
	if err != nil {
		return "", err
	}
	return f, nil
}
