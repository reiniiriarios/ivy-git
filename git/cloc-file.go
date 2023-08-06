package git

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"unicode"

	enry "github.com/go-enry/go-enry/v2"
)

// ClocFile is collecting to line count result.
type ClocFile struct {
	Code     int64
	Comments int64
	Blanks   int64
	Bytes    int64
	Name     string
	Lang     string
}

// ClocFiles is gocloc result set.
type ClocFiles []ClocFile

func (cf ClocFiles) Len() int {
	return len(cf)
}

func (cf ClocFiles) Swap(i, j int) {
	cf[i], cf[j] = cf[j], cf[i]
}

func (cf ClocFiles) Less(i, j int) bool {
	if cf[i].Code == cf[j].Code {
		return cf[i].Name < cf[j].Name
	}
	return cf[i].Code > cf[j].Code
}

func (g *Git) readLanguageFromFileOnBranch(file string) string {
	// Swap to forward slashes. Even on windows, git show requires this.
	file = strings.ReplaceAll(file, "\\", "/")

	// Run command to show a specific file on a specific branch in a specific repo.
	cmd := []string{
		"-C",
		g.Repo.Directory,
		"--no-pager",
		"show",
		g.Repo.Main + ":" + file,
	}
	content, err := g.run(cmd...)
	if err != nil {
		println("Error (readLanguageFromFileOnBranch)", err.Error())
		return ""
	}
	lang := enry.GetLanguage(file, []byte(content))

	return lang
}

func (g *Git) clocAnalyzeFileOnBranch(file string, language *Language) (*ClocFile, error) {
	cfile := &ClocFile{}

	// Swap to forward slashes. Even on windows, git show requires this.
	file = strings.ReplaceAll(file, "\\", "/")

	// We're going to pipe the output of the git command and scan the pipe reader.
	reader, writer := io.Pipe()

	// Analyze pipe stream.
	scannerStopped := make(chan struct{})
	go func() {
		defer close(scannerStopped)
		cfile = clocAnalyzeReader(file, language, reader)
	}()

	// Run command to show a specific file on a specific branch in a specific repo.
	command := []string{
		"-C",
		g.Repo.Directory,
		"--no-pager",
		"show",
		g.Repo.Main + ":" + file,
	}
	cmd := exec.Command("git", command...)

	// On windows, we need to hide the command prompt.
	hideCmdPrompt(cmd)

	// Set stdout and make go.
	cmd.Stdout = writer
	// debug
	// cmd.Stderr = writer
	err := cmd.Start()
	go func() {
		err = cmd.Wait()
		writer.Close()
	}()

	// And we're done.
	<-scannerStopped

	return cfile, err
}

type ScanByteCounter struct {
	BytesRead int64
}

func (s *ScanByteCounter) Wrap(split bufio.SplitFunc) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (int, []byte, error) {
		adv, tok, err := split(data, atEOF)
		s.BytesRead += int64(adv)
		return adv, tok, err
	}
}

func clocAnalyzeReader(filename string, language *Language, file *io.PipeReader) *ClocFile {
	clocFile := &ClocFile{
		Name: filename,
		Lang: language.Data.Name,
	}

	isFirstLine := true
	inComments := [][2]string{}
	buf := getByteSlice()
	defer putByteSlice(buf)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf.Bytes(), 1024*1024)
	counter := ScanByteCounter{}
	splitFunc := counter.Wrap(bufio.ScanLines)
	scanner.Split(splitFunc)

scannerloop:
	for scanner.Scan() {
		lineOrg := scanner.Text()
		line := strings.TrimSpace(lineOrg)

		if len(strings.TrimSpace(line)) == 0 {
			clocFile.Blanks++
			continue
		}

		// shebang line is 'code'
		if isFirstLine && strings.HasPrefix(line, "#!") {
			clocFile.Code++
			isFirstLine = false
			continue
		}

		if len(inComments) == 0 {
			if isFirstLine {
				line = trimBOM(line)
			}

		singleloop:
			for _, singleComment := range language.lineComments {
				if strings.HasPrefix(line, singleComment) {
					// check if single comment is a prefix of multi comment
					for _, ml := range language.multiLines {
						if ml[0] != "" && strings.HasPrefix(line, ml[0]) {
							break singleloop
						}
					}
					clocFile.Comments++
					continue scannerloop
				}
			}

			if len(language.multiLines) == 0 {
				clocFile.Code++
				continue scannerloop
			}
		}

		if len(inComments) == 0 && !containsComment(line, language.multiLines) {
			clocFile.Code++
			continue scannerloop
		}

		lenLine := len(line)
		if len(language.multiLines) == 1 && len(language.multiLines[0]) == 2 && language.multiLines[0][0] == "" {
			clocFile.Code++
			continue
		}
		codeFlags := make([]bool, len(language.multiLines))
		for pos := 0; pos < lenLine; {
			for idx, ml := range language.multiLines {
				begin, end := ml[0], ml[1]
				lenBegin := len(begin)

				if pos+lenBegin <= lenLine && strings.HasPrefix(line[pos:], begin) && (begin != end || len(inComments) == 0) {
					pos += lenBegin
					inComments = append(inComments, [2]string{begin, end})
					continue
				}

				if n := len(inComments); n > 0 {
					last := inComments[n-1]
					if pos+len(last[1]) <= lenLine && strings.HasPrefix(line[pos:], last[1]) {
						inComments = inComments[:n-1]
						pos += len(last[1])
					}
				} else if pos < lenLine && !unicode.IsSpace(nextRune(line[pos:])) {
					codeFlags[idx] = true
				}
			}
			pos++
		}

		isCode := true
		for _, b := range codeFlags {
			if !b {
				isCode = false
			}
		}

		if isCode {
			clocFile.Code++
		} else {
			clocFile.Comments++
		}
	}

	clocFile.Bytes = counter.BytesRead

	return clocFile
}
