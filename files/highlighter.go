package files

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// https://github.com/alecthomas/chroma/tree/master/styles
const HIGHLIGHT_STYLE = "vs"

func HighlightFile(path string) (File, error) {
	f := File{
		Filename: filepath.Base(path),
		Filepath: path,
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
	f.Raw = string(bytes)
	return nil
}

func (f *File) highlight() error {
	lexer := lexers.Match(f.Filepath)
	if lexer == nil {
		lexer = lexers.Analyse(f.Raw)
		if lexer == nil {
			return nil
		}
	}
	f.Lang = lexer.Config().Name

	iterator, err := lexer.Tokenise(nil, f.Raw)
	if err != nil {
		return err
	}

	style := styles.Get(HIGHLIGHT_STYLE)
	if style == nil {
		style = styles.Fallback
	}
	formatter := LinesFormatter()

	f.Highlight = formatter.Format(style, iterator)

	return nil
}
