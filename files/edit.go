package files

import (
	"bufio"
	"io"
	"os"
	"sort"
)

func DeleteAndReplaceLinesFromFile(filename string, delete_lines []int64, replace_lines map[int64]int64) error {
	// Open the source file.
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Flip the replace map to lookup by replace line.
	replace_lines_reverse := make(map[int64]int64, len(replace_lines))
	for k, v := range replace_lines {
		replace_lines_reverse[v] = k
	}

	// Create a map of contents by line number => string to replace with.
	replace_lines_contents := make(map[int64]string)
	if len(replace_lines) > 0 {
		scanner := bufio.NewScanner(file)
		var rep_ln int64
		for rep_ln = 1; scanner.Scan(); rep_ln++ {
			if ln, exists := replace_lines_reverse[rep_ln]; exists {
				rep_text := scanner.Text()
				replace_lines_contents[ln] = rep_text
			}

			if err != nil {
				break
			}
		}
	}

	// Create a temp file to write to.
	tmpfile, err := os.Create(filename + "~")
	if err != nil {
		return err
	}
	defer tmpfile.Close()

	sort.Slice(delete_lines, func(i, j int) bool { return delete_lines[i] < delete_lines[j] })

	// Move back to beginning of file.
	file.Seek(0, io.SeekStart)
	// Create reader to take advantage of reading line endings.
	reader := bufio.NewReader(file)
	var line_counter int64 = 0
	delete_counter := 0
	for {
		line, err := reader.ReadString('\n')
		if err == nil || err == io.EOF {
			line_counter++
			// Delete and replace lines, writing into the temp file.
			if delete_counter < len(delete_lines) && line_counter == delete_lines[delete_counter] {
				// Do nothing.
				delete_counter++
			} else if replacement, exists := replace_lines_contents[line_counter]; exists {
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
			if err != io.EOF {
				println("Error reading file:", err)
			}
			break
		}
	}

	// Delete the original.
	err = os.Remove(filename)
	if err != nil {
		return err
	}

	// Replace the original with the temp file.
	err = os.Rename(filename+"~", filename)
	if err != nil {
		return err
	}

	return nil
}
