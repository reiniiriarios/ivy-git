package git

import (
	"errors"
	"net/url"
	"regexp"
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
	r, err := g.run("remote")
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
		// Ignore errors here, there may not be a currently selected branch.
		return "", nil
	}
	if b == "" {
		return "", nil
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

	r, err := g.run("config", "branch."+branch+".remote")
	r = parseOneLine(r)
	// If not found configured, get first remote. This won't ordinarily happen?
	if err != nil || r == "" {
		all, err := g.getRemoteNames()
		if err != nil {
			return "", err
		}
		if len(all) < 1 {
			// No remotes found.
			return "", nil
		}
		r = all[0]
	}

	return r, nil
}

var gitUrlSSHRegex = regexp.MustCompile(`^(?:[^@]+@)?([^:]+):(.+)$`)

func (g *Git) GetRemotes() ([]Remote, error) {
	remotes := []Remote{}
	rmap := make(map[string]int)

	rs, err := g.run("remote", "-v")
	if err != nil {
		return remotes, err
	}
	rl := parseLines(rs)

	for _, r := range rl {
		d := strings.Fields(r)
		if len(d) == 3 {
			name := d[0]
			uri := d[1]
			direction := d[2]

			if i, exists := rmap[name]; exists {
				if direction == "(fetch)" {
					remotes[i].Fetch = true
				} else if direction == "(push)" {
					remotes[i].Push = true
				}
			} else {
				var fetch, push bool
				if direction == "(fetch)" {
					fetch = true
					push = false
				} else if direction == "(push)" {
					fetch = false
					push = true
				}

				site := ""
				userRepo := ""

				url, err := url.Parse(uri)
				// Errors possible on valid URLs for SSH, e.g.:
				// "git@github.com:user/repo.git": first path segment in URL cannot contain colon
				if err == nil {
					if url.Hostname() != "" {
						site = getSiteName(url.Hostname())
					} else if url.Scheme != "" {
						// In the case of SSH urls, the Scheme is probably the "site", e.g.
						// sitealias:/home/user/repo.git
						site = url.Scheme
					}

					userRepo = strings.Trim(url.Path, "/")
					if len(userRepo) > 4 && userRepo[len(userRepo)-4:] == ".git" {
						userRepo = userRepo[:len(userRepo)-4]
					}
				} else {
					// Try to parse SSH urls and return site and path, e.g.
					// git@github.com:user/repo.git => github.com, user/repo.git
					matches := gitUrlSSHRegex.FindAllStringSubmatch(uri, -1)
					if len(matches) > 0 && len(matches[0]) > 1 {
						site = matches[0][1]
						userRepo = matches[0][2]
					} else {
						// If all fails, set the userRepo as the uri itself, leaving site blank.
						userRepo = uri
					}
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

	ls, err := g.run("ls-remote", remote)
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
	_, err := g.run("fetch", remote, "--prune")
	return err
}

func (g *Git) AddRemote(name string, fetch_url string, push_url string) error {
	if name == "" {
		return errors.New("no remote name specified")
	}
	if fetch_url == "" {
		return errors.New("no fetch url specified")
	}

	_, err := g.run("remote", "add", name, fetch_url)
	if err != nil {
		return err
	}

	if push_url != "" {
		_, err = g.run("remote", "set-url", name, "--push", push_url)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Git) RemoveRemote(name string) error {
	if name == "" {
		return errors.New("no remote name specified")
	}
	_, err := g.run("remote", "rm", name)
	return err
}

func (g *Git) isOnGitHub() bool {
	r, err := g.run("remote", "-v")
	if err != nil {
		// ignore errors
		return false
	}
	return strings.Contains(strings.ToLower(r), "github")
}
