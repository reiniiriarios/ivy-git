package git

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func (g *Git) parseGitAttributes() map[string]string {
	translations := make(map[string]string)
	if g.Repo.Main == "" {
		return translations
	}

	ga := filepath.Join(g.Repo.Main, ".gitattributes")
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
