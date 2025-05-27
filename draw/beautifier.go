package draw

import (
	"image/color"
	"github.com/jcocozza/deck/parser"
	"regexp"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	regexp_link = regexp.MustCompile(`https?://[^\s]+`)
)

func Color(t parser.ContentType, level int, theme Theme) color.Color {
	switch t {
	case parser.Header:
		switch level {
		case 1:
			return theme.Heading
		case 2:
			return theme.SubHeading
		default:
			return theme.SubHeading
		}
	case parser.File: // syntax highlighting some day
		fallthrough
	case parser.List:
		fallthrough
	case parser.Text:
		return theme.Text
	default:
		panic("invaid content type")
	}
}

type Beautifier interface {
	Beautify(text []string, t parser.ContentType, level int, theme Theme) []ImageItem
	Face(size int) font.Face
}

func NewBeutifier(f *opentype.Font) Beautifier {
	return &BeautifierImpl{f: f}
}

type textColor struct {
	Text  string
	Color color.Color
}

func tokenize(str string, re *regexp.Regexp, matchCol color.Color, defaultCol color.Color) []textColor {
	var tokens []textColor
	lastIndex := 0
	matches := re.FindAllStringIndex(str, -1) // [[start, end], ...]
	for _, match := range matches {
		start, end := match[0], match[1]

		// Add the non-matching text before this match
		if start > lastIndex {
			tokens = append(tokens, textColor{
				Text:  str[lastIndex:start],
				Color: defaultCol, // or a default color
			})
		}

		// Add the matching colored text
		tokens = append(tokens, textColor{
			Text:  str[start:end],
			Color: matchCol,
		})

		lastIndex = end
	}
	// Add any remaining text after the last match
	if lastIndex < len(str) {
		tokens = append(tokens, textColor{
			Text:  str[lastIndex:],
			Color: defaultCol,
		})
	}
	return tokens
}

type BeautifierImpl struct {
	f *opentype.Font
}

func (b *BeautifierImpl) Face(size int) font.Face {
	face, err := FontFace(b.f, size)
	if err != nil {
		panic(err)
	}
	return face
}

func (b *BeautifierImpl) Beautify(text []string, t parser.ContentType, level int, theme Theme) []ImageItem {
	items := []ImageItem{}
	var defaultCol color.Color
	var size int // TODO: this needs to be injected somehow
	switch t {
	case parser.Paragraph:
		// nothing happens here, a paragraph should never have text
	case parser.Header:
		switch level {
		case 1:
			defaultCol = theme.Heading
			size = 18
		case 2:
			defaultCol = theme.SubHeading
			size = 16
		default:
			defaultCol = theme.SubHeading
			size = 16
		}
	case parser.File: // syntax highlighting some day
		fallthrough
	case parser.List:
		fallthrough
	case parser.Text:
		defaultCol = theme.Text
		size = 14
	default:
		defaultCol = theme.Text
		size = 14
	}
	for _, str := range text {
		tokens := tokenize(str, regexp_link, theme.Link, defaultCol)
		face := b.Face(size)
		for _, tkn := range tokens {
			item := &TextImageItem{Text: tkn.Text, Color: tkn.Color, Face: face}
			items = append(items, item)
		}
		items = append(items, &NewLineItem{Face: face})
	}
	return items
}
