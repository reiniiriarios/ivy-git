package git

import (
	"strconv"
	"strings"
	"time"
)

// Match to graph.ts.
const UNCOMMITED_HASH = "#"

// Displayed in commit list.
const DATE_FORMAT = "Jan 2, 2006, 03:04:05 pm"

const GIT_LOG_SEP = "-act45j3o9y78__jyo9ct-a4ojy9actyo_ct4oy9j-"

type Commit struct {
	Id              uint64
	Hash            string
	Parents         []string
	RefName         string
	AuthorName      string
	AuthorEmail     string
	AuthorTimestamp int64
	AuthorDatetime  string
	Subject         string
	Branches        []Ref
	Tags            []Ref
	Remotes         []Ref
	Heads           []Ref
	Stash           bool
	Merge           bool
	Labeled         bool
	Color           uint16
	X               uint16
}

type Ref struct {
	Hash      string
	Name      string
	ShortName string
	Annotated bool
}

type Refs struct {
	HEAD     Ref
	Branches []Ref
	Tags     []Ref
	Remotes  []Ref
	Heads    []Ref
}

type CommitDetails struct {
	Commits []Commit
	Lookup  map[string]uint64
	HEAD    Ref
}

// Get commit details from `git log`.
func (g *Git) getLog() ([]Commit, map[string]uint64, error) {
	var commits []Commit
	lookup := make(map[string]uint64)

	// Include:
	// %H  - hash
	// %P  - parent hash
	// %an - Author Name
	// %ae - Author Email
	// %at - Author Time
	// %s  - Subject
	// https://git-scm.com/docs/pretty-formats
	data := []string{"%H", "%P", "%an", "%ae", "%at", "%s"}
	format := strings.Join(data, GIT_LOG_SEP)
	// Include:
	// - branches
	// - tags
	// - commits mentioned by reflogs
	// - all remotes
	// - HEAD
	c, err := g.RunCwd("--no-pager", "log", "--format='"+format+"'", "--branches", "--tags", "--glob=refs/remotes", "HEAD")
	if err != nil {
		return commits, lookup, err
	}

	var i uint64 = 0
	cs := parseLines(c)
	for _, cm := range cs {
		parts := strings.Split(cm, GIT_LOG_SEP)
		if len(parts) == len(data) {

			// Get parents
			parents := []string{}
			if parts[1] != "" {
				parents = strings.Split(parts[1], " ")
			}

			// Get timestamp and formatted datetime for author
			ts, err := strconv.ParseInt(parts[4], 10, 64)
			dt := ""
			if err != nil {
				println(err.Error())
			} else {
				dt = time.Unix(ts, 0).Format(DATE_FORMAT)
			}

			// Index by SHA
			commits = append(commits, Commit{
				Hash:            parts[0],
				Parents:         parents,
				AuthorName:      parts[2],
				AuthorEmail:     parts[3],
				AuthorTimestamp: ts,
				AuthorDatetime:  dt,
				Subject:         parts[5],
				Merge:           len(parents) > 1,
			})
			lookup[parts[0]] = i
			i++

		} else if strings.Trim(cm, " ") != "" {
			println("unable to parse commit message:", cm)
		}
	}

	return commits, lookup, nil
}

