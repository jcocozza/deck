package draw

import (
	"image/color"
	"regexp"
	"sort"

	"github.com/jcocozza/deck/internal/slide"
)

type ThemeElement struct {
	Size  int
	Color color.Color
}

type Theme struct {
	Background   color.Color
	Header       ThemeElement
	SubHeader    ThemeElement
	SubSubHeader ThemeElement
	Link         ThemeElement
	Default      ThemeElement
}

// White background, all text is black
var DefaultTheme = Theme{
	Background:   color.RGBA{R: 255, G: 255, B: 255, A: 255}, // white
	Header:       ThemeElement{Size: 18, Color: color.Black},  // blue
	SubHeader:    ThemeElement{Size: 16, Color: color.Black},  // orange
	SubSubHeader: ThemeElement{Size: 16, Color: color.Black},  // orange
	Link:         ThemeElement{Size: 14, Color: color.Black},    // hyperlink blue
	Default:      ThemeElement{Size: 14, Color: color.Black},       // black
}

// TODO: make this decent; it is dumb right now
var DefaultColorTheme = Theme{
	Background:   color.RGBA{R: 255, G: 255, B: 255, A: 255}, // white
	Header:       ThemeElement{Size: 18, Color: color.RGBA{R: 68, G: 114, B: 196, A: 255}},  // blue
	SubHeader:    ThemeElement{Size: 16, Color: color.RGBA{R: 237, G: 125, B: 49, A: 255}},  // orange
	SubSubHeader: ThemeElement{Size: 16, Color: color.RGBA{R: 237, G: 125, B: 49, A: 255}},  // orange
	Link:         ThemeElement{Size: 14, Color: color.RGBA{R: 5, G: 99, B: 193, A: 255}},    // hyperlink blue
	Default:      ThemeElement{Size: 14, Color: color.RGBA{R: 0, G: 0, B: 0, A: 255}},       // black
}

func (t *Theme) GetElement(ty slide.SlideLineType) ThemeElement {
	switch ty {
	case slide.Header:
		return t.Header
	case slide.Subheader:
		return t.SubHeader
	case slide.Subsubheader:
		return t.SubSubHeader
	default:
		return t.Default
	}
}

func (t *Theme) GetElementByPattern(ty patternType) ThemeElement {
	switch ty {
	case link:
		return t.Link
	default:
		return t.Default
	}
}

type prettystring struct {
	Text    string
	T       ThemeElement
	Bold    bool
	Italics bool
}

type patternType int

const (
	defaultpt patternType = iota
	link
)

var (
	patterns = map[patternType]*regexp.Regexp{
		link: regexp.MustCompile(`https?://[^\s]+`),
	}
)

type match struct {
	Start int
	End   int
	T     ThemeElement
}

func tokenizeLine(line string, theme Theme) []prettystring {
	var matches []match
	for pttrnT, re := range patterns {
		locs := re.FindAllStringIndex(line, -1)
		for _, loc := range locs {
			matches = append(matches, match{
				Start: loc[0],
				End:   loc[1],
				T:     theme.GetElementByPattern(pttrnT),
			})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Start < matches[j].Start
	})
	// remove overlaps (keep first)
	var filtered []match
	lastEnd := -1
	for _, m := range matches {
		if m.Start >= lastEnd {
			filtered = append(filtered, m)
			lastEnd = m.End
		}
	}
	var tokens []prettystring
	pos := 0
	for _, m := range filtered {
		if pos < m.Start {
			tokens = append(tokens, prettystring{
				Text: line[pos:m.Start],
				T:    theme.GetElementByPattern(defaultpt),
			})
		}
		tokens = append(tokens, prettystring{
			Text: line[m.Start:m.End],
			T:    m.T,
		})
		pos = m.End
	}
	if pos < len(line) {
		tokens = append(tokens, prettystring{
			Text: line[pos:],
			T:    theme.GetElementByPattern(defaultpt),
		})
	}
	return tokens
}

func MakePretty(line slide.SlideLine, theme Theme) []prettystring {
	switch line.T {
	case slide.Header, slide.Subheader, slide.Subsubheader:
		return []prettystring{{Text: line.Text, T: theme.GetElement(line.T)}}
	default:
		return tokenizeLine(line.Text, theme)
	}
}
