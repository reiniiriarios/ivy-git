package git

import (
	"errors"
)

// Switch branch on currently selected repo.
func (g *Git) SwitchBranch(branch string, remote string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	var err error
	if remote == "" {
		_, err = g.RunCwd("checkout", branch)
	} else if g.BranchExists(branch) {
		err = g.PullBranch(branch, true)
		if err != nil {
			return err
		}
		_, err = g.RunCwd("checkout branch")
	} else {
		_, err = g.RunCwd("checkout", "-b", branch, remote+"/"+branch)
	}
	return err
}

func (g *Git) PushBranch(branch string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}

	remote, err := g.getBranchRemote(branch, false)
	set_upstream := err != nil || remote == ""
	if set_upstream {
		remote = g.getRemoteFallback()
	}

	err = g.PushRemoteBranch(remote, branch, set_upstream)
	return err
}

func (g *Git) PushRemoteBranch(remote string, branch string, set_upstream bool) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	if remote == "" {
		return errors.New("no remote name specified")
	}

	var err error
	if set_upstream {
		_, err = g.RunCwd("push", "--set-upstream", remote, branch)
	} else {
		_, err = g.RunCwd("push", remote, branch+":"+branch)
	}
	return err
}

func (g *Git) PullBranch(branch string, rebase bool) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}

	remote, err := g.getBranchRemote(branch, true)
	if err != nil {
		return err
	}
	err = g.PullRemoteBranch(remote, branch, rebase)
	return err
}

func (g *Git) PullRemoteBranch(remote string, branch string, rebase bool) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	if remote == "" {
		return errors.New("no remote name specified")
	}

	var err error
	if rebase {
		_, err = g.RunCwd("pull", remote, branch+":"+branch, "--rebase")
	} else {
		_, err = g.RunCwd("pull", remote, branch+":"+branch)
	}
	return err
}

func (g *Git) ResetBranchToRemote(branch string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}

	remote, err := g.getBranchRemote(branch, true)
	if err != nil {
		return err
	}
	err = g.fetchBranchRemote(branch, remote)
	if err != nil {
		return err
	}

	current_branch, err := g.GetCurrentBranch()
	if err != nil {
		return err
	}
	err = g.SwitchBranch(branch, "")
	if err != nil {
		return err
	}

	_, err = g.RunCwd("reset", "--hard", remote+"/"+branch)

	g.SwitchBranch(current_branch, "")

	return err
}

// Delete a branch.
func (g *Git) DeleteBranch(branch string, force bool, delete_on_remotes bool) (bool, error) {
	if branch == "" {
		return false, errors.New("no branch name specified")
	}

	delete := "-d"
	if force {
		delete = "-D"
	}

	_, err := g.RunCwd("branch", delete, branch)
	if err != nil {
		println(err.Error())
		println(errorCode(err))
		if errorCode(err) == MustForceDeleteBranch {
			return true, err
		}
		return false, err
	}

	if delete_on_remotes {
		remotes, err := g.getRemoteNames()
		if err != nil {
			return false, err
		}
		for _, remote := range remotes {
			if g.branchExistsOnRemote(branch, remote) {
				_, err := g.RunCwd("push", delete, remote, branch)
				if err != nil {
					if errorCode(err) == MustForceDeleteBranch {
						return true, err
					}
					return false, err
				}
			}
		}
	}

	return false, nil
}

// Delete a remote branch.
func (g *Git) DeleteRemoteBranch(branch string, remote string, force bool) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	if remote == "" {
		return errors.New("no remote name specified")
	}

	delete := "-d"
	if force {
		delete = "-D"
	}
	_, err := g.RunCwd("push", delete, remote, branch)
	return err
}

// Rename a branch locally and on all remotes.
func (g *Git) RenameBranch(branch string, new_name string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	if new_name == "" {
		return errors.New("no new branch name specified")
	}

	_, err := g.RunCwd("branch", "-m", branch, new_name)
	if err != nil {
		return err
	}

	remotes, err := g.getRemoteNames()
	if err != nil {
		return err
	}
	for _, remote := range remotes {
		if g.branchExistsOnRemote(branch, remote) {
			_, err := g.RunCwd("push", remote, ":"+branch, new_name)
			if err != nil {
				return err
			}
			_, err = g.RunCwd("push", "--set-upstream", remote, new_name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Git) RebaseOnBranch(branch string) error {
	if branch == "" {
		return errors.New("no branch name specified")
	}
	_, err := g.RunCwd("rebase", branch)
	return err
}

func (g *Git) CreateBranch(name string, at_hash string, checkout bool) error {
	if name == "" {
		return errors.New("no branch name specified")
	}

	cmd := []string{}
	if checkout {
		cmd = append(cmd, "checkout", "-b")
	} else {
		cmd = append(cmd, "branch")
	}
	cmd = append(cmd, name)
	if at_hash != "" {
		cmd = append(cmd, at_hash)
	}
	_, err := g.RunCwd(cmd...)

	return err
}
