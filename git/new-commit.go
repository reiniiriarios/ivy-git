package git

import "errors"

func (g *Git) MakeCommit(subject string, body string) error {
	if subject == "" {
		return errors.New("no commit subject specified")
	}

	//todo: partials
	if g.isStagedEmpty() {
		err := g.StageAll()
		if err != nil {
			return err
		}
	}

	cmd := []string{"commit", "-m", subject}
	if body != "" {
		// https://git-scm.com/docs/git-commit#Documentation/git-commit.txt--mltmsggt
		cmd = append(cmd, "-m", body)
	}
	_, err := g.RunCwd(cmd...)

	return err
}

func (g *Git) RevertCommit(hash string) error {
	if hash == "" {
		return errors.New("no commit hash specified")
	}
	_, err := g.RunCwd("revert", hash)
	return err
}

func (g *Git) CherryPick(hash string, record bool, no_commit bool) error {
	cmd := []string{"cherry-pick", "--allow-empty"}
	if record {
		cmd = append(cmd, "-x")
	}
	if no_commit {
		cmd = append(cmd, "--no-commit")
	}
	cmd = append(cmd, hash)

	_, err := g.RunCwd(cmd...)
	return err
}
