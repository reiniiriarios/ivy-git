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

const REF_MAX_NAME_LENGTH = 50

type Commit struct {
	Id              uint64
	Hash            string
	Parents         []string
	Children        []string
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
	HEAD            bool
	DropPossible    bool
}

type CommitDetails struct {
	Commits []Commit
	Lookup  map[string]uint64
	HEAD    Ref
}

// Get commit details from `git log`.
func (g *Git) getLog(limit uint64, offset uint64) ([]Commit, map[string]uint64, error) {
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
	count := "--max-count=" + strconv.FormatUint(limit, 10)
	skip := "--skip=" + strconv.FormatUint(offset, 10)
	c, err := g.RunCwd("--no-pager", "log", "--format="+format, count, skip, "--branches", "--tags", "--glob=refs/remotes", "HEAD")
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

type CommitsSigned map[string]string

// Get simple signature status of commit list.
func (g *Git) GetCommitsSignStatus(limit uint64, offset uint64) (CommitsSigned, error) {
	commits := make(CommitsSigned)

	// Include:
	// %G?
	//   G = good (valid)
	//   B = bad
	//   U = unknown validity
	//   X = expired
	//   Y = good signature, expired key
	//   E = missing key
	//   N = no signature
	// https://git-scm.com/docs/pretty-formats
	data := []string{"%H", "%G?"}
	format := strings.Join(data, GIT_LOG_SEP)
	count := "--max-count=" + strconv.FormatUint(limit, 10)
	skip := "--skip=" + strconv.FormatUint(offset, 10)
	c, err := g.RunCwd("--no-pager", "log", "--format="+format, count, skip, "--branches", "--tags", "--glob=refs/remotes", "HEAD")
	if err != nil {
		return commits, err
	}

	lines := parseLines(c)
	for _, commit := range lines {
		parts := strings.Split(commit, GIT_LOG_SEP)
		if len(parts) == len(data) {
			commits[parts[0]] = parts[1]
		}
	}

	return commits, nil
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
func (g *Git) getCommits(limit uint64, offset uint64) ([]Commit, map[string]uint64, Ref, error) {
	commits, lookup, err := g.getLog(limit, offset)
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
func (g *Git) getCommitList(limit uint64, offset uint64) (Ref, []Commit, error) {
	commits, lookup, HEAD, err := g.getCommits(limit, offset)
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
func (g *Git) GetCommitsAndGraph(limit uint64, offset uint64) (Ref, []Commit, Graph, error) {
	HEAD, commits, err := g.getCommitList(limit, offset)
	if err != nil {
		return HEAD, commits, Graph{}, err
	}

	// Build all graph data.
	graph := g.BuildGraph(HEAD, &commits)

	// Add color and x-coord to commits from graph data.
	// Add whether the commit should have a label.
	for i := range commits {
		if graph.Vertices[graph.CommitLookup[commits[i].Hash]].BranchId != -1 {
			bid := graph.Vertices[graph.CommitLookup[commits[i].Hash]].BranchId
			if bid != -1 {
				commits[i].Color = graph.Limbs[bid].Color
			}
		}
		commits[i].X = graph.Vertices[graph.CommitLookup[commits[i].Hash]].X
		commits[i].HEAD = commits[i].Hash == HEAD.Hash
		commits[i].Labeled =
			len(commits[i].Heads) > 0 ||
				len(commits[i].Branches) > 0 ||
				len(commits[i].Tags) > 0 ||
				len(commits[i].Remotes) > 0 ||
				commits[i].HEAD ||
				commits[i].Stash
		commits[i].Id = uint64(i)
		commits[i].DropPossible = isDroppable(uint64(i), &graph.CommitLookup, &commits)
		for child := range graph.Vertices[graph.CommitLookup[commits[i].Hash]].Children {
			commits[i].Children = append(commits[i].Children, graph.Vertices[child].Hash)
		}
	}

	return HEAD, commits, graph, nil
}

// Whether dropping the commit is (somewhat easily) possible.
func isDroppable(commit_id uint64, lookup *map[string]int64, commits *[]Commit) bool {
	return droppableTopologyPasses(commit_id, lookup, commits) == 1
}

// Search commit topology recursively.
//
//	  -1  Topology fails (fail)
//		 0  Head not found (fail)
//		 1  Head found     (pass)
func droppableTopologyPasses(commit_id uint64, lookup *map[string]int64, commits *[]Commit) int8 {
	if len((*commits)[commit_id].Parents) > 1 {
		// Merge commit
		return -1
	}
	if len((*commits)[commit_id].Children) > 1 {
		// Splits to two branches
		return -1
	}
	if len((*commits)[commit_id].Children) == 1 {
		// Search recursively through first child
		topology := droppableTopologyPasses((*commits)[(*lookup)[(*commits)[commit_id].Children[0]]].Id, lookup, commits)
		if topology != 0 {
			// Either topology has failed (-1) or HEAD was found (1)
			return topology
		}
	}
	if (*commits)[commit_id].HEAD {
		return 1
	}
	return 0
}
