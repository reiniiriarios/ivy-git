package git

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Get a simple list of untracked files, no other data.
func (g *Git) GetUntrackedFiles() (string, error) {
	l, err := g.RunCwd("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return "", err
	}
	return l, nil
}

func (g *Git) GetUncommittedDiff() (string, error) {
	diff, err := g.RunCwd("--no-pager", "diff", "HEAD^")
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) GetDiffRemoteCurrent() (string, error) {
	branch, err := g.GetCurrentBranch()
	if err != nil {
		return "", err
	}
	remote, err := g.getBranchRemote(branch, true)
	if err != nil {
		return "", err
	}
	diff, err := g.GetDiffRemote(remote, branch)
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) GetDiffRemote(remote string, branch string) (string, error) {
	if branch == "" {
		return "", errors.New("no branch name specified")
	}
	if remote == "" {
		return "", errors.New("no remote name specified")
	}

	diff, err := g.RunCwd("--no-pager", "diff", remote+"/"+branch)
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) GetDiffStaged() (string, error) {
	diff, err := g.RunCwd("--no-pager", "diff", "--staged")
	if err != nil {
		return "", err
	}
	return diff, nil
}

func (g *Git) findMergeBase(hash1 string, hash2 string) (string, error) {
	if hash1 == "" || hash2 == "" {
		return "", errors.New("no commit hash specified")
	}

	b, err := g.RunCwd("merge-base", hash1, hash2)
	if err != nil {
		return "", err
	}
	b = parseOneLine(b)
	return b, nil
}

func (g *Git) GetWorkingFileParsedDiff(file string, status string, staged bool) (Diff, error) {
	raw, err := g.getWorkingFileDiff(file, status, staged)
	if err != nil {
		return Diff{}, err
	}
	diff := Diff{
		Raw: raw,
	}
	err = diff.parse()
	if err != nil {
		return Diff{}, err
	}
	return diff, nil
}

func (g *Git) getWorkingFileDiff(file string, status string, staged bool) (string, error) {
	cmd := []string{"diff", "-w", "--no-ext-diff", "--patch-with-raw", "-z", "--no-color"}

	var d string

	if status == FileUntracked {
		// --no-index emulates exit codes from `diff`, will return 1 when changes found
		// https://github.com/git/git/blob/1f66975deb8402131fbf7c14330d0c7cdebaeaa2/diff-no-index.c#L300
		cmd = append(cmd, "--no-index", "--", "/dev/null", file)
		d, _ = g.RunCwdNoError(cmd...)
	} else {
		if staged {
			cmd = append(cmd, "--cached", "--", file)
		} else {
			cmd = append(cmd, "--", file)
		}
		var err error
		d, err = g.RunCwd(cmd...)
		if err != nil {
			return "", err
		}
	}

	return d, nil
}

const (
	DiffAddLine     = "DiffAddLine"
	DiffDeleteLine  = "DiffDeleteLine"
	DiffContextLine = "DiffContextLine"
)

type Diff struct {
	Raw             string
	Hunks           []DiffHunk
	Binary          bool
	SelectableLines uint64
	SelectedLines   uint64
}

type DiffHunk struct {
	Header   string
	StartOld int64
	EndOld   int64
	StartNew int64
	EndNew   int64
	Heading  string
	Lines    []DiffLine
}

type DiffLine struct {
	Line      string
	Type      string
	RawLineNo int64
	OldLineNo int64
	NewLineNo int64
	NoNewline bool
	MiniHunk  int64
	Selected  bool
}

const HiddenBidiCharsRegex = "/[\u202A-\u202E]|[\u2066-\u2069]/"

