package git

import (
	"errors"
	"fmt"
	"ivy-git/files"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type OursTheirs int8

const (
	DiffAddLine         = "DiffAddLine"
	DiffDeleteLine      = "DiffDeleteLine"
	DiffContextLine     = "DiffContextLine"
	DiffChangeStartLine = "DiffChangeStartLine"
	DiffChangeEndLine   = "DiffChangeEndLine"
	DiffChangeFlipLine  = "DiffChangeFlipLine"

	DiffOurs        OursTheirs = -1 // Between <<<<<<< and =======
	DiffNeither     OursTheirs = 0  // Outside <<<<<<< and >>>>>>>
	DiffTheirs      OursTheirs = 1  // Between ======= and >>>>>>>
	DiffBoth        OursTheirs = 2  // Between <<<<<<< and ======= first, ======= and >>>>>>> second
	DiffBothInverse OursTheirs = -2 // Between ======= and >>>>>>> first, <<<<<<< and ======= second
)

type Diff struct {
	File            string
	Raw             string
	Hunks           []DiffHunk
	Binary          bool
	Empty           bool
	SelectableLines uint64
	SelectedLines   uint64
	NumConflicts    uint16
	Conflicts       map[int64]DiffConflict // map[MiniHunk_id]DiffConflict
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
	Line       string
	Type       string
	RawLineNo  int64 // Line number in the diff
	OldLineNo  int64 // Line number in the old version
	NewLineNo  int64 // Line number in the new version
	CurLineNo  int64 // Line number in the current file (during conflicts)
	NoNewline  bool
	MiniHunk   int64
	Selected   bool
	OursTheirs OursTheirs
}

type DiffConflict struct {
	Ours       []DiffLine
	Theirs     []DiffLine
	Markers    []int64 // <<<<<<<, =======, >>>>>>>
	Resolution OursTheirs
}

type DiffConflicts []DiffConflict

const HiddenBidiCharsRegex = "/[\u202A-\u202E]|[\u2066-\u2069]/"

func (d *Diff) parse() error {
	// Do not use parseLines() here, will strip relevant data.
	lines := strings.Split(strings.ReplaceAll(d.Raw, "\r\n", "\n"), "\n")
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

		// Skip empty lines (e.g. trailing line)
		if len(lines[ln]) == 0 {
			continue
		}

		// Parse other lines
		prefix := lines[ln][:1]
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
			after++
		case "-":
			d.SelectableLines++
			d.SelectedLines++
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
			before++
		case " ":
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
			before++
			after++
		default:
			return errors.New("malformed diff, unrecognized line type")
		}
	}
	if len(d.Hunks) == 0 {
		d.Empty = true
	}
	return nil
}

