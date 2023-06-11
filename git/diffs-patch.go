package git

import "fmt"

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
