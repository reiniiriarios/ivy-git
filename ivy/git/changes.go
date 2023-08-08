package git

import (
	"errors"
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
	Them     string
	Us       string
	Flag     string
}

// Get list of changed files.
// Returns X changes, Y changes, conflicts
func (g *Git) GitListChanges() (map[string]*Change, map[string]*Change, map[string]*Change, error) {
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
	var changesC []Change

	c, err := g.run("status", "--untracked-files", "--porcelain", "-z")
	if err != nil {
		return nil, nil, nil, err
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
		if x == "R" || y == "R" || x == "C" || y == "C" {
			// Renames and copies get two lines of data, the second line is the old filename.
			old_file = changes[i+1]
			i++
		}

		// DD, AU, UD, UA, DU, AA, UU
		if x == "U" || y == "U" || (x == "A" && y == "A") || (x == "D" && y == "D") {
			changesC = append(changesC, Change{
				File:     file,
				Basename: filepath.Base(file),
				Dir:      filepath.Dir(file),
				OldFile:  old_file,
				Them:     x,
				Us:       y,
				Flag:     getUnmergedStatusFlag(x, y),
			})
		} else {
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
	}

	changesXmap := sortChanges(changesX)
	changesYmap := sortChanges(changesY)
	changesCmap := sortChanges(changesC)

	return changesXmap, changesYmap, changesCmap, nil
}

func sortChanges(changes []Change) map[string]*Change {
	if len(changes) > 0 {
		sort.Slice(changes, func(i, j int) bool {
			return strings.ToUpper(changes[i].Basename) < strings.ToUpper(changes[j].Basename)
		})
	}
	changesmap := make(map[string]*Change)
	for i := range changes {
		changesmap[changes[i].File] = &changes[i]
	}
	return changesmap
}

func (g *Git) getFileStatus(file string) (string, string, error) {
	s, err := g.run("status", "--untracked-files", "--porcelain", "-z", "--", file)
	if err != nil {
		return "", "", err
	}
	// XY filename
	s = parseOneLine(s)
	if len(strings.TrimSpace(s)) < 4 {
		return "", "", errors.New("cannot determine file status, file not found")
	}
	return s[0:1], s[1:2], nil
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

func getUnmergedStatusFlag(x string, y string) string {
	switch x + y {
	case "DD":
		return "unmerged-deleted"
	case "AA":
		return "unmerged-added"
	case "AU":
		return "unmerged-added-us"
	case "UA":
		return "unmerged-added-them"
	case "DU":
		return "unmerged-deleted-us"
	case "UD":
		return "unmerged-deleted-them"
	case "UU":
		return "unmerged-modified"
	}

	return "unknown"
}

func fileIsNew(status string) bool {
	return status == FileUntracked || status == FileAdded
}
