package git

import (
	"bufio"
	"os"
	"path/filepath"
)

func (g *Git) deleteLinesFromFile(filename string, lines []int64) error {
	file, err := os.Open(filepath.Join(g.Repo.Directory, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	tmpfile, err := os.Open(filepath.Join(g.Repo.Directory, filename+"~~"))
	if err != nil {
		return err
	}
	defer tmpfile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(tmpfile)
	var line_counter int64 = 0
	delete_counter := 0
	for scanner.Scan() {
		line_counter++
		if line_counter == lines[delete_counter] {
			delete_counter++
			if delete_counter >= len(lines) {
				break
			}
		} else {
			writer.WriteString(scanner.Text())
		}
	}

	return nil
}
