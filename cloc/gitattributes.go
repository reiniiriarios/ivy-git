package cloc

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func parseGitAttributes(dir string) map[string]string {
	translations := make(map[string]string)
	ga := filepath.Join(dir, ".gitattributes")
	if _, err := os.Stat(ga); err != nil {
		// ignore err
		return translations
	}
	file, err := os.Open(ga)
	if err != nil {
		// ignore err
		return translations
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r := regexp.MustCompile(`\*\.([a-z0-9]+).*(?:diff|linguist-language)=([a-z0-9]+).*`)
		matches := r.FindStringSubmatch(scanner.Text())
		if len(matches) == 3 {
			translations[strings.ToLower(matches[1])] = strings.ToLower(matches[2])
		}
	}
	// err := scanner.Err()
	// ignore err

	return translations
}
