package git

import (
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

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

// Parse a single line result.
func parseOneLine(s string) string {
	return strings.Trim(strings.Trim(strings.ReplaceAll(s, "\r\n", "\n"), "\n"), "'")
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

func mdToHTML(md string) string {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	mdb := []byte(md)
	doc := p.Parse(mdb)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// render
	maybeUnsafeHTML := markdown.Render(doc, renderer)

	// safe
	html := bluemonday.UGCPolicy().SanitizeBytes(maybeUnsafeHTML)

	return string(html)
}
