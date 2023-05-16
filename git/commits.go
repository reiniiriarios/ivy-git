package git

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type CommitAddl struct {
	Hash               string
	Body               string
	BodyHtml           string
	CommitterName      string
	CommitterEmail     string
	CommitterTimestamp int64
	CommitterDatetime  string
}

// Get additional commit details not listed in the table.
func (g *Git) GetCommitDetails(hash string) (CommitAddl, error) {
	// Include:
	// %an - Committer Name
	// %ae - Committer Email
	// %at - Committer Time
	// %b  - Body
	// https://git-scm.com/docs/pretty-formats
	data := []string{"%cn", "%ce", "%ct", "%b"}
	format := strings.Join(data, GIT_LOG_SEP)
	c, err := g.RunCwd("--no-pager", "log", hash, "--format='"+format+"'", "--max-count=1")
	if err != nil {
		return CommitAddl{}, err
	}

	c = strings.Trim(strings.Trim(strings.Trim(c, "\n"), "\r"), "'")
	parts := strings.Split(c, GIT_LOG_SEP)
	if len(parts) != len(data) {
		return CommitAddl{}, errors.New("error parsing git log")
	}

	// Get timestamp and formatted datetime for committer
	ts, err := strconv.ParseInt(parts[2], 10, 64)
	dt := ""
	if err != nil {
		println(err.Error())
	} else {
		dt = time.Unix(ts, 0).Format(DATE_FORMAT)
	}

	return CommitAddl{
		Hash:               hash,
		Body:               parts[3],
		BodyHtml:           mdToHTML(parts[3]),
		CommitterName:      parts[0],
		CommitterEmail:     parts[1],
		CommitterTimestamp: ts,
		CommitterDatetime:  dt,
	}, nil
}

type FileStat struct {
	File    string
	Name    string
	Dir     string
	Path    []string
	OldFile string
	OldName string
	OldDir  string
	Added   uint32
	Deleted uint32
	Status  string
}

type FileStatDir struct {
	Name  string
	Path  []string
	Files []FileStat
	Dirs  []FileStatDir
}

// Get commit diff summary from diff-tree --numstat and --name-status.
func (g *Git) GetCommitDiffSummary(hash string) (FileStatDir, error) {
	filestats := []FileStat{}

	// Get the number of lines added and deleted from each file.
	n, err := g.RunCwd("diff-tree", "--numstat", "-r", "--root", "--find-renames", "-z", hash)
	if err != nil {
		return FileStatDir{}, err
	}

	n = strings.Trim(strings.Trim(strings.Trim(n, "\n"), "\r"), "'")
	// The -z option splits lines by NUL.
	nl := strings.Split(n, "\x00")
	// The first line is the hash, skip.
	for i := 1; i < len(nl); i++ {
		nf := strings.Fields(nl[i])
		if len(nf) == 3 {
			a, _ := strconv.ParseInt(nf[0], 10, 32)
			d, _ := strconv.ParseInt(nf[1], 10, 32)
			name := filepath.Base(nf[2])
			dir := filepath.Dir(nf[2])
			path := strings.Split(strings.ReplaceAll(dir, "\\", "/"), "/")
			filestats = append(filestats, FileStat{
				File:    nf[2],
				Name:    name,
				Dir:     dir,
				Path:    path,
				Added:   uint32(a),
				Deleted: uint32(d),
			})
		} else if len(nf) == 2 {
			// If there are two fields parsed, the next two lines are the
			// previous name and the new name.
			a, _ := strconv.ParseInt(nf[0], 10, 32)
			d, _ := strconv.ParseInt(nf[1], 10, 32)
			i++
			oldfile := nl[i]
			i++
			file := nl[i]
			name := filepath.Base(file)
			dir := filepath.Dir(file)
			path := strings.Split(strings.ReplaceAll(dir, "\\", "/"), "/")
			oldname := filepath.Base(oldfile)
			olddir := filepath.Dir(oldfile)
			filestats = append(filestats, FileStat{
				File:    file,
				Name:    name,
				Dir:     dir,
				Path:    path,
				OldFile: oldfile,
				OldName: oldname,
				OldDir:  olddir,
				Added:   uint32(a),
				Deleted: uint32(d),
			})
		}
	}

	// Get the status of each file in the commit.
	s, err := g.RunCwd("diff-tree", "--name-status", "-r", "--root", "--find-renames", "-z", hash)
	if err != nil {
		return FileStatDir{}, err
	}

	s = strings.Trim(strings.Trim(strings.Trim(s, "\n"), "\r"), "'")
	// The -z option splits lines by NUL.
	sl := strings.Split(s, "\x00")
	// The first line is the hash, skip.
	// Each line is either the status or the file. Parse two lines at a time.
	for i := 1; i < len(sl)-1; i += 2 {
		for f := range filestats {
			// Renames get three lines of data.
			if sl[i][:1] == "R" {
				if filestats[f].File == sl[i+2] {
					filestats[f].Status = "R"
					i++
					break
				}
			} else {
				if filestats[f].File == sl[i+1] {
					filestats[f].Status = sl[i]
					break
				}
			}
		}
	}

	// Parse files into directory tree.
	files := FileStatDir{}
	for _, f := range filestats {
		c := &files
		for n, p := range f.Path {
			if p != "." {
				added := false
				for j := range c.Dirs {
					if c.Dirs[j].Name == p {
						added = true
						c = &c.Dirs[j]
					}
				}
				if !added {
					c.Dirs = append(c.Dirs, FileStatDir{
						Name: p,
						Path: append(f.Path[:n], p),
					})
					c = &c.Dirs[len(c.Dirs)-1]
				}
			}
			if n == len(f.Path)-1 {
				c.Files = append(c.Files, f)
			}
		}
	}

	// Trim tree.
	trimDirs(&files)

	return files, nil
}

// Trim dirs with no contents except one dir.
// This collapses dirs to, e.g.
//
//	foo / bar
//	  baz.go
//
// if foo doesn't have any files changed and only one subdir, bar.
func trimDirs(dir *FileStatDir) {
	if dir.Name != "" && len(dir.Dirs) == 1 && len(dir.Files) == 0 {
		dir.Dirs[0].Name = dir.Name + " / " + dir.Dirs[0].Name
		*dir = dir.Dirs[0]
		trimDirs(dir)
	} else {
		for d := range dir.Dirs {
			trimDirs(&dir.Dirs[d])
		}
	}
}
