package format

import (
	"fmt"
	"strings"
	"unicode"
)

type linetype string

const (
	header       linetype = "header"
	subheader    linetype = "sub-header"
	subsubheader linetype = "sub-sub-header"
	comment      linetype = "comment"
	fileCenter   linetype = "file-center"
	fileLeft     linetype = "file-left"
	fileRight    linetype = "file-right"
	fileBottom   linetype = "file-bottom"
	fileTop      linetype = "file-top"
	emptyLine    linetype = "emptyLine"
	emptySlide   linetype = "emptySlide"
	text         linetype = "text"
)

var prefixes = map[linetype]string{
	header:       "# ",
	subheader:    "## ",
	subsubheader: "### ",
	comment:      "// ",
	fileCenter:   "@",
	fileLeft:     "@l:",
	fileRight:    "@r:",
	fileBottom:   "@b:",
	fileTop:      "@t:",
	emptySlide:   "\\",
}

type lexline struct {
	t    linetype
	text string
}

// for debugging
func (l *lexline) String() string {
	return fmt.Sprintf("[%s] %s", prefixes[l.t], l.text)
}

type Lexer interface {
	Lex(lines []string) []lexline
}

func NewLexer() Lexer {
	return &LinesLexer{}
}

type LinesLexer struct{}

func (l *LinesLexer) lexln(line string) lexline {
	line = strings.TrimRightFunc(line, unicode.IsSpace)
	switch {

	case strings.HasPrefix(line, prefixes[subsubheader]):
		return lexline{t: subsubheader, text: strings.TrimPrefix(line, prefixes[subsubheader])}
	case strings.HasPrefix(line, prefixes[subheader]):
		return lexline{t: subheader, text: strings.TrimPrefix(line, prefixes[subheader])}
	case strings.HasPrefix(line, prefixes[header]):
		return lexline{t: header, text: strings.TrimPrefix(line, prefixes[header])}

	case strings.HasPrefix(line, prefixes[comment]):
		return lexline{t: comment, text: line}

	case strings.HasPrefix(line, prefixes[fileLeft]):
		return lexline{t: fileLeft, text: strings.TrimPrefix(line, prefixes[fileLeft])}
	case strings.HasPrefix(line, prefixes[fileRight]):
		return lexline{t: fileRight, text: strings.TrimPrefix(line, prefixes[fileRight])}
	case strings.HasPrefix(line, prefixes[fileBottom]):
		return lexline{t: fileBottom, text: strings.TrimPrefix(line, prefixes[fileBottom])}
	case strings.HasPrefix(line, prefixes[fileTop]):
		return lexline{t: fileTop, text: strings.TrimPrefix(line, prefixes[fileTop])}
	case strings.HasPrefix(line, prefixes[fileCenter]):
		return lexline{t: fileCenter, text: strings.TrimPrefix(line, prefixes[fileCenter])}

	case len(line) == 0:
		return lexline{t: emptyLine, text: line}

	case strings.HasPrefix(line, prefixes[emptySlide]):
		return lexline{t: emptySlide, text: ""}

	default:
		return lexline{t: text, text: line}
	}
}

func (l *LinesLexer) Lex(lines []string) []lexline {
	rawLines := make([]lexline, len(lines))
	for i, line := range lines {
		rawLines[i] = l.lexln(line)
	}
	return rawLines
}
