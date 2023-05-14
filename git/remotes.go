package git

import (
	"net/url"
	"strings"
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

func (g *Git) GetRemotes() ([]Remote, error) {
	remotes := []Remote{}
	rmap := make(map[string]int)

	rs, err := g.RunCwd("remote", "-v")
	if err != nil {
		return remotes, err
	}
	rl := parseLines(rs)

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
					return remotes, err
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

				var ad, bd uint32 = 0, 0
				currentBranch, err := g.GetCurrentBranch()
				if err != nil {
					ad, bd, err = g.getAheadBehind(currentBranch, d[0])
					if err != nil {
						numCommits, err := g.getNumCommitsOnBranch(currentBranch)
						if err == nil {
							ad = numCommits
						}
					}
				}

				rmap[d[0]] = len(remotes)
				remotes = append(remotes, Remote{
					Name:     d[0],
					Url:      d[1],
					Fetch:    f,
					Push:     p,
					Site:     s,
					Repo:     ur,
					User:     u,
					RepoName: rn,
					Ahead:    ad,
					Behind:   bd,
				})
			}
		}
	}

	return remotes, nil
}

// Get the name of the main branch for a specific remote.
func (g *Git) getMainBranchForRemote(remote string) string {
	ls, err := g.RunCwd("ls-remote", remote)
	if err != nil {
		println(err.Error())
		return ""
	}
	rs := parseLines(ls)

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

func (g *Git) FetchRemote(remote string) error {
	_, err := g.RunCwd("fetch", remote, "--prune")
	return err
}
