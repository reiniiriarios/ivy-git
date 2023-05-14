package main

import (
	"net/url"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Remote struct {
	Name       string
	Url        string
	Fetch      bool
	Push       bool
	Site       string
	Repo       string
	User       string
	RepoName   string
	MainBranch string
	Ahead      uint32
	Behind     uint32
	LastUpdate int64
}

type RemotesResponse struct {
	Response string
	Message  string
	Remotes  []Remote
}

func (a *App) GetRemotes() RemotesResponse {
	remotes := []Remote{}
	rmap := make(map[string]int)

	rs, err := a.GitCwd("remote", "-v")
	if err != nil {
		return RemotesResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}
	rl := a.getLines(rs)

	for _, r := range rl {
		d := strings.Fields(r)
		if len(d) == 3 {
			if i, exists := rmap[d[0]]; exists {
				if d[2] == "(fetch)" {
					remotes[i].Fetch = true
				} else if d[2] == "(push)" {
					remotes[i].Push = true
				}
			} else {
				var f, p bool
				if d[2] == "(fetch)" {
					f = true
					p = false
				} else if d[2] == "(push)" {
					f = false
					p = true
				}

				url, err := url.Parse(d[1])
				if err != nil {
					return RemotesResponse{
						Response: "error",
						Message:  err.Error(),
					}
				}

				s := url.Hostname()
				if s == "github.com" {
					s = "GitHub"
				} else if s == "bitbucket.org" {
					s = "Bitbucket"
				} else if s == "gitlab.com" {
					s = "GitLab"
				} else if s == "dev.azure.com" {
					s = "Azure"
				}

				ur := url.Path
				if ur[:1] == "/" {
					ur = ur[1:]
				}
				if len(ur) > 4 && ur[len(ur)-4:] == ".git" {
					ur = ur[:len(ur)-4]
				}

				u := ""
				rn := ""
				if strings.Count(ur, "/") == 1 {
					urn := strings.Split(ur, "/")
					u = urn[0]
					rn = urn[1]
				}

				m := a.getMainBranchForRemote(d[0])
				var ad, bd uint32 = 0, 0
				if a.BranchExists(m) {
					ad, bd, _ = a.getAheadBehind(m, d[0])
				}

				rmap[d[0]] = len(remotes)
				remotes = append(remotes, Remote{
					Name:       d[0],
					Url:        d[1],
					Fetch:      f,
					Push:       p,
					Site:       s,
					Repo:       ur,
					User:       u,
					RepoName:   rn,
					MainBranch: m,
					Ahead:      ad,
					Behind:     bd,
				})
			}
		}
	}

	return RemotesResponse{
		Response: "success",
		Remotes:  remotes,
	}
}

// Get the name of the main branch for a specific remote.
func (a *App) getMainBranchForRemote(remote string) string {
	ls, err := a.GitCwd("ls-remote", remote)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return ""
	}
	rs := a.getLines(ls)

	head := ""
	maybe_main := ""
	refs := make(map[string]string)
	for _, l := range rs {
		r := strings.Fields(l)
		if len(r) == 2 {
			if r[1] == "HEAD" {
				head = r[0]
			} else {
				refs[r[0]] = r[1]
				if len(r) >= 4 && r[1][len(r)-4:] == "main" {
					maybe_main = "main"
				} else if len(r) >= 6 && r[1][len(r)-6:] == "master" {
					maybe_main = "master"
				}
			}
		}
	}

	if head != "" {
		if ref, exists := refs[head]; exists {
			if strings.Contains(ref, "/") {
				refp := strings.Split(ref, "/")
				return refp[len(refp)-1]
			}
			return ref
		}
	}
	return maybe_main
}

func (a *App) FetchRemote(remote string) GenericResponse {
	_, err := a.GitCwd("fetch", remote, "--prune")
	if err != nil {
		return GenericResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}
	return GenericResponse{
		Response: "success",
	}
}
