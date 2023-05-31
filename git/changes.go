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

// Get list of changed files.
func (g *Git) GitListChanges() (map[string]*Change, map[string]*Change, error) {
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
	changesXmap := make(map[string]*Change)
	changesYmap := make(map[string]*Change)

	c, err := g.RunCwd("status", "--untracked-files", "--porcelain", "-z")
	if err != nil {
		return changesXmap, changesYmap, err
	}

	// https://git-scm.com/docs/git-status
	c = parseOneLine(c)
	// The -z option splits lines by NUL.
	changes := strings.Split(c, "\x00")

	for i := 0; i < len(changes); i++ {
		if strings.Trim(changes[i], " ") == "" {
			continue
		}

		x := changes[i][0:1]
		y := changes[i][1:2]
		file := changes[i][3:]

		old_file := ""
		if x == "R" || y == "R" {
			// Renames get two lines of data, the second line is the old filename.
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

	for i := range changesX {
		changesXmap[changesX[i].File] = &changesX[i]
	}
	for i := range changesY {
		changesYmap[changesY[i].File] = &changesY[i]
	}

	return changesXmap, changesYmap, nil
}

// `git status --porcelain`
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

const (
	FileUnknownStatus = ""
	FileNotUpdated    = " "
	FileUntracked     = "?"
	FileIgnored       = "!"
	FileAdded         = "A"
	FileCopied        = "C"
	FileDeleted       = "D"
	FileModified      = "M"
	FileRenamed       = "R"
	FileTypeChanged   = "T"
	FileUnmerged      = "U"
)

func getStatusFlag(status string) string {
	switch status {
	case FileAdded:
		return "added"
	case FileCopied:
		return "copied"
	case FileDeleted:
		return "deleted"
	case FileModified:
		return "modified"
	case FileRenamed:
		return "renamed"
	case FileTypeChanged:
		return "type-changed"
	case FileUnmerged:
		return "unmerged"
	case FileUntracked:
		return "untracked"
	case FileIgnored:
		return "ignored"
	case FileNotUpdated:
		return "not-updated"
	}

	return "unknown"
}

func fileIsNew(status string) bool {
	return status == FileUntracked || status == FileAdded
}
