package git

import (
	"strconv"
	"strings"
	"time"
)

// Get all stashes from `git reflog`.
func (g *Git) getStashes() []Commit {
	var stashes []Commit

	// Include:
	// %H  - hash
	// %P  - parent hash
	// %gd - refname
	// %an - Author Name
	// %ae - Author Email
	// %at - Author Time
	// %s  - Subject
	// https://git-scm.com/docs/pretty-formats
	data := []string{"%H", "%P", "%gd", "%an", "%ae", "%at", "%s"}
	format := strings.Join(data, GIT_LOG_SEP)
	// todo: consider replacing with `git stash list` to prevent "bad revision" errors
	s, err := g.run("reflog", "--format="+format, "refs/stash", "--")
	if err != nil {
		// if no stashes:
		// fatal: bad revision 'refs/stash'
		return stashes
	}

	ss := parseLines(s)
	for _, st := range ss {
		parts := strings.Split(st, GIT_LOG_SEP)
		if len(parts) == len(data) {

			// Get parents.
			// Stashes may have two or three parents. In the following diagrams,
			//   I = index at time of stash
			//   S = stash
			//   U = uncommitted files at time of stash (-u flag)
			//
			//          .----S
			//         /    /
			//   -----H----I
			//
			//          .----S----.
			//         /    /    /
			//   -----H----I    U
			//
			parents := []string{}
			if parts[1] != "" {
				parents = strings.Split(parts[1], " ")
			}

			// Get timestamp and formatted datetime for author
			ts, err := strconv.ParseInt(parts[5], 10, 64)
			dt := ""
			if err == nil {
				dt = time.Unix(ts, 0).Format(DATE_FORMAT)
			}

			stashes = append(stashes, Commit{
				Hash:            parts[0],
				Parents:         parents,
				RefName:         parts[2],
				AuthorName:      parts[3],
				AuthorEmail:     parts[4],
				AuthorTimestamp: ts,
				AuthorDatetime:  dt,
				Subject:         parts[6],
				Stash:           true,
			})
		}
	}

	return stashes
}

func (g *Git) MakeStash(subject string) error {
	if g.isStagedEmpty() {
		err := g.StageAll()
		if err != nil {
			return err
		}
	}

	_, err := g.run("stash", "push", "--staged", "--message", subject)

	return err
}

// Pop a stash, optionally reinstating index.
// Stash should be a ref in the form of `stash@{<revision>}`.
func (g *Git) PopStash(stash string, index bool) error {
	cmd := []string{"stash", "pop"}
	if index {
		cmd = append(cmd, "--index")
	}
	cmd = append(cmd, stash)
	_, err := g.run(cmd...)
	return err
}

// Apply a stash, optionally reinstating index.
// Stash should be a ref in the form of `stash@{<revision>}`.
func (g *Git) ApplyStash(stash string, index bool) error {
	cmd := []string{"stash", "apply"}
	if index {
		cmd = append(cmd, "--index")
	}
	cmd = append(cmd, stash)
	_, err := g.run(cmd...)
	return err
}

// Drop a stash.
// Stash should be a ref in the form of `stash@{<revision>}`.
func (g *Git) DropStash(stash string) error {
	_, err := g.run("stash", "drop", stash)
	return err
}

// Create a branch from a stash.
// Stash should be a ref in the form of `stash@{<revision>}`.
func (g *Git) CreateBranchFromStash(stash string, branch_name string) error {
	_, err := g.run("stash", "branch", branch_name, stash)
	return err
}
