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
		// This is sorted enough, no need to go by deletions or else.
		return c[i].Insertions > c[j].Insertions
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

const CONTRIBUTORS_LOG_LIMIT = 1000

func (g *Git) AddContributorsSince(contributors Contributors, start_hash string) (Contributors, string, error) {
	if g.Repo.Main == "" {
		return contributors, start_hash, errors.New("no main branch set for current repo")
	}
	final_hash := g.lastCommitOnMain()

	// If this hasn't been run before, start at the beginning.
	boundary := false
	if start_hash == "" {
		boundary = true
		start_hash = g.getInitialCommit()
	}
	// Last hash processed.
	current_hash := ""

	contributors_map := contributors.toMap()

	data := []string{"%H"}
	if GIT_RESPECT_MAILMAP {
		data = append(data, "%aN", "%aE")
	} else {
		data = append(data, "%an", "%ae")
	}
	format := strings.Join(data, GIT_LOG_SEP) + "\n"

	// The -z option does not include the --shortstat data and, in fact, places a \x00 between
	// the commit data specified in --format and the --shortstat data, grouping the data
	// incorrectly. As there are newlines between each commit's data as well as within its data,
	// we add a random string by which we separate and parse commit data.
	split := "_-ghowuigah9384h2g0h8222-_"
	// Loop n results at a time for larger repos.
	for i := 0; ; /* breaks below */ i += CONTRIBUTORS_LOG_LIMIT {
		cmd := []string{
			"--no-pager",
			"log",
			start_hash + ".." + final_hash,
			"--format=" + split + format,
			"--shortstat",
			"--max-count=" + strconv.Itoa(CONTRIBUTORS_LOG_LIMIT),
			"--skip=" + strconv.Itoa(i),
		}
		// If starting on the initial commit, we should include that commit as well.
		if boundary {
			cmd = append(cmd, "--boundary")
		}
		// if bad revision, this ensures correct error is thrown
		cmd = append(cmd, "--")

		c, err := g.run(cmd...)
		if err != nil {
			return contributors, start_hash, err
		}

		// Remove extra whitespace.
		c = strings.ReplaceAll(c, "\r\n", "\n")
		multinewline := regexp.MustCompile(`(\r\n?|\n){2,}`)
		c = multinewline.ReplaceAllString(c, "\n")
		c = strings.TrimSpace(c)
		if c == "" {
			break
		}

		r_ins := regexp.MustCompile(`([0-9]+) insertion`)
		r_del := regexp.MustCompile(`([0-9]+) deletion`)

		// Loop lines
		lines := strings.Split(c, split)
		if len(lines) == 0 {
			break
		}
		for _, line := range lines {
			line = multinewline.ReplaceAllString(line, "\n")

			// Split commit data and shortstat.
			parts := strings.Split(strings.TrimSpace(line), "\n")
			if len(parts) < 2 && strings.TrimSpace(parts[0]) == "" {
				continue
			}

			// Split commit data.
			parts_a := strings.Split(parts[0], GIT_LOG_SEP)
			if len(parts_a) != len(data) {
				// Malformed data. This may occur on extraneous lines.
				continue
			}
			// Keep this updated to the most recent one processed.
			current_hash = parts_a[0]
			// Relevant data
			name := parts_a[1]
			email := parts_a[2]

			// If contributor not listed yet, add to list.
			if _, exists := contributors_map[email]; !exists {
				contributors_map[email] = Contributor{
					Name:  name,
					Email: email,
				}
			}

			contributor := contributors_map[email]
			contributor.Commits++

			// Empty commits won't have this data.
			if len(parts) > 1 {
				// Add to insertions and deletions counts.
				match_ins := r_ins.FindAllStringSubmatch(parts[1], 1)
				if len(match_ins) > 0 && len(match_ins[0]) > 0 {
					u, err := strconv.ParseUint(match_ins[0][1], 0, 64)
					// ignore errors
					if err == nil {
						contributor.Insertions += u
					}
				}
				match_del := r_del.FindAllStringSubmatch(parts[1], 1)
				if len(match_del) > 0 && len(match_del[0]) > 0 {
					u, err := strconv.ParseUint(match_del[0][1], 0, 64)
					// ignore errors
					if err == nil {
						contributor.Deletions += u
					}
				}
			}

			contributors_map[email] = contributor
		}

		// If the last processed commit is the last commit on main, we're done.
		if current_hash == start_hash {
			break
		}
	}

	return contributors_map.toArray(), g.lastCommitOnMain(), nil
}

// Get number of commits behind the tip of main a commit is.
func (g *Git) CommitBehindMain(hash string) (uint64, error) {
	if g.Repo.Main == "" {
		return 0, errors.New("no main branch set for current repo")
	}
	return g.CommitBehindCount(hash, g.lastCommitOnMain())
}

// Get number of commits behind one commit is from another.
func (g *Git) CommitBehindCount(older_hash string, newer_hash string) (uint64, error) {
	if older_hash == newer_hash {
		return 0, nil
	}
	c, err := g.run("rev-list", older_hash+".."+newer_hash, "--count", "--")
	if err != nil {
		return 0, err
	}
	c = parseOneLine(c)
	n, err := strconv.ParseUint(c, 10, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}
