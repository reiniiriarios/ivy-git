package git

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func (g *Git) deleteAndReplaceLinesFromFile(filename string, delete_lines []int64, replace_lines map[int64]string) error {
	file_path := filepath.Join(g.Repo.Directory, filename)
	// Open the source file.
	file, err := os.Open(file_path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a temp file to write to.
	tmpfile, err := os.Create(file_path + "~")
	if err != nil {
		return err
	}
	defer tmpfile.Close()

	sort.Slice(delete_lines, func(i, j int) bool { return delete_lines[i] < delete_lines[j] })

	// Delete and replace lines, writing into the temp file.
	reader := bufio.NewReader(file)
	var line_counter int64 = 0
	delete_counter := 0
	for {
		line, err := reader.ReadString('\n')
		if err == nil || err == io.EOF {
			line_counter++
			if delete_counter < len(delete_lines) && line_counter == delete_lines[delete_counter] {
				// Do nothing.
				delete_counter++
			} else if replacement, exists := replace_lines[line_counter]; exists {
				var crlf string
				if len(line) >= 2 && line[len(line)-2:len(line)-1] == "\r" {
					crlf = "\r\n"
				} else {
					crlf = "\n"
				}
				tmpfile.WriteString(replacement + crlf)
			} else {
				tmpfile.WriteString(line)
			}
		}
		if err != nil {
			break
		}
	}

	// Delete the original.
	err = os.Remove(file_path)
	if err != nil {
		return err
	}

	// And replace the original with the temp file.
	err = os.Rename(file_path+"~", file_path)
	if err != nil {
		return err
	}

	return nil
}
