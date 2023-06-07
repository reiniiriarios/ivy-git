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
	s, err := g.RunCwd("reflog", "--format="+format, "refs/stash", "--")
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
