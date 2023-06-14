package git

import (
	"errors"
	"strings"
)

func (g *Git) NumTags() uint64 {
	b, err := g.run("tag")
	if err != nil {
		return 0
	}
	lines := parseLines(b)
	return uint64(len(lines))
}

func (g *Git) PushTag(name string) error {
	if name == "" {
		return errors.New("no tag name specified")
	}

	remote, err := g.GetRemoteForCurrentBranch()
	if err != nil {
		return err
	}
	_, err = g.run("push", remote, name)
	return err
}

func (g *Git) DeleteTag(name string) error {
	if name == "" {
		return errors.New("no tag name specified")
	}

	remote, err := g.GetRemoteForCurrentBranch()
	if err != nil {
		return err
	}
	_, err = g.run("tag", "-d", name)
	if err != nil {
		return err
	}
	_, err = g.run("push", remote, ":refs/tags/"+name)
	return err
}

func (g *Git) AddTag(hash string, name string, message string, push bool) error {
	if hash == "" {
		return errors.New("no commit hash specified")
	}
	if name == "" {
		return errors.New("no tag name specified")
	}

	var err error
	if message != "" {
		_, err = g.run("tag", "-a", name, hash, "-m", message)
	} else {
		_, err = g.run("tag", name, hash)
	}
	if err != nil {
		return err
	}

	if push {
		err = g.PushTag(name)
		return err
	}

	return nil
}

func (g *Git) getRemoteTags(remote string) ([]string, error) {
	if remote == "" {
		return []string{}, errors.New("no remote name specified")
	}

	tags := []string{}
	t, err := g.run("ls-remote", "--tags", remote)
	if err != nil {
		return tags, err
	}
	lines := parseLines(t)
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 2 {
			tags = append(tags, parts[1])
		}
	}
	return tags, nil
}