// Get ref details from `git show-ref`.
func (g *Git) getRefs() (Refs, error) {
	var refs Refs

	r, err := g.RunCwd("show-ref", "--dereference", "--head")
	if err != nil {
		println(err.Error())
		return refs, err
	}

	// For the purposes of displaying a coherent tree,
	// we're denoting the following:
	// - refs/heads/[...]                 = branches
	// - refs/tags/[...]                  = tags
	// - HEAD and refs/remotes/[...]/HEAD = heads
	// - refs/remotes/[...]/[...]         = remotes
	rs := parseLines(r)
	for _, r := range rs {
		rr := strings.Split(r, " ")
		if len(rr) >= 2 {
			hash := rr[0]
			name := strings.Join(rr[1:], " ")
			if strings.HasPrefix(name, "refs/heads/") {
				n := name[11:]
				s := n
				rrr := strings.Split(n, "/")
				if len(rrr) >= 2 {
					s = rrr[0]
				}
				refs.Branches = append(refs.Branches, Ref{
					Hash:      hash,
					Name:      n,
					ShortName: s,
				})
			} else if strings.HasPrefix(name, "refs/tags/") {
				annotated := strings.HasSuffix(name, "^{}")
				if annotated {
					name = name[10 : len(name)-3]
				} else {
					name = name[10:]
				}
				refs.Tags = append(refs.Tags, Ref{
					Hash:      hash,
					Name:      name,
					Annotated: annotated,
				})
			} else if strings.HasPrefix(name, "refs/remotes/") {
				n := name[13:]
				s := n
				rrr := strings.Split(n, "/")
				if len(rrr) >= 2 {
					s = rrr[0]
				}
				ref := Ref{
					Hash:      hash,
					Name:      n,
					ShortName: s,
				}
				if name[len(name)-4:] == "HEAD" {
					refs.Heads = append(refs.Remotes, ref)
				} else {
					refs.Remotes = append(refs.Remotes, ref)
				}
			} else if name == "HEAD" {
				refs.HEAD = Ref{
					Hash: hash,
					Name: name,
				}
			} else if !strings.HasPrefix(name, "refs/stash") {
				// Ignore stash, but anything else log a warning.
				println("Error parsing ref:", name)
			}
		}
	}

	return refs, nil
}

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
	// TODO: consider replacing with `git stash list` to prevent "bad revision" errors
	s, err := g.RunCwd("reflog", "--format='"+format+"'", "refs/stash", "--")
	if err != nil {
		println(err.Error())
		return stashes
	}

	ss := parseLines(s)
	for _, st := range ss {
		parts := strings.Split(st, GIT_LOG_SEP)
		if len(parts) == len(data) {

			// Get parents
			parents := []string{}
			if parts[1] != "" {
				parents = strings.Split(parts[1], " ")
				// Only keep the first parent for stashes.
				parents = parents[0:1]
			}

			// Get timestamp and formatted datetime for author
			ts, err := strconv.ParseInt(parts[5], 10, 64)
			dt := ""
			if err != nil {
				println(err.Error())
			} else {
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

// Get the number of changed files that are uncommitted.
func (g *Git) getNumUncommitedChanges() int {
	c, err := g.RunCwd("status", "--untracked-files=all", "--porcelain")
	if err != nil {
		println(err.Error())
		return 0
	}
	return strings.Count(c, "\n")
}

// Compile commits and refs for tree view.
func (g *Git) getCommits() ([]Commit, map[string]uint64, Ref, error) {
	commits, lookup, err := g.getLog()
	if err != nil {
		println(err.Error())
		return commits, lookup, Ref{}, err
	}

	refs, err := g.getRefs()
	if err != nil {
		println(err.Error())
		return commits, lookup, Ref{}, err
	}

	for _, b := range refs.Branches {
		if c, exists := lookup[b.Hash]; exists {
			commits[c].Branches = append(commits[c].Branches, b)
		}
	}

	for _, h := range refs.Heads {
		if c, exists := lookup[h.Hash]; exists {
			commits[c].Heads = append(commits[c].Heads, h)
		}
	}

	for _, t := range refs.Tags {
		if c, exists := lookup[t.Hash]; exists {
			commits[c].Tags = append(commits[c].Tags, t)
		}
	}

	for _, r := range refs.Remotes {
		if c, exists := lookup[r.Hash]; exists {
			commits[c].Remotes = append(commits[c].Remotes, r)
		}
	}

	return commits, lookup, refs.HEAD, nil
}

// Get list of commits and all associated details for display.
// Returns HEAD, []Commit
func (g *Git) getCommitList() (Ref, []Commit, error) {
	commits, lookup, HEAD, err := g.getCommits()
	if err != nil {
		return Ref{}, nil, err
	}

	// Add stashes.
	stashes := g.getStashes()
	for _, s := range stashes {
		// Confirm the stash hash doesn't match a commit.
		if _, exists := lookup[s.Hash]; !exists {
			// If the stash parent hash matches a commit, add to the list.
			for i := range commits {
				if commits[i].Hash == s.Parents[0] {
					if len(commits) == i {
						commits = append(commits, s)
					} else {
						commits = append(commits[:i+1], commits[i:]...)
						commits[i] = s
					}
					break
				}
			}
		}
	}

	// Add uncommitted changes.
	u := g.getNumUncommitedChanges()
	if u > 0 {
		t := time.Now()
		commits = append([]Commit{{
			Hash:            UNCOMMITED_HASH,
			Parents:         []string{HEAD.Hash},
			AuthorName:      "*",
			AuthorEmail:     "",
			AuthorTimestamp: t.Unix(),
			AuthorDatetime:  "*",
			Subject:         "Uncommited Changes (" + strconv.Itoa(u) + ")",
		}}, commits...)
	}

	return HEAD, commits, nil
}

// Get HEAD, Commits List, and Graph
func (g *Git) GetCommitsAndGraph() (Ref, []Commit, Graph, error) {
	HEAD, commits, err := g.getCommitList()
	if err != nil {
		return HEAD, commits, Graph{}, err
	}

	// Build all graph data.
	graph := g.BuildGraph(HEAD, commits)

	// Add color and x-coord to commits from graph data.
	// Add whether the commit should have a label.
	for i := range commits {
		if graph.Vertices[graph.CommitLookup[commits[i].Hash]].BranchId != -1 {
			bid := graph.Vertices[graph.CommitLookup[commits[i].Hash]].BranchId
			if bid != -1 {
				commits[i].Color = graph.Branches[bid].Color
			}
		}
		commits[i].X = graph.Vertices[graph.CommitLookup[commits[i].Hash]].X
		commits[i].Labeled =
			len(commits[i].Heads) > 0 ||
				len(commits[i].Branches) > 0 ||
				len(commits[i].Tags) > 0 ||
				len(commits[i].Remotes) > 0 ||
				commits[i].Hash == HEAD.Hash ||
				commits[i].Stash
		commits[i].Id = uint64(i)
	}

	return HEAD, commits, graph, nil
}