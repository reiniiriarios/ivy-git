package git

import (
	"path/filepath"
	"sort"
	"strings"
)

type Change struct {
	File     string
	Basename string
	Dir      string
	OldFile  string
	Letter   string
	Flag     string
}

type ChangesResponse struct {
	Response string
	Message  string
	ChangesX []Change
	ChangesY []Change
}

// Get list of changed files.
func (g *Git) GitListChanges() ([]Change, []Change, error) {
	// When a merge is occurring and the merge was successful, or outside of a merge situation,
	//   X shows the status of the index and Y shows the status of the working tree.
	// When a merge conflict has occurred and has not yet been resolved,
	//   X and Y show the state introduced by each head of the merge, relative to the common ancestor.
	//   These paths are said to be unmerged.
	// When a path is untracked, X and Y are always the same, since they are unknown to the index.
	//   ?? is used for untracked paths. Ignored files are not listed unless --ignored is used;
	//   if it is, ignored files are indicated by !!.

	var changesX []Change
	var changesY []Change

	c, err := g.RunCwd("status", "--untracked-files", "--porcelain", "-z")
	if err != nil {
		return changesX, changesY, err
	}

	// https://git-scm.com/docs/git-status
	c = parseOneLine(c)
	// The -z option splits lines by NUL.
	changes := strings.Split(c, "\x00")

	for i := 0; i < len(changes); i++ {
		if strings.Trim(changes[i], " ") == "" {
			continue
		}

		println(changes[i])
		// Renames get three lines of data.
		x := changes[i][0:1]
		y := changes[i][1:2]
		file := changes[i][3:]

		old_file := ""
		if x == "R" || y == "R" {
			old_file = changes[i+1]
			i++
		}

		if x != " " && x != "?" {
			changesX = append(changesX, Change{
				File:     file,
				Basename: filepath.Base(file),
				Dir:      filepath.Dir(file),
				OldFile:  old_file,
				Letter:   x,
				Flag:     getStatusFlag(x),
			})
		}
		if y != " " {
			changesY = append(changesY, Change{
				File:     file,
				Basename: filepath.Base(file),
				Dir:      filepath.Dir(file),
				OldFile:  old_file,
				Letter:   y,
				Flag:     getStatusFlag(y),
			})
		}
	}

	// Sort X changes by alpha.
	if len(changesX) > 0 {
		sort.Slice(changesX, func(i, j int) bool {
			return strings.ToUpper(changesX[i].Basename) < strings.ToUpper(changesX[j].Basename)
		})
	}

	// Sort Y changes by alpha.
	if len(changesY) > 0 {
		sort.Slice(changesY, func(i, j int) bool {
			return strings.ToUpper(changesY[i].Basename) < strings.ToUpper(changesY[j].Basename)
		})
	}

	return changesX, changesY, nil
}

// Get status flag for `git status --porcelain`
//
// X          Y     Meaning
// -------------------------------------------------
//
//	[AMD]   not updated
//
// M        [ MTD]  updated in index
// T        [ MTD]  type changed in index
// A        [ MTD]  added to index
// D                deleted from index
// R        [ MTD]  renamed in index
// C        [ MTD]  copied in index
// [MTARC]          index and work tree matches
// [ MTARC]    M    work tree changed since index
// [ MTARC]    T    type changed in work tree since index
// [ MTARC]    D    deleted in work tree
//
//	R    renamed in work tree
//	C    copied in work tree
//
// -------------------------------------------------
// D           D    unmerged, both deleted
// A           U    unmerged, added by us
// U           D    unmerged, deleted by them
// U           A    unmerged, added by them
// D           U    unmerged, deleted by us
// A           A    unmerged, both added
// U           U    unmerged, both modified
// -------------------------------------------------
// ?           ?    untracked
// !           !    ignored
// -------------------------------------------------
func getStatusFlag(status string) string {
	switch status {
	case "A":
		return "added"
	case "C":
		return "copied"
	case "D":
		return "deleted"
	case "M":
		return "modified"
	case "R":
		return "renamed"
	case "T":
		return "type-changed"
	case "U":
		return "unmerged"
	case "?":
		return "untracked"
	case "!":
		return "ignored"
	case " ":
		return "not-updated"
	}

	return "unknown"
}
