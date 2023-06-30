package git

import (
	"errors"
	"os"
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

func (g *Git) GetLastCommitHash() (string, error) {
	h, err := g.run("--no-pager", "log", "--format=%H", "--max-count=1")
	if err != nil {
		if errorCode(err) == NoCommitsYet || errorCode(err) == BadRevision || errorCode(err) == UnknownRevisionOrPath || errorCode(err) == ExitStatus1 {
			return "", nil
		}
		return "", err
	}
	h = parseOneLine(h)
	return h, nil
}

// Get subject and body of most recent commit message.
func (g *Git) GetLastCommitMessage() (CommitMessage, error) {
	m, err := g.run("--no-pager", "log", "--format=%s"+GIT_LOG_SEP+"%b", "--max-count=1")
	if err != nil {
		return CommitMessage{}, err
	}
	ms := strings.Split(m, GIT_LOG_SEP)
	return CommitMessage{
		Subject: strings.TrimSpace(ms[0]),
		Body:    strings.TrimSpace(ms[1]),
	}, nil
}

// Get additional commit details not listed in the table.
func (g *Git) GetCommitDetails(hash string, date_format string) (CommitAddl, error) {
	if hash == "" {
		return CommitAddl{}, errors.New("no commit hash specified")
	}

	// Include:
	// %cn - Committer Name
	// %ce - Committer Email
	// %at - Committer Time
	// %b  - Body
	// https://git-scm.com/docs/pretty-formats
	data := []string{}
	if GIT_RESPECT_MAILMAP {
		data = append(data, "%cN", "%cE")
	} else {
		data = append(data, "%cn", "%ce")
	}
	data = append(data, "%ct", "%b")

	format := strings.Join(data, GIT_LOG_SEP)
	c, err := g.run("--no-pager", "log", hash, "--format="+format, "--max-count=1")
	if err != nil {
		return CommitAddl{}, err
	}

	c = parseOneLine(c)
	parts := strings.Split(c, GIT_LOG_SEP)
	if len(parts) != len(data) {
		return CommitAddl{}, errors.New("error parsing git log")
	}

	// Get timestamp and formatted datetime for committer
	ts, err := strconv.ParseInt(parts[2], 10, 64)
	dt := ""
	if err == nil {
		dt = time.Unix(ts, 0).Format(date_format)
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
	OldRel  string
	Added   uint32
	Deleted uint32
	Binary  bool
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
	if hash == "" {
		return FileStatDir{}, errors.New("no commit hash specified")
	}

	filestats := []FileStat{}

	parents, err := g.getCommitParents(hash)
	if err != nil {
		return FileStatDir{}, err
	}

	// Get the number of lines added and deleted from each file.
	var numstat string
	var merge_base string
	var ignore_first_line bool
	if len(parents) > 1 {
		merge_base, err = g.findMergeBase(parents...)
		if err != nil {
			return FileStatDir{}, err
		}
		numstat, err = g.run("diff-tree", "--numstat", "-r", "--root", "--find-renames", "-z", hash, merge_base)
		ignore_first_line = false
	} else {
		numstat, err = g.run("diff-tree", "--numstat", "-r", "--root", "--find-renames", "-z", hash)
		ignore_first_line = true
	}
	if err != nil {
		return FileStatDir{}, err
	}

	numstat = parseOneLine(numstat)
	// The -z option splits lines by NUL.
	nl := strings.Split(numstat, "\x00")

	// If a second hash is not supplied, the first line is the commit hash,
	// and so we start counting at one, ignoring the first line of output.
	i := 0
	if ignore_first_line {
		i = 1
	}

	for ; i < len(nl); i++ {
		nf := strings.Fields(nl[i])
		binary := false
		var a int64 = 0
		var d int64 = 0
		if len(nf) >= 2 {
			if nf[0] == "-" && nf[1] == "-" {
				binary = true
			} else {
				a, _ = strconv.ParseInt(nf[0], 10, 32)
				d, _ = strconv.ParseInt(nf[1], 10, 32)
			}
		}

		if len(nf) == 3 {
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
				Binary:  binary,
			})
		} else if len(nf) == 2 {
			// If there are two fields parsed, the next two lines are the previous name and the new name.
			i++
			oldfile := nl[i]
			i++
			file := nl[i]
			name := filepath.Base(file)
			dir := filepath.Dir(file)
			path := strings.Split(strings.ReplaceAll(dir, "\\", "/"), "/")
			oldname := filepath.Base(oldfile)
			olddir := filepath.Dir(oldfile)
			oldrel, _ := filepath.Rel(dir, olddir)
			if strings.HasSuffix(oldrel, string(os.PathSeparator)+".") {
				oldrel = oldrel[:len(oldrel)-2]
			}
			oldrel = strings.ReplaceAll(oldrel, "\\", "/")

			filestats = append(filestats, FileStat{
				File:    file,
				Name:    name,
				Dir:     dir,
				Path:    path,
				OldFile: oldfile,
				OldName: oldname,
				OldDir:  olddir,
				OldRel:  oldrel,
				Added:   uint32(a),
				Deleted: uint32(d),
				Binary:  binary,
			})
		}
	}

	// Get the status of each file in the commit.
	var name_status string
	if len(parents) > 1 {
		name_status, err = g.run("diff-tree", "--name-status", "-r", "--root", "--find-renames", "-z", hash, merge_base)
	} else {
		name_status, err = g.run("diff-tree", "--name-status", "-r", "--root", "--find-renames", "-z", hash)
	}
	if err != nil {
		return FileStatDir{}, err
	}

	name_status = parseOneLine(name_status)
	// The -z option splits lines by NUL.
	sl := strings.Split(name_status, "\x00")
	// If not a merge commit, he first line is the hash, skip.
	start := 1
	if len(parents) > 1 {
		start = 0
	}
	// Each line is either the status or the file. Parse two lines at a time.
	for i := start; i < len(sl)-1; i += 2 {
		for f := range filestats {
			// Renames and copies get three lines of data.
			if sl[i][:1] == "R" || sl[i][:1] == "C" {
				if filestats[f].File == sl[i+2] {
					filestats[f].Status = sl[i][:1]
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

// Get commit parents hashes
func (g *Git) getCommitParents(hash string) ([]string, error) {
	if hash == "" {
		return []string{}, errors.New("no commit hash specified")
	}

	c, err := g.run("--no-pager", "log", "--format=%P", "--max-count=1", hash)
	if err != nil {
		return []string{}, err
	}
	c = parseOneLine(c)
	parents := strings.Split(c, " ")

	return parents, nil
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

type CommitSignature struct {
	Status      string
	Name        string
	Key         string
	Description string
}

// Get simple signature status of commit list.
func (g *Git) GetCommitSignature(hash string) (CommitSignature, error) {
	if hash == "" {
		return CommitSignature{}, errors.New("no commit hash specified")
	}

	// Include:
	// %G?
	//   G = good (valid)
	//   B = bad
	//   U = unknown validity
	//   X = expired
	//   Y = good signature, expired key
	//   E = missing key
	//   N = no signature
	// %GS
	//   Signer name
	// %GK
	//   Key
	// https://git-scm.com/docs/pretty-formats
	data := []string{"%G?", "%GS", "%GK"}
	format := strings.Join(data, GIT_LOG_SEP)
	c, err := g.run("--no-pager", "log", hash, "--format="+format, "--max-count=1")
	if err != nil {
		return CommitSignature{}, err
	}

	c = parseOneLine(c)
	parts := strings.Split(c, GIT_LOG_SEP)
	if len(parts) != len(data) {
		return CommitSignature{}, errors.New("error parsing commit signature")
	}

	var desc string
	switch parts[0] {
	case "G":
		desc = "Good signature"
	case "B":
		desc = "Bad signature"
	case "U":
		desc = "Signature has unknown validity"
	case "X":
		desc = "Signature is expired"
	case "Y":
		desc = "Signed with expired key"
	case "E":
		desc = "No key available to validate signature"
	case "N":
		desc = "None"
	}

	return CommitSignature{
		Status:      parts[0],
		Name:        parts[1],
		Key:         parts[2],
		Description: desc,
	}, nil
}
