package git

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Contributor struct {
	Name       string
	Email      string
	Commits    uint64
	Insertions uint64
	Deletions  uint64
}

type Contributors []Contributor

type ContributorsMap map[string]Contributor

func (c Contributors) Len() int {
	return len(c)
}

func (c Contributors) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Contributors) Less(i, j int) bool {
	if c[i].Commits == c[j].Commits {
		return c[i].Insertions < c[j].Insertions
	}
	return c[i].Commits > c[j].Commits
}

func (m *ContributorsMap) toArray() Contributors {
	c := make(Contributors, 0, len(*m))
	for _, v := range *m {
		c = append(c, v)
	}
	return c
}

func (c *Contributors) toMap() ContributorsMap {
	m := make(ContributorsMap)
	for _, v := range *c {
		m[v.Email] = v
	}
	return m
}

func (g *Git) AddContributorsSince(contributors Contributors, hash string) (Contributors, string, error) {
	if g.Repo.Main == "" {
		return contributors, hash, errors.New("no main branch set for current repo")
	}
	last_commit_on_main := g.lastCommitOnMain()

	// If this hasn't been run before, start at the beginning.
	if hash == "" {
		hash = g.getInitialCommit()
	}

	data := []string{"%an", "%ae"}
	format := strings.Join(data, GIT_LOG_SEP)
	cmd := []string{"--no-pager", "log", hash + ".." + last_commit_on_main, "--reverse", "--format=" + format, "--shortstat", "-z"}
	if hash != "" {
		cmd = append(cmd, hash)
	}
	c, err := g.run(cmd...)
	if err != nil {
		return contributors, hash, err
	}

	contributors_map := contributors.toMap()

	// Remove extra whitespace after null separator and before shortstat.
	r := regexp.MustCompile("\x00[\r\n ]+")
	c = r.ReplaceAllString(c, "\x00")

	r_ins := regexp.MustCompile(`([0-9]+) insertion`)
	r_del := regexp.MustCompile(`([0-9]+) deletion`)

	// Loop lines
	lines := parseLines(c)
	for _, line := range lines {

		// Split commit data and shortstat.
		parts := strings.Split(line, "\x00")
		if len(parts) != 2 {
			continue
		}

		// Split commit data.
		parts_a := strings.Split(parts[0], GIT_LOG_SEP)
		if len(parts_a) != len(data) {
			continue
		}
		name := parts_a[0]
		email := parts_a[1]

		// If contributor not listed yet, add to list.
		if _, exists := contributors_map[email]; !exists {
			contributors_map[email] = Contributor{
				Name:  name,
				Email: email,
			}
		}

		contributor := contributors_map[email]
		contributor.Commits++

		// Add to insertions and deletions counts.
		match_ins := r_ins.FindAllStringSubmatch(parts[1], 1)
		if len(match_ins) > 0 && len(match_ins[0]) > 0 {
			u, err := strconv.ParseUint(match_ins[0][1], 0, 64)
			// ignore errors
			if err == nil {
				contributor.Insertions += u
			} else {
				println(err.Error())
			}
		}
		match_del := r_del.FindAllStringSubmatch(parts[1], 1)
		if len(match_del) > 0 && len(match_del[0]) > 0 {
			u, err := strconv.ParseUint(match_del[0][1], 0, 64)
			// ignore errors
			if err == nil {
				contributor.Deletions += u
			} else {
				println(err.Error())
			}
		}

		contributors_map[email] = contributor
	}

	return contributors_map.toArray(), last_commit_on_main, nil
}
