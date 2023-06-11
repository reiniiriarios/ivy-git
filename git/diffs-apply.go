package git

import (
	"ivy-git/files"
	"path/filepath"
)

// Resolve diff conflicts by editing the file.
func (g *Git) ResolveDiffConflicts(d Diff) error {
	del, rep := d.parseConflictResolutions()
	err := files.DeleteAndReplaceLinesFromFile(filepath.Join(g.Repo.Directory, d.File), del, rep)
	return err
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
				for j := d.Conflicts[i].Theirs[0].CurLineNo; j <= d.Conflicts[i].Theirs[len(d.Conflicts[i].Theirs)-1].CurLineNo; j++ {
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
					for j := d.Conflicts[i].Ours[0].CurLineNo; j <= d.Conflicts[i].Ours[len(d.Conflicts[i].Ours)-1].CurLineNo; j++ {
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
