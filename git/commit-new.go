package git

import "errors"

func (g *Git) MakeCommit(subject string, body string, amend bool) error {
	if subject == "" {
		return errors.New("no commit subject specified")
	}

	if g.isStagedEmpty() {
		err := g.StageAll()
		if err != nil {
			return err
		}
	}

	cmd := []string{"commit"}
	if amend {
		cmd = append(cmd, "--amend")
	}
	cmd = append(cmd, "--message", subject)
	if body != "" {
		// https://git-scm.com/docs/git-commit#Documentation/git-commit.txt--mltmsggt
		cmd = append(cmd, "--message", body)
	}
	_, err := g.run(cmd...)

	return err
}

func (g *Git) RevertCommit(hash string) error {
	if hash == "" {
		return errors.New("no commit hash specified")
	}
	_, err := g.run("revert", hash)
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

	_, err := g.run(cmd...)
	return err
}
