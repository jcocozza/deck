package draw

import (
	"image/color"
	"presentation/parser"
	"regexp"
)

var (
	regexp_link = regexp.MustCompile(`https?://[^\s]+`)
)

type Theme struct {
	Background     color.Color // lt1
	Foreground     color.Color // dk1
	Heading        color.Color // dk2 or accent1
	SubHeading     color.Color // accent2
	Text           color.Color // dk1
	CardBackground color.Color // lt2
	Link           color.Color // hlink
	VisitedLink    color.Color // folHlink
	Primary        color.Color // accent1
	Secondary      color.Color // accent3 or accent5
}

var DefaultTheme = Theme{
	Background:     color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White (lt1)
	Foreground:     color.RGBA{R: 0, G: 0, B: 0, A: 255},       // Black (dk1)
	Heading:        color.RGBA{R: 68, G: 114, B: 196, A: 255},  // Accent1 (blue)
	SubHeading:     color.RGBA{R: 237, G: 125, B: 49, A: 255},  // Accent2 (orange)
	Text:           color.RGBA{R: 0, G: 0, B: 0, A: 255},       // Same as Foreground
	CardBackground: color.RGBA{R: 242, G: 242, B: 242, A: 255}, // Light gray (lt2)
	Link:           color.RGBA{R: 5, G: 99, B: 193, A: 255},    // Hyperlink blue
	VisitedLink:    color.RGBA{R: 149, G: 79, B: 114, A: 255},  // Visited purple
	Primary:        color.RGBA{R: 68, G: 114, B: 196, A: 255},  // Accent1 again
	Secondary:      color.RGBA{R: 165, G: 165, B: 165, A: 255}, // Neutral gray (accent5-like)
}

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
}

func NewBeutifier() Beautifier {
	return &BeautifierImpl{}
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

type BeautifierImpl struct{}

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
		tokens := tokenize(str, regexp_link, theme.Link, defaultCol) // TODO: sizes should NOT be hardcoded
		for _, tkn := range tokens {
			face, err := FontFace(readTestFont(), size)
			if err != nil {panic(err)}
			items = append(items, &TextImageItem{Text: tkn.Text, Color: tkn.Color, Face: face})
		}
	}
	return items
}