func (d *Diff) parse() error {
	lines := parseLines(d.Raw)
	var ln int
	for ln = 0; ln < len(lines); ln++ {
		if strings.HasPrefix(lines[ln], "Binary files ") && strings.HasSuffix(lines[ln], "differ") {
			d.Binary = true
			break
		}
		if strings.HasPrefix(lines[ln], "+++") {
			d.Binary = false
			break
		}
	}
	if ln == len(lines)-1 {
		// Couldn't parse header
		return nil
	}
	if d.Binary {
		// No parsing binary files
		return nil
	}

	// Continue ln where we left off with the last loop
	current_hunk := -1
	var before, after int64
	var mini_hunk int64 = 0
	var mini_hunk_editable bool = false
	for ; ln < len(lines); ln++ {
		// Parse hunk heading
		if strings.HasPrefix(lines[ln], "@@") {
			hunk, err := parseHunkHeading(lines[ln])
			if err != nil {
				return err
			}
			d.Hunks = append(d.Hunks, hunk)
			current_hunk++
			before = d.Hunks[current_hunk].StartOld
			after = d.Hunks[current_hunk].StartNew
			continue
		}
		// Keep looking until we've found at least the first hunk
		if current_hunk == -1 {
			continue
		}

		// Parse other lines
		prefix := ""
		if len(lines[ln]) > 0 {
			prefix = lines[ln][:1]
		}
		switch prefix {
		case "\\":
			if len(lines[ln]) < 12 {
				return errors.New("malformed diff, ending not recognized")
			}
			if len(d.Hunks[current_hunk].Lines) < 1 {
				return errors.New("malformed diff, invalid length")
			}
			// Last line didn't have a newline, don't add this line, and we're done
			d.Hunks[current_hunk].Lines[len(d.Hunks[current_hunk].Lines)-1].NoNewline = true
			return nil
		case "+":
			d.SelectableLines++
			d.SelectedLines++
			after++
			mini_hunk_editable = true
			d.Hunks[current_hunk].Lines = append(d.Hunks[current_hunk].Lines, DiffLine{
				Line:      lines[ln][1:],
				Type:      DiffAddLine,
				RawLineNo: int64(ln),
				OldLineNo: -1,
				NewLineNo: after,
				MiniHunk:  mini_hunk,
				Selected:  true,
			})
		case "-":
			d.SelectableLines++
			d.SelectedLines++
			before++
			mini_hunk_editable = true
			d.Hunks[current_hunk].Lines = append(d.Hunks[current_hunk].Lines, DiffLine{
				Line:      lines[ln][1:],
				Type:      DiffDeleteLine,
				RawLineNo: int64(ln),
				OldLineNo: before,
				NewLineNo: -1,
				MiniHunk:  mini_hunk,
				Selected:  true,
			})
		case " ":
			before++
			after++
			// If last line was editable and this one isn't, the mini hunk has ended.
			if mini_hunk_editable {
				mini_hunk++
			}
			mini_hunk_editable = false
			d.Hunks[current_hunk].Lines = append(d.Hunks[current_hunk].Lines, DiffLine{
				Line:      lines[ln][1:],
				Type:      DiffContextLine,
				RawLineNo: int64(ln),
				OldLineNo: before,
				NewLineNo: after,
				MiniHunk:  -1,
			})
		default:
			println(lines[ln])
			return errors.New("malformed diff, unrecognized line type")
		}
	}
	return nil
}

func parseHunkHeading(line string) (DiffHunk, error) {
	r := regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)
	matches := r.FindStringSubmatch(line)
	if len(matches) < 5 {
		return DiffHunk{}, errors.New("malformed diff header")
	}

	start_old, _ := strconv.ParseInt(matches[1], 10, 64)
	start_new, _ := strconv.ParseInt(matches[3], 10, 64)
	var end_old, end_new int64
	if matches[2] != "" {
		end_old, _ = strconv.ParseInt(matches[2], 10, 64)
	} else {
		end_old = 1
	}
	if matches[4] != "" {
		end_new, _ = strconv.ParseInt(matches[4], 10, 64)
	} else {
		end_old = 1
	}

	heading := strings.Join(strings.Split(line, "@@")[2:], "@@")

	return DiffHunk{
		Header:   line,
		Heading:  heading,
		StartOld: start_old,
		EndOld:   end_old,
		StartNew: start_new,
		EndNew:   end_new,
	}, nil
}

