package git

import "strings"

// Parse a multiline string into an array.
func parseLines(s string) []string {
	var r []string
	l := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	for _, v := range l {
		v = strings.Trim(v, "'")
		r = append(r, v)
	}
	return r
}
