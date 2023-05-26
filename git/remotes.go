package git

import (
	"errors"
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

func (g *Git) getRemoteNames() ([]string, error) {
	r, err := g.RunCwd("remote")
	if err != nil {
		return []string{}, err
	}
	lines := parseLines(r)
	remotes := []string{}
	for _, l := range lines {
		if l != "" {
			remotes = append(remotes, l)
		}
	}
	return remotes, nil
}

func (g *Git) GetRemoteForCurrentBranch() (string, error) {
	b, err := g.GetCurrentBranch()
	if err != nil {
		return "", err
	}
	r, err := g.getRemoteForBranch(b)
	if err != nil {
		return "", err
	}
	return r, nil
}

func (g *Git) getRemoteForBranch(branch string) (string, error) {
	if branch == "" {
		return "", errors.New("no branch name specified")
	}

	r, err := g.RunCwd("config", "branch."+branch+".remote")
	r = parseOneLine(r)
	// If not found configured, get first remote. This won't ordinarily happen?
	if err != nil || r == "" {
		all, err := g.getRemoteNames()
		if err != nil {
			return "", err
		}
		if len(all) < 1 {
			return "", errors.New("no remotes found")
		}
		r = all[0]
	}

	return r, nil
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
				var fetch, push bool
				if d[2] == "(fetch)" {
					fetch = true
					push = false
				} else if d[2] == "(push)" {
					fetch = false
					push = true
				}

				url, err := url.Parse(d[1])
				if err != nil {
					return remotes, err
				}

				site := getSiteName(url.Hostname())

				userRepo := strings.Trim(url.Path, "/")
				if len(userRepo) > 4 && userRepo[len(userRepo)-4:] == ".git" {
					userRepo = userRepo[:len(userRepo)-4]
				}

				user := ""
				repoName := ""
				if strings.Count(userRepo, "/") == 1 {
					urn := strings.Split(userRepo, "/")
					user = urn[0]
					repoName = urn[1]
				}

				var ahead, behind uint32 = 0, 0
				currentBranch, err := g.GetCurrentBranch()
				if err == nil {
					ahead, behind, err = g.getAheadBehind(currentBranch, d[0])
					if err != nil {
						numCommits, err := g.NumCommitsOnBranch(currentBranch)
						if err == nil {
							ahead = uint32(numCommits)
						}
					}
				}

				rmap[d[0]] = len(remotes)
				remotes = append(remotes, Remote{
					Name:     d[0],
					Url:      d[1],
					Fetch:    fetch,
					Push:     push,
					Site:     site,
					Repo:     userRepo,
					User:     user,
					RepoName: repoName,
					Ahead:    ahead,
					Behind:   behind,
				})
			}
		}
	}

	return remotes, nil
}

// Get the name of the main branch for a specific remote.
func (g *Git) getMainBranchForRemote(remote string) string {
	if remote == "" {
		return ""
	}

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
	if remote == "" {
		return errors.New("no remote name specified")
	}
	_, err := g.RunCwd("fetch", remote, "--prune")
	return err
}
