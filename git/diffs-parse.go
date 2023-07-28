package git

import (
	"errors"
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
	raw             string
	Flags           []string
	Hunks           []DiffHunk
	Binary          bool
	Empty           bool
	SelectableLines uint64
	SelectedLines   uint64
	NumConflicts    uint16
	Conflicts       map[int64]DiffConflict // map[MiniHunk_id]DiffConflict
	Lang            string
	Highlight       map[int]string
}

type DiffHunk struct {
	Header   string
	StartOld int64
	EndOld   int64
	StartNew int64
	EndNew   int64
	StartCur int64
	EndCur   int64
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
	lines := strings.Split(strings.ReplaceAll(d.raw, "\r\n", "\n"), "\n")
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
			// Last line didn't have a newline, don't add this line, just mark the previous
			d.Hunks[current_hunk].Lines[len(d.Hunks[current_hunk].Lines)-1].NoNewline = true
		case "+":
			d.SelectableLines++
			d.SelectedLines++
			mini_hunk_editable = true
			d.Hunks[current_hunk].Lines = append(d.Hunks[current_hunk].Lines, DiffLine{
				Line:      lines[ln][1:],
				Type:      DiffAddLine,
				RawLineNo: int64(ln) + 1,
				OldLineNo: -1,
				NewLineNo: after,
				CurLineNo: after,
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
				RawLineNo: int64(ln) + 1,
				OldLineNo: before,
				NewLineNo: -1,
				CurLineNo: -1,
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
				RawLineNo: int64(ln) + 1,
				OldLineNo: before,
				NewLineNo: after,
				CurLineNo: after,
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
	lines := strings.Split(strings.ReplaceAll(d.raw, "\r\n", "\n"), "\n")
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
				RawLineNo:  int64(ln) + 1,
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
				RawLineNo:  int64(ln) + 1,
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
				RawLineNo:  int64(ln) + 1,
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
					RawLineNo:  int64(ln) + 1,
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
					RawLineNo:  int64(ln) + 1,
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
					RawLineNo:  int64(ln) + 1,
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
		d.Hunks[current_hunk].EndCur = cur
	}
	if len(d.Hunks) == 0 {
		d.Empty = true
	}
	return nil
}

var hunkHeaderRegex = regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)

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
	matches := hunkHeaderRegex.FindStringSubmatch(line)
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
		StartCur: start_new,
		EndCur:   end_new,
	}, nil
}