// Format patch header.
//
//	--- a/from
//	+++ b/to
func patchHeader(from string, to string) string {
	if from == "" {
		from = "/dev/null"
	} else {
		from = "a/" + from
	}

	if to == "" {
		to = "/dev/null"
	} else {
		to = "b/" + to
	}

	return fmt.Sprintf("--- %s\n+++ %s\n", from, to)
}

// Create a patch to stage a file or commit from staging.
// new_file should be true for file status == New || Untracked
func (d *Diff) createPatch(from string, to string, new_file bool) string {
	patch := ""

	for _, hunk := range d.Hunks {
		hunkBuffer := ""
		oldCount := 0
		newCount := 0
		changes := false

		for _, line := range hunk.Lines {
			if line.Type == DiffContextLine {
				hunkBuffer += " " + line.Line + "\n"
				oldCount++
				newCount++
			} else if line.Selected {
				changes = true
				if line.Type == DiffAddLine {
					hunkBuffer += "+" + line.Line + "\n"
					newCount++
				} else if line.Type == DiffDeleteLine {
					hunkBuffer += "-" + line.Line + "\n"
					oldCount++
				}
			} else {
				// Unselected lines should be ignored in new files.
				if new_file {
					continue
				}

				// Unselected addition lines should not be included.
				if line.Type == DiffAddLine {
					continue
				}

				// Unselected deleted lines should be included instead as context lines.
				// line.Type == DiffDeleteLine
				hunkBuffer += " " + line.Line + "\n"
				oldCount++
				newCount++
			}

			if line.NoNewline {
				hunkBuffer += "\\ No newline at end of file\n"
			}
		}

		// If there are no changes in this hunk, skip it.
		if !changes {
			continue
		}

		patch += fmt.Sprintf("@@ -%d,%d +%d,%d @@ %s\n", hunk.StartOld, oldCount, hunk.StartNew, newCount, hunk.Heading)
		patch += hunkBuffer
	}

	if patch == "" {
		// Don't generate a header if there's no patch.
		return ""
	}
	return patchHeader(from, to) + patch
}

// Create a patch to unstage changes or discard from unstaged files.
func (d *Diff) createDiscardPatch(filename string) string {
	patch := ""

	for _, hunk := range d.Hunks {
		hunkBuffer := ""
		oldCount := 0
		newCount := 0
		changes := false

		for _, line := range hunk.Lines {
			if line.Type == DiffContextLine {
				hunkBuffer += " " + line.Line + "\n"
				oldCount++
				newCount++
			} else if line.Selected {
				changes = true
				// Reverse the changes. Add:- Delete:+
				if line.Type == DiffAddLine {
					hunkBuffer += "-" + line.Line + "\n"
					newCount++
				} else if line.Type == DiffDeleteLine {
					hunkBuffer += "+" + line.Line + "\n"
					oldCount++
				}
			} else {
				// Unselected delete lines should be ignored.
				if line.Type == DiffDeleteLine {
					continue
				}

				// Unselected addition lines should remain the same.
				// line.Type == DiffAddLine
				hunkBuffer += " " + line.Line + "\n"
				oldCount++
				newCount++
			}

			if line.NoNewline {
				hunkBuffer += "\\ No newline at end of file\n"
			}
		}

		// If there are no changes in this hunk, skip it.
		if !changes {
			continue
		}

		// Swap old and new.
		patch += fmt.Sprintf("@@ -%d,%d +%d,%d @@ %s\n", hunk.StartNew, newCount, hunk.StartOld, oldCount, hunk.Heading)
		patch += hunkBuffer
	}

	if patch == "" {
		// Don't generate a header if there's no patch.
		return ""
	}
	return patchHeader(filename, filename) + patch
}
