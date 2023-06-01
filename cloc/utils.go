package cloc

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
)

func trimBOM(line string) string {
	l := len(line)
	if l >= 3 {
		if line[0] == 0xef && line[1] == 0xbb && line[2] == 0xbf {
			trimLine := line[3:]
			return trimLine
		}
	}
	return line
}

func containsComment(line string, multiLines [][]string) bool {
	for _, mlcomm := range multiLines {
		for _, comm := range mlcomm {
			if strings.Contains(line, comm) {
				return true
			}
		}
	}
	return false
}

func nextRune(s string) rune {
	for _, r := range s {
		return r
	}
	return 0
}

func checkMD5Sum(path string, fileCache map[string]struct{}) (ignore bool) {
	content, err := os.ReadFile(path)
	if err != nil {
		return true
	}

	// calc md5sum
	hash := md5.Sum(content)
	c := fmt.Sprintf("%x", hash)
	if _, ok := fileCache[c]; ok {
		return true
	}

	fileCache[c] = struct{}{}
	return false
}

func ignoreFile(file string) bool {
	// todo: more of this
	return strings.HasSuffix(file, "package-lock.json")
}

func parseAllFiles(files []string, languages *DefinedLanguages, translations map[string]string) (result map[string]*Language) {
	result = make(map[string]*Language, 0)
	fileCache := make(map[string]struct{})

	for _, file := range files {
		if ignoreFile(file) {
			continue
		}
		ext, ok := getFileType(file)
		if !ok {
			continue
		}
		// Swap with language defined in gitattributes if present.
		if translation, exists := translations[ext]; exists {
			ext = translation
		}
		targetExt, ok := Exts[ext]
		if !ok {
			continue
		}
		if checkMD5Sum(file, fileCache) {
			continue
		}
		if _, exists := result[targetExt]; !exists {
			result[targetExt] = NewLanguage(
				languages.Langs[targetExt].Data.Name,
				languages.Langs[targetExt].lineComments,
				languages.Langs[targetExt].multiLines)
		}
		result[targetExt].Files = append(result[targetExt].Files, file)
	}
	return result
}
