package git

import (
	"errors"
	"strings"
)

func (g *Git) PushTag(name string) error {
	if name == "" {
		return errors.New("no tag name specified")
	}

	remote, err := g.getRemoteForCurrentBranch()
	if err != nil {
		return err
	}
	_, err = g.RunCwd("push", remote, name)
	return err
}

func (g *Git) DeleteTag(name string) error {
	if name == "" {
		return errors.New("no tag name specified")
	}

	remote, err := g.getRemoteForCurrentBranch()
	if err != nil {
		return err
	}
	_, err = g.RunCwd("tag", "-d", name)
	if err != nil {
		return err
	}
	_, err = g.RunCwd("push", remote, ":refs/tags/"+name)
	return err
}

func (g *Git) AddTag(hash string, name string, annotated bool, message string, push bool) error {
	if hash == "" {
		return errors.New("no commit hash specified")
	}
	if name == "" {
		return errors.New("no tag name specified")
	}

	var err error
	if annotated {
		if message == "" {
			_, err = g.RunCwd("tag", "-a", name, hash)
		} else {
			_, err = g.RunCwd("tag", "-a", name, hash, "-m", message)
		}
	} else {
		_, err = g.RunCwd("tag", name, hash)
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
	t, err := g.RunCwd("ls-remote", "--tags", remote)
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
