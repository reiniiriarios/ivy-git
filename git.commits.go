package main

import (
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const GIT_LOG_SEP = "-act45j3o9y78__jyo9ct-a4ojy9actyo_ct4oy9j-"

type Commit struct {
	Hash        string
	Parent      string
	AuthorName  string
	AuthorEmail string
	AuthorDate  string // todo
	Subject     string
}

type CommitResponse struct {
	Response string
	Message  string
	Commits  []Commit
}

func (a *App) GetCommits() ([]Commit, error) {
	var commits []Commit

	format := "%H" + GIT_LOG_SEP + "%P" + GIT_LOG_SEP + "%an" + GIT_LOG_SEP + "%ae" + GIT_LOG_SEP + "%ad" + GIT_LOG_SEP + "%s"
	c, err := a.GitCwd("--no-pager", "log", "--format='"+format+"'")
	if err != nil {
		return commits, err
	}

	cs := strings.Split(strings.ReplaceAll(c, "\r\n", "\n"), "\n")
	for _, cm := range cs {
		cm = strings.Trim(cm, "'")
		parts := strings.Split(cm, GIT_LOG_SEP)
		if len(parts) == 6 {
			commits = append(commits, Commit{
				Hash:        parts[0],
				Parent:      parts[1],
				AuthorName:  parts[2],
				AuthorEmail: parts[3],
				AuthorDate:  parts[4],
				Subject:     parts[5],
			})
		} else if strings.Trim(cm, " ") != "" {
			runtime.LogError(a.ctx, "unable to parse commit message")
			runtime.LogError(a.ctx, cm)
		}
	}

	return commits, nil
}

func (a *App) GetCommitsForTree() CommitResponse {
	commits, err := a.GetCommits()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return CommitResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	return CommitResponse{
		Response: "success",
		Commits:  commits,
	}
}
