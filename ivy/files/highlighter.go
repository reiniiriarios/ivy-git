package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// https://github.com/alecthomas/chroma/tree/master/styles
const HIGHLIGHT_STYLE = "vs"

func HighlightFile(path string) (File, error) {
	return HighlightFileSelection(path, highlightRanges{})
}

func HighlightDiff(path string, diff string) (File, error) {
	diff = cleanDiffToCode(diff)

	f := File{
		Filename: filepath.Base(path),
		Filepath: path,
		Diff:     true,
		raw:      diff,
	}

	err := f.highlight()
	if err != nil {
		return f, err
	}

	return f, nil
}

// Strip hunk headers and line prefixes from diff.
func cleanDiffToCode(diff string) string {
	clean := ""
	lines := strings.Split(diff, "\n")
	for _, l := range lines {
		if len(l) > 0 &&
			!strings.HasPrefix(l, "@@") &&
			!strings.HasPrefix(l, ">>>>>>>") &&
			!strings.HasPrefix(l, "<<<<<<<") &&
			!strings.HasPrefix(l, "=======") {
			clean += l[1:] + "\n"
		} else {
			// Add blank lines for the diff markup so the line numbering is correct.
			clean += "\n"
		}
	}
	return clean
}

func HighlightFileSelection(path string, ranges highlightRanges) (File, error) {
	f := File{
		Filename: filepath.Base(path),
		Filepath: path,
		ranges:   ranges,
	}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			f.NotFound = true
			return f, nil
		}
		return f, err
	}

	err = f.getSize()
	if err != nil {
		return f, err
	}
	if f.Size > MAX_FILE_SIZE {
		f.TooLarge = true
		return f, nil
	}

	err = f.getContents()
	if err != nil {
		return f, err
	}

	err = f.highlight()
	if err != nil {
		return f, err
	}

	return f, nil
}

func (f *File) getSize() error {
	fp, err := os.Open(f.Filepath)
	if err != nil {
		return err
	}
	defer fp.Close()

	size, err := getFileSize(f.Filepath)
	if err != nil {
		return err
	}
	f.Size = size

	return nil
}

func (f *File) getContents() error {
	bytes, err := ioutil.ReadFile(f.Filepath)
	if err != nil {
		return err
	}
	f.raw = string(bytes)
	return nil
}

func (f *File) highlight() error {
	lexer := lexers.Match(f.Filepath)
	if lexer == nil {
		lexer = lexers.Analyse(f.raw)
		if lexer == nil {
			return nil
		}
	}
	f.Lang = lexer.Config().Name

	iterator, err := lexer.Tokenise(nil, f.raw)
	if err != nil {
		return err
	}

	style := styles.Get(HIGHLIGHT_STYLE)
	if style == nil {
		style = styles.Fallback
	}
	formatter := LinesFormatter(HighlightLines(f.ranges))

	f.Highlight = formatter.Format(style, iterator)

	return nil
}
