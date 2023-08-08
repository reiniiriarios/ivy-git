package files

import (
	"fmt"
	"html"
	"sort"

	"github.com/alecthomas/chroma/v2"
)

// Option sets an option of the HTML formatter.
type Option func(f *Formatter)

// HighlightLines higlights the given line ranges with the Highlight style.
//
// A range is the beginning and ending of a range as 1-based line numbers, inclusive.
func HighlightLines(ranges [][2]int) Option {
	return func(f *Formatter) {
		f.highlightRanges = ranges
		sort.Sort(f.highlightRanges)
	}
}

// New formatter.
func LinesFormatter(options ...Option) *Formatter {
	f := &Formatter{}
	for _, option := range options {
		option(f)
	}
	return f
}

type Formatter struct {
	highlightRanges highlightRanges
}

type highlightRanges [][2]int

func (h highlightRanges) Len() int           { return len(h) }
func (h highlightRanges) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h highlightRanges) Less(i, j int) bool { return h[i][0] < h[j][0] }

func (f *Formatter) Format(style *chroma.Style, iterator chroma.Iterator) map[int]string {
	iterator_tokens := iterator.Tokens()
	h_lines := make(map[int]string)
	line_split := chroma.SplitTokensIntoLines(iterator_tokens)
	highlightIndex := 0

	for index, tokens := range line_split {
		// 1-based line number.
		line_no := 1 + index
		highlight, next := f.shouldHighlight(highlightIndex, line_no)
		if next {
			highlightIndex++
		}

		// Start of Line
		line_highlight := ""
		line_highlight += `<span`

		if highlight {
			// Line + LineHighlight
			line_highlight += fmt.Sprintf(` class="%s %s"`, f.class(chroma.Line), f.class(chroma.LineHighlight))
			line_highlight += `>`
		} else {
			line_highlight += fmt.Sprintf("%s>", f.classAttr(chroma.Line))
		}

		line_highlight += fmt.Sprintf(`<span%s>`, f.classAttr(chroma.CodeLine))

		for _, token := range tokens {
			html := html.EscapeString(token.String())
			attr := f.classAttr(token.Type)
			if attr != "" {
				html = fmt.Sprintf("<span%s>%s</span>", attr, html)
			}
			line_highlight += html
		}

		line_highlight += `</span>` // End of CodeLine
		line_highlight += `</span>` // End of Line

		h_lines[line_no] = line_highlight
	}

	return h_lines
}

func (f *Formatter) shouldHighlight(highlightIndex, line int) (bool, bool) {
	next := false
	for highlightIndex < len(f.highlightRanges) && line > f.highlightRanges[highlightIndex][1] {
		highlightIndex++
		next = true
	}
	if highlightIndex < len(f.highlightRanges) {
		hrange := f.highlightRanges[highlightIndex]
		if line >= hrange[0] && line <= hrange[1] {
			return true, next
		}
	}
	return false, next
}

func (f *Formatter) class(t chroma.TokenType) string {
	for t != 0 {
		if cls, ok := chroma.StandardTypes[t]; ok {
			return cls
		}
		t = t.Parent()
	}
	return chroma.StandardTypes[t]
}

func (f *Formatter) classAttr(tt chroma.TokenType) string {
	cls := f.class(tt)
	if cls == "" {
		return ""
	}
	return fmt.Sprintf(` class="%s"`, cls)
}
