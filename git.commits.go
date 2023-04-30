package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const GIT_LOG_SEP = "-act45j3o9y78__jyo9ct-a4ojy9actyo_ct4oy9j-"

type Commit struct {
	Hash            string
	Parents         []string
	AuthorName      string
	AuthorEmail     string
	AuthorTimestamp int64
	AuthorDatetime  string
	Subject         string
	Tags            []string
	Remotes         []string
	Heads           []string
}

type Ref struct {
	Hash string
	Name string
}

type Refs struct {
	HEAD    string
	Tags    []Ref
	Remotes []Ref
	Heads   []Ref
}

type CommitResponse struct {
	Response string
	Message  string
	HEAD     string
	Commits  []Commit
}

// Get commit details from `git log`.
func (a *App) GetCommits() ([]Commit, map[string]uint64, error) {
	var commits []Commit
	lookup := make(map[string]uint64)

	format := strings.Join([]string{"%H", "%P", "%an", "%ae", "%at", "%s"}, GIT_LOG_SEP)
	c, err := a.GitCwd("--no-pager", "log", "--format='"+format+"'")
	if err != nil {
		return commits, lookup, err
	}

	var i uint64 = 0
	cs := strings.Split(strings.ReplaceAll(c, "\r\n", "\n"), "\n")
	for _, cm := range cs {
		cm = strings.Trim(cm, "'")
		parts := strings.Split(cm, GIT_LOG_SEP)
		if len(parts) == 6 {

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
				dt = time.Unix(ts, 0).Format("Jan 1, 2006, 03:04:05 pm")
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
func (a *App) GetRefs() (Refs, error) {
	var refs Refs

	r, err := a.GitCwd("show-ref", "--dereference", "--head")
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return refs, err
	}

	rs := strings.Split(strings.ReplaceAll(r, "\r\n", "\n"), "\n")
	for _, r := range rs {
		rr := strings.Split(strings.Trim(r, "'"), " ")
		if len(rr) >= 2 {
			hash := rr[0]
			name := strings.Join(rr[1:], " ")
			if strings.HasPrefix(name, "refs/heads/") {
				refs.Heads = append(refs.Heads, Ref{
					Hash: hash,
					Name: name[11:],
				})
			} else if strings.HasPrefix(name, "refs/tags/") {
				refs.Tags = append(refs.Tags, Ref{
					Hash: hash,
					Name: name[10:],
				})
			} else if strings.HasPrefix(name, "refs/remotes/") {
				refs.Remotes = append(refs.Remotes, Ref{
					Hash: hash,
					Name: name[13:],
				})
			} else if name == "HEAD" {
				refs.HEAD = hash
			} else {
				runtime.LogWarning(a.ctx, "Error parsing ref: "+name)
			}
		}
	}

	return refs, nil
}

// Compile commits and refs for tree view.
func (a *App) GetCommitsForTree() CommitResponse {
	commits, lookup, err := a.GetCommits()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return CommitResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	refs, err := a.GetRefs()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return CommitResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	for _, h := range refs.Heads {
		if c, exists := lookup[h.Hash]; exists {
			commits[c].Heads = append(commits[c].Heads, h.Name)
		}
	}

	for _, t := range refs.Tags {
		if c, exists := lookup[t.Hash]; exists {
			commits[c].Tags = append(commits[c].Tags, t.Name)
		}
	}

	for _, r := range refs.Remotes {
		if c, exists := lookup[r.Hash]; exists {
			commits[c].Remotes = append(commits[c].Remotes, r.Name)
		}
	}

	return CommitResponse{
		Response: "success",
		HEAD:     refs.HEAD,
		Commits:  commits,
	}
}
