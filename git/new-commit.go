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
