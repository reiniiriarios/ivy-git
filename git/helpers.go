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

func getSiteName(hostname string) string {
	if hostname == "github.com" {
		return "GitHub"
	}
	if hostname == "bitbucket.org" {
		return "Bitbucket"
	}
	if hostname == "gitlab.com" {
		return "GitLab"
	}
	if hostname == "dev.azure.com" {
		return "Azure"
	}
	return hostname
}
