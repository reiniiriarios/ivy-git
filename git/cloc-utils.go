package git

import (
	"bytes"
	"strings"
	"sync"
)

var bsPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func getByteSlice() *bytes.Buffer {
	v := bsPool.Get().(*bytes.Buffer)
	v.Reset()
	return v
}

func putByteSlice(bs *bytes.Buffer) {
	bsPool.Put(bs)
}

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

func clocIgnoreFile(file string) bool {
	// todo: more of this
	return strings.HasSuffix(file, "package-lock.json")
}

func (g *Git) clocAllFiles(files []string, languages *DefinedLanguages, translations map[string]string) (result map[string]*Language) {
	result = make(map[string]*Language, 0)
	// fileCache := make(map[string]struct{})

	for _, file := range files {
		if clocIgnoreFile(file) {
			continue
		}
		ext, ok := g.getFileType(file)
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