func (d *Diff) parseConflicts() error {
	d.Conflicts = make(map[int64]DiffConflict)

	// Do not use parseLines() here, will strip relevant data.
	lines := strings.Split(strings.ReplaceAll(d.Raw, "\r\n", "\n"), "\n")
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
	var before, after, cur int64
	var change_mini_hunk int64 = 0
	var conflict_ours_theirs OursTheirs = DiffNeither
	var conflict_before int16 = 0
	var conflict_after int16 = 0
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
			cur = after

			continue
		}
		// Keep looking until we've found at least the first hunk
		if current_hunk == -1 {
			continue
		}

		// Skip empty lines (e.g. trailing line)
		if len(lines[ln]) == 0 {
			continue
		}

		var mini_hunk int64 = -1
		var new_line DiffLine

		if strings.HasPrefix(lines[ln][1:], "<<<<<<<") {
			change_mini_hunk++
			d.NumConflicts++
			new_line = DiffLine{
				Line:       lines[ln][1:],
				Type:       DiffChangeStartLine,
				RawLineNo:  int64(ln),
				OldLineNo:  -1,
				NewLineNo:  -1,
				CurLineNo:  cur,
				MiniHunk:   change_mini_hunk,
				OursTheirs: DiffOurs,
			}
			// Add this marker line to its conflict.
			c := d.Conflicts[change_mini_hunk]
			c.Markers = append(c.Markers, cur)
			d.Conflicts[change_mini_hunk] = c
			cur++
			// Start the "ours" side of the conflict.
			conflict_ours_theirs = DiffOurs

		} else if strings.HasPrefix(lines[ln][1:], "=======") {
			new_line = DiffLine{
				Line:       lines[ln][1:],
				Type:       DiffChangeFlipLine,
				RawLineNo:  int64(ln),
				OldLineNo:  -1,
				NewLineNo:  -1,
				CurLineNo:  cur,
				MiniHunk:   change_mini_hunk,
				OursTheirs: DiffNeither,
			}
			// Add this marker line to its conflict.
			c := d.Conflicts[change_mini_hunk]
			c.Markers = append(c.Markers, cur)
			d.Conflicts[change_mini_hunk] = c
			cur++
			// Start the "theirs" side of the conflict.
			conflict_ours_theirs = DiffTheirs
			// Roll back before and after to before the "ours" side.
			before -= int64(conflict_before)
			after -= int64(conflict_after)
			// Reset these counters.
			conflict_before = 0
			conflict_after = 0

		} else if strings.HasPrefix(lines[ln][1:], ">>>>>>>") {
			new_line = DiffLine{
				Line:       lines[ln][1:],
				Type:       DiffChangeEndLine,
				RawLineNo:  int64(ln),
				OldLineNo:  -1,
				NewLineNo:  -1,
				CurLineNo:  cur,
				MiniHunk:   change_mini_hunk,
				OursTheirs: DiffTheirs,
			}
			// Add this marker line to its conflict.
			c := d.Conflicts[change_mini_hunk]
			c.Markers = append(c.Markers, cur)
			d.Conflicts[change_mini_hunk] = c
			cur++
			// End conflict.
			conflict_ours_theirs = DiffNeither
			change_mini_hunk++

		} else {
			if conflict_ours_theirs != 0 {
				mini_hunk = change_mini_hunk
			}
			// Parse other lines
			prefix := lines[ln][:1]
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
				new_line = DiffLine{
					Line:       lines[ln][1:],
					Type:       DiffAddLine,
					RawLineNo:  int64(ln),
					OldLineNo:  -1,
					NewLineNo:  after,
					CurLineNo:  cur,
					MiniHunk:   mini_hunk,
					OursTheirs: conflict_ours_theirs,
				}
				// Rolling new line number.
				after++
				cur++
				// Rolling new line number offset.
				if conflict_ours_theirs == DiffOurs {
					conflict_after++
				}
			case "-":
				new_line = DiffLine{
					Line:       lines[ln][1:],
					Type:       DiffDeleteLine,
					RawLineNo:  int64(ln),
					OldLineNo:  before,
					NewLineNo:  -1,
					CurLineNo:  -1,
					MiniHunk:   mini_hunk,
					OursTheirs: conflict_ours_theirs,
				}
				// Rolling old line number.
				before++
				// Rolling old line number offset.
				if conflict_ours_theirs == DiffOurs {
					conflict_before++
				}
			case " ":
				new_line = DiffLine{
					Line:       lines[ln][1:],
					Type:       DiffContextLine,
					RawLineNo:  int64(ln),
					OldLineNo:  before,
					NewLineNo:  after,
					CurLineNo:  cur,
					MiniHunk:   mini_hunk,
					OursTheirs: conflict_ours_theirs,
				}
				// Rolling old and new line number.
				before++
				after++
				cur++
				// Rolling old and new line number offsets.
				if conflict_ours_theirs == -1 {
					conflict_before++
					conflict_after++
				}
			default:
				return errors.New("malformed diff, unrecognized line type")
			}
			// Add to list of conflicts.
			if conflict_ours_theirs == DiffOurs {
				c := d.Conflicts[mini_hunk]
				c.Ours = append(d.Conflicts[mini_hunk].Ours, new_line)
				d.Conflicts[mini_hunk] = c
			} else if conflict_ours_theirs == DiffTheirs {
				c := d.Conflicts[mini_hunk]
				c.Theirs = append(d.Conflicts[mini_hunk].Theirs, new_line)
				d.Conflicts[mini_hunk] = c
			}
		}
		d.Hunks[current_hunk].Lines = append(d.Hunks[current_hunk].Lines, new_line)
	}
	if len(d.Hunks) == 0 {
		d.Empty = true
	}
	return nil
}

