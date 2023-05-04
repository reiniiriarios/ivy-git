package main

import (
	"path/filepath"
	"sort"
	"strings"
)

type Change struct {
	File     string
	Basename string
	Dir      string
	Letter   string
	Flag     string
}

type ChangesResponse struct {
	Response string
	Message  string
	ChangesX []Change
	ChangesY []Change
}

// FRONTEND: Get list of changed files.
func (a *App) GitListChanges() ChangesResponse {
	c, err := a.GitCwd("status", "--untracked-files", "--porcelain")
	if err != nil {
		return ChangesResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	// https://git-scm.com/docs/git-status
	var changesX []Change
	var changesY []Change
	cs := a.getLines(c)
	for _, change := range cs {
		if strings.Trim(change, " ") != "" {
			x := change[0:1]
			file := change[2:]
			if x != " " && x != "?" {
				changesX = append(changesX, Change{
					File:     file,
					Basename: filepath.Base(file),
					Dir:      filepath.Dir(file),
					Letter:   x,
					Flag:     getStatusFlag(x),
				})
			}
			y := change[1:2]
			if y != " " {
				changesY = append(changesY, Change{
					File:     file,
					Basename: filepath.Base(file),
					Dir:      filepath.Dir(file),
					Letter:   y,
					Flag:     getStatusFlag(y),
				})
			}
		}
	}

	if len(changesX) > 0 {
		sort.Slice(changesX, func(i, j int) bool {
			return strings.ToUpper(changesX[i].Basename) < strings.ToUpper(changesX[j].Basename)
		})
	}

	if len(changesY) > 0 {
		sort.Slice(changesY, func(i, j int) bool {
			return strings.ToUpper(changesY[i].Basename) < strings.ToUpper(changesY[j].Basename)
		})
	}

	return ChangesResponse{
		Response: "success",
		ChangesX: changesX,
		ChangesY: changesY,
	}
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
