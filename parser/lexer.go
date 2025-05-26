package parser

import (
	"strings"
	"unicode"
)

type linetype int

const (
	title linetype = iota
	subtitle
	comment
	file
	emptyLine
	text
	list
)

var prefixes = [...]string{
	title:    "# ",
	subtitle: "## ",
	comment:  "// ",
	file:     "@",
}

func haslstprefix(line string) bool {
	// Check for unordered list prefixes
	if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "+ ") {
		return true
	}
	// Check for numbered list like "1. ", "23. "
	for i, r := range line {
		if unicode.IsDigit(r) {
			continue
		}
		if r == '.' && i > 0 && i+1 < len(line) && line[i+1] == ' ' {
			return true
		}
		break
	}
	return false
}

type lexline struct {
	t    linetype
	text string
}

type Lexer interface {
	Lex(lines []string) []lexline
}

func NewLexer() Lexer {
	return &LexerImpl{}
}

type LexerImpl struct{}

func (l *LexerImpl) lexln(line string) lexline {
	//line = strings.TrimSpace(line)
	// only trim to the right to support user enabled indentation
	line = strings.TrimRightFunc(line, unicode.IsSpace)
	switch {
	case strings.HasPrefix(line, prefixes[title]):
		return lexline{t: title, text: strings.TrimPrefix(line, prefixes[title])}
	case strings.HasPrefix(line, prefixes[subtitle]):
		return lexline{t: subtitle, text: strings.TrimPrefix(line, prefixes[subtitle])}
	case strings.HasPrefix(line, prefixes[comment]):
		return lexline{t: comment, text: line}
	case strings.HasPrefix(line, prefixes[file]):
		return lexline{t: file, text: line}
	case len(line) == 0:
		return lexline{t: emptyLine, text: line}
	case haslstprefix(line):
		return lexline{t: list, text: line}
	default:
		return lexline{t: text, text: line}
	}
}

func (l *LexerImpl) Lex(lines []string) []lexline {
	rawLines := make([]lexline, len(lines))
	for i, line := range lines {
		rawLines[i] = l.lexln(line)
	}
	return rawLines
}