// Parse a hunk header.
//
// Currently only supports default format:
//
//	@@ -n,n +n,n @@ heading
//	Where -n,n is the old line number range and +n,n is the new range.
//
// todo: allow support for custom hunk headers.
//
// https://git-scm.com/docs/diff-format
// https://git-scm.com/docs/git-diff
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

// Parse conflict resolutions into a list of lines to delete and
// a map of replacement strings for some lines.
func (d *Diff) parseConflictResolutions() ([]int64, map[int64]int64) {
	var delete_lines []int64
	replace_lines := make(map[int64]int64)

	for i := range d.Conflicts {
		switch d.Conflicts[i].Resolution {

		case DiffOurs:
			if len(d.Conflicts[i].Theirs) > 0 {
				// Loop from the first line number to the last. This ensures
				// lines outside of displayed hunks are included.
				for j := d.Conflicts[i].Theirs[0].CurLineNo; j <= d.Conflicts[i].Theirs[len(d.Conflicts[i].Theirs)-1].CurLineNo; j++ {
					delete_lines = append(delete_lines, j)
				}
			}

		case DiffTheirs:
			if len(d.Conflicts[i].Ours) > 0 {
				// Loop from the first line number to the last. This ensures
				// lines outside of displayed hunks are included.
				println("delete ours")
				for j := d.Conflicts[i].Ours[0].CurLineNo; j <= d.Conflicts[i].Ours[len(d.Conflicts[i].Ours)-1].CurLineNo; j++ {
					println(j)
					delete_lines = append(delete_lines, j)
				}
			}

		case DiffBoth:
			// Nothing to do.

		case DiffBothInverse:
			// If Ours...
			ours_len := 0
			theirs_len := 0
			if len(d.Conflicts[i].Ours) > 0 {
				// Start at the beginning of Ours and add Theirs as replacement lines.
				ln := d.Conflicts[i].Ours[0].CurLineNo
				println("start", ln)
				for j := d.Conflicts[i].Theirs[0].CurLineNo; j <= d.Conflicts[i].Theirs[len(d.Conflicts[i].Theirs)-1].CurLineNo; j++ {
					println("rep", ln, "with", j)
					replace_lines[ln] = j
					ln++
					theirs_len++
				}
				// Move the ======= delete line to the new center of ours vs theirs.
				d.Conflicts[i].Markers[1] = ln
				ln++
				// If Ours and Theirs...
				if len(d.Conflicts[i].Theirs) > 0 {
					// Add Ours below Theirs.
					println("theirs", ln)
					for j := d.Conflicts[i].Ours[0].CurLineNo; j <= d.Conflicts[i].Ours[len(d.Conflicts[i].Ours)-1].CurLineNo; j++ {
						println("rep", ln, "with", j)
						replace_lines[ln] = j
						ln++
						ours_len++
					}
				}
				// If just Theirs...
			} else if len(d.Conflicts[i].Theirs) > 0 {
				ln := d.Conflicts[i].Theirs[0].CurLineNo
				// Move the ======= delete line to the new center of ours vs theirs.
				d.Conflicts[i].Markers[1] = ln - 1
				for j := d.Conflicts[i].Ours[0].CurLineNo; j <= d.Conflicts[i].Ours[len(d.Conflicts[i].Ours)-1].CurLineNo; j++ {
					replace_lines[ln] = j
					ln++
					ours_len++
				}
			}
		}

		delete_lines = append(delete_lines, d.Conflicts[i].Markers...)
	}

	return delete_lines, replace_lines
}

// Resolve diff conflicts by editing the file.
func (g *Git) ResolveDiffConflicts(d Diff) error {
	del, rep := d.parseConflictResolutions()
	err := files.DeleteAndReplaceLinesFromFile(filepath.Join(g.Repo.Directory, d.File), del, rep)
	return err
}
