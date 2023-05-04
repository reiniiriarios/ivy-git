package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const GIT_LOG_SEP = "-act45j3o9y78__jyo9ct-a4ojy9actyo_ct4oy9j-"
const UNCOMMITED_HASH = "#"
const DATE_FORMAT = "Jan 1, 2006, 03:04:05 pm"

type Commit struct {
	Hash            string
	Parents         []string
	AuthorName      string
	AuthorEmail     string
	AuthorTimestamp int64
	AuthorDatetime  string
	Subject         string
	Branches        []Ref
	Tags            []Ref
	Remotes         []Ref
	Heads           []Ref
	Stash           Stash
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

type CommitResponse struct {
	Response string
	Message  string
	HEAD     Ref
	Commits  []Commit
}

// Get commit details from `git log`.
func (a *App) getLog() ([]Commit, map[string]uint64, error) {
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
	c, err := a.GitCwd("--no-pager", "log", "--format='"+format+"'", "--branches", "--tags", "--reflog", "--glob=refs/remotes", "HEAD")
	if err != nil {
		return commits, lookup, err
	}

	// Start counting commits at 1, 0 will be uncommited changes (if any)
	var i uint64 = 1
	cs := a.getLines(c)
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
				runtime.LogError(a.ctx, err.Error())
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
			})
			lookup[parts[0]] = i
			i++

		} else if strings.Trim(cm, " ") != "" {
			runtime.LogError(a.ctx, "unable to parse commit message")
			runtime.LogError(a.ctx, cm)
		}
	}

	return commits, lookup, nil
}

// Get ref details from `git show-ref`.
func (a *App) getRefs() (Refs, error) {
	var refs Refs

	r, err := a.GitCwd("show-ref", "--dereference", "--head")
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return refs, err
	}

	// For the purposes of displaying a coherent tree,
	// we're denoting the following:
	// - refs/heads/[...]                 = branches
	// - refs/tags/[...]                  = tags
	// - HEAD and refs/remotes/[...]/HEAD = heads
	// - refs/remotes/[...]/[...]         = remotes
	rs := a.getLines(r)
	for _, r := range rs {
		rr := strings.Split(r, " ")
		if len(rr) >= 2 {
			hash := rr[0]
			name := strings.Join(rr[1:], " ")
			if strings.HasPrefix(name, "refs/heads/") {
				refs.Branches = append(refs.Branches, Ref{
					Hash: hash,
					Name: name[11:],
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
			} else {
				runtime.LogWarning(a.ctx, "Error parsing ref: "+name)
			}
		}
	}

	return refs, nil
}

// Compile commits and refs for tree view.
func (a *App) getCommits() ([]Commit, map[string]uint64, Ref, error) {
	commits, lookup, err := a.getLog()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return commits, lookup, Ref{}, err
	}

	refs, err := a.getRefs()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
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

func (a *App) GetCommitList() CommitResponse {
	commits, lookup, HEAD, err := a.getCommits()
	if err != nil {
		return CommitResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	u := a.getNumUncommitedChanges()

	if u > 0 {
		t := time.Now()
		// The first index of commits is reserved for uncommited changes.
		commits[0] = Commit{
			Hash:            UNCOMMITED_HASH,
			Parents:         []string{HEAD.Hash},
			AuthorName:      "*",
			AuthorEmail:     "",
			AuthorTimestamp: t.Unix(),
			AuthorDatetime:  t.Format(DATE_FORMAT),
			Subject:         "Uncommited Changes (" + strconv.Itoa(u) + ")",
		}
	}

	g := Graph{
		Vertices: a.getVertices(commits, lookup, HEAD),
	}
	g.BuildPaths()

	return CommitResponse{
		Response: "success",
		HEAD:     HEAD,
		Commits:  commits,
	}
}

func (a *App) getNumUncommitedChanges() int {
	c, err := a.GitCwd("status", "--untracked-files=all", "--porcelain")
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return 0
	}
	return strings.Count(c, "\n")
}

type Stash struct {
	Hash            string
	Parent          string
	RefName         string
	AuthorName      string
	AuthorEmail     string
	AuthorTimestamp int64
	AuthorDatetime  string
	Subject         string
}

func (a *App) getStashes() []Stash {
	var stashes []Stash

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
	s, err := a.GitCwd("reflog", "--format='"+format+"'", "--glob=refs/stash")
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return stashes
	}

	ss := a.getLines(s)
	for _, st := range ss {
		parts := strings.Split(st, GIT_LOG_SEP)
		if len(parts) == len(data) {
			// Get timestamp and formatted datetime for author
			ts, err := strconv.ParseInt(parts[5], 10, 64)
			dt := ""
			if err != nil {
				runtime.LogError(a.ctx, err.Error())
			} else {
				dt = time.Unix(ts, 0).Format(DATE_FORMAT)
			}

			stashes = append(stashes, Stash{
				Hash:            parts[0],
				Parent:          parts[1],
				RefName:         parts[2],
				AuthorName:      parts[3],
				AuthorEmail:     parts[4],
				AuthorTimestamp: ts,
				AuthorDatetime:  dt,
				Subject:         parts[7],
			})
		}
	}

	return stashes
}