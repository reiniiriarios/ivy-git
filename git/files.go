package git

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileLineType string

const (
	LineNormal        = "LineNormal"
	LineOurs          = "LineOurs"
	LineTheirs        = "LineTheirs"
	LineConflictStart = "LineConflictStart"
	LineConflictFlip  = "LineConflictFlip"
	LineConflictEnd   = "LineConflictEnd"
)

type File struct {
	File         string
	Raw          string
	Lines        []FileLine
	NumConflicts uint16
	Resolved     bool
}

type FileLine struct {
	RawLineNo  uint64
	LineNo     uint64
	Line       string
	Type       FileLineType
	Conflict   uint16
	OursTheirs int8
}

func (g *Git) GetParsedConflictFile(file string) (File, error) {
	f, err := getFileContents(filepath.Join(g.Repo.Directory, file))
	if err != nil {
		return File{}, err
	}
	fileobj := File{
		File: file,
		Raw:  f,
	}
	fileobj.parse()
	return fileobj, nil
}

func (f *File) parse() error {
	var conflict_side int8 = DiffNeither
	lines := strings.Split(strings.ReplaceAll(f.Raw, "\r\n", "\n"), "\n")
	var line_no uint64 = 1
	var ours_count uint64 = 1
	for ln := 0; ln < len(lines); ln++ {
		if strings.HasPrefix(lines[ln], "<<<<<<<") {
			f.NumConflicts++
			f.Lines = append(f.Lines, FileLine{
				RawLineNo:  uint64(ln),
				LineNo:     line_no,
				Type:       LineConflictStart,
				Line:       lines[ln],
				Conflict:   f.NumConflicts,
				OursTheirs: DiffOurs,
			})
			conflict_side = DiffOurs
		} else if strings.HasPrefix(lines[ln], "=======") {
			f.Lines = append(f.Lines, FileLine{
				RawLineNo:  uint64(ln),
				LineNo:     line_no,
				Type:       LineConflictFlip,
				Line:       lines[ln],
				Conflict:   f.NumConflicts,
				OursTheirs: DiffNeither,
			})
			conflict_side = DiffTheirs
			// Reset the line counter
			line_no -= ours_count
			ours_count = 0
		} else if strings.HasPrefix(lines[ln], ">>>>>>>") {
			f.Lines = append(f.Lines, FileLine{
				RawLineNo:  uint64(ln),
				LineNo:     line_no,
				Type:       LineConflictEnd,
				Line:       lines[ln],
				Conflict:   f.NumConflicts,
				OursTheirs: DiffTheirs,
			})
			conflict_side = DiffNeither
		} else {
			var line_type FileLineType
			line_conflict := f.NumConflicts
			if conflict_side == -1 {
				line_type = LineOurs
				ours_count++
			} else if conflict_side == 1 {
				line_type = LineTheirs
			} else {
				line_type = LineNormal
				line_conflict = 0
			}
			f.Lines = append(f.Lines, FileLine{
				RawLineNo:  uint64(ln),
				LineNo:     line_no,
				Type:       line_type,
				Line:       lines[ln],
				Conflict:   line_conflict,
				OursTheirs: conflict_side,
			})
			line_no++
		}
	}
	return nil
}

func getFileContents(file string) (string, error) {
	if _, err := os.Stat(file); err != nil {
		if !os.IsNotExist(err) {
			println(err)
		}
		return "", err
	}

	msg, err := ioutil.ReadFile(file)
	if err != nil {
		println(err)
		return "", err
	}

	return string(msg), nil
}
