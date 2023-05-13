package main

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

type CommitAddl struct {
	Hash               string
	Body               string
	CommitterName      string
	CommitterEmail     string
	CommitterTimestamp int64
	CommitterDatetime  string
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
	c, err := a.GitCwd("--no-pager", "log", "--format='"+format+"'", "--branches", "--tags", "--glob=refs/remotes", "HEAD")
	if err != nil {
		return commits, lookup, err
	}

	var i uint64 = 0
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
				Merge:           len(parents) > 1,
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
				runtime.LogWarning(a.ctx, "Error parsing ref: "+name)
			}
		}
	}

	return refs, nil
}

// Get all stashes from `git reflog`.
func (a *App) getStashes() []Commit {
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
	s, err := a.GitCwd("reflog", "--format='"+format+"'", "refs/stash", "--")
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return stashes
	}

	ss := a.getLines(s)
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
				runtime.LogError(a.ctx, err.Error())
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
func (a *App) getNumUncommitedChanges() int {
	c, err := a.GitCwd("status", "--untracked-files=all", "--porcelain")
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return 0
	}
	return strings.Count(c, "\n")
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

type CommitsResponse struct {
	Response string
	Message  string
	HEAD     Ref
	Commits  []Commit
	Graph    Graph
}

// FRONTEND: Get list of commits and all associated details for display.
func (a *App) GetCommitList() CommitsResponse {
	commits, lookup, HEAD, err := a.getCommits()
	if err != nil {
		return CommitsResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	// Add stashes.
	stashes := a.getStashes()
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
	u := a.getNumUncommitedChanges()
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

	vertices, lookup2 := a.getVertices(commits, HEAD)

	// Build all graph data.
	g := Graph{
		Vertices: vertices,
	}
	g.BuildPaths()

	// Add color and x-coord to commits from graph data.
	// Add whether the commit should have a label.
	for i := range commits {
		if vertices[lookup2[commits[i].Hash]].BranchId != -1 {
			bid := vertices[lookup2[commits[i].Hash]].BranchId
			if bid != -1 {
				commits[i].Color = g.Branches[bid].Color
			}
		}
		commits[i].X = vertices[lookup2[commits[i].Hash]].X
		commits[i].Labeled =
			len(commits[i].Heads) > 0 ||
				len(commits[i].Branches) > 0 ||
				len(commits[i].Tags) > 0 ||
				len(commits[i].Remotes) > 0 ||
				commits[i].Hash == HEAD.Hash ||
				commits[i].Stash
		commits[i].Id = uint64(i)
	}

	return CommitsResponse{
		Response: "success",
		HEAD:     HEAD,
		Commits:  commits,
		Graph:    g,
	}
}

type CommitResponse struct {
	Response string
	Message  string
	Commit   CommitAddl
}

// FRONTEND: Get additional commit details not listed in the table.
func (a *App) GetCommitDetails(hash string) CommitResponse {
	// Include:
	// %an - Committer Name
	// %ae - Committer Email
	// %at - Committer Time
	// %b  - Body
	// https://git-scm.com/docs/pretty-formats
	data := []string{"%cn", "%ce", "%ct", "%b"}
	format := strings.Join(data, GIT_LOG_SEP)
	c, err := a.GitCwd("--no-pager", "log", hash, "--format='"+format+"'", "--max-count=1")
	if err != nil {
		return CommitResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	c = strings.Trim(strings.Trim(strings.Trim(c, "\n"), "\r"), "'")
	parts := strings.Split(c, GIT_LOG_SEP)
	if len(parts) != len(data) {
		return CommitResponse{
			Response: "error",
			Message:  "Error fetching commit details.",
		}
	}

	// Get timestamp and formatted datetime for committer
	ts, err := strconv.ParseInt(parts[2], 10, 64)
	dt := ""
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	} else {
		dt = time.Unix(ts, 0).Format(DATE_FORMAT)
	}

	return CommitResponse{
		Response: "success",
		Commit: CommitAddl{
			Hash:               hash,
			Body:               parts[3],
			CommitterName:      parts[0],
			CommitterEmail:     parts[1],
			CommitterTimestamp: ts,
			CommitterDatetime:  dt,
		},
	}
}

type FileStat struct {
	File    string
	Name    string
	Dir     string
	Path    []string
	OldFile string
	OldName string
	OldDir  string
	Added   uint32
	Deleted uint32
	Status  string
}

type FileStatDir struct {
	Name  string
	Path  []string
	Files []FileStat
	Dirs  []FileStatDir
}

type CommitDiffSummaryResponse struct {
	Response string
	Message  string
	Files    FileStatDir
}

// FRONTEND: Get commit diff summary from diff-tree --numstat and --name-status.
func (a *App) GetCommitDiffSummary(hash string) CommitDiffSummaryResponse {
	filestats := []FileStat{}

	// Get the number of lines added and deleted from each file.
	n, err := a.GitCwd("diff-tree", "--numstat", "-r", "--root", "--find-renames", "-z", hash)
	if err != nil {
		return CommitDiffSummaryResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}
	n = strings.Trim(strings.Trim(strings.Trim(n, "\n"), "\r"), "'")
	// The -z option splits lines by NUL.
	nl := strings.Split(n, "\x00")
	// The first line is the hash, skip.
	for i := 1; i < len(nl); i++ {
		nf := strings.Fields(nl[i])
		if len(nf) == 3 {
			a, _ := strconv.ParseInt(nf[0], 10, 32)
			d, _ := strconv.ParseInt(nf[1], 10, 32)
			name := filepath.Base(nf[2])
			dir := filepath.Dir(nf[2])
			path := strings.Split(strings.ReplaceAll(dir, "\\", "/"), "/")
			filestats = append(filestats, FileStat{
				File:    nf[2],
				Name:    name,
				Dir:     dir,
				Path:    path,
				Added:   uint32(a),
				Deleted: uint32(d),
			})
		} else if len(nf) == 2 {
			// If there are two fields parsed, the next two lines are the
			// previous name and the new name.
			a, _ := strconv.ParseInt(nf[0], 10, 32)
			d, _ := strconv.ParseInt(nf[1], 10, 32)
			i++
			oldfile := nl[i]
			i++
			file := nl[i]
			name := filepath.Base(file)
			dir := filepath.Dir(file)
			path := strings.Split(strings.ReplaceAll(dir, "\\", "/"), "/")
			oldname := filepath.Base(oldfile)
			olddir := filepath.Dir(oldfile)
			filestats = append(filestats, FileStat{
				File:    file,
				Name:    name,
				Dir:     dir,
				Path:    path,
				OldFile: oldfile,
				OldName: oldname,
				OldDir:  olddir,
				Added:   uint32(a),
				Deleted: uint32(d),
			})
		}
	}

	// Get the status of each file in the commit.
	s, err := a.GitCwd("diff-tree", "--name-status", "-r", "--root", "--find-renames", "-z", hash)
	if err != nil {
		return CommitDiffSummaryResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}
	s = strings.Trim(strings.Trim(strings.Trim(s, "\n"), "\r"), "'")
	// The -z option splits lines by NUL.
	sl := strings.Split(s, "\x00")
	// The first line is the hash, skip.
	// Each line is either the status or the file. Parse two lines at a time.
	for i := 1; i < len(sl)-1; i += 2 {
		for f := range filestats {
			// Renames get three lines of data.
			if sl[i][:1] == "R" {
				if filestats[f].File == sl[i+2] {
					filestats[f].Status = "R"
					i++
					break
				}
			} else {
				if filestats[f].File == sl[i+1] {
					filestats[f].Status = sl[i]
					break
				}
			}
		}
	}

	// Parse files into directory tree.
	files := FileStatDir{}
	for _, f := range filestats {
		c := &files
		for n, p := range f.Path {
			if p != "." {
				added := false
				for j := range c.Dirs {
					if c.Dirs[j].Name == p {
						added = true
						c = &c.Dirs[j]
					}
				}
				if !added {
					c.Dirs = append(c.Dirs, FileStatDir{
						Name: p,
						Path: append(f.Path[:n], p),
					})
					c = &c.Dirs[len(c.Dirs)-1]
				}
			}
			if n == len(f.Path)-1 {
				c.Files = append(c.Files, f)
			}
		}
	}

	// Trim tree.
	trimDirs(&files)

	return CommitDiffSummaryResponse{
		Response: "success",
		Files:    files,
	}
}

// Trim dirs with no contents except one dir.
// This collapses dirs to, e.g.
//
//	foo / bar
//	  baz.go
//
// if foo doesn't have any files changed and only one subdir, bar.
func trimDirs(dir *FileStatDir) {
	if dir.Name != "" && len(dir.Dirs) == 1 && len(dir.Files) == 0 {
		dir.Dirs[0].Name = dir.Name + " / " + dir.Dirs[0].Name
		*dir = dir.Dirs[0]
		trimDirs(dir)
	} else {
		for d := range dir.Dirs {
			trimDirs(&dir.Dirs[d])
		}
	}
}
