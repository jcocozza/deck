package format

import (
	"fmt"

	"github.com/jcocozza/deck/internal/slide"
)

type Parser interface {
	Parse(lines []lexline) ([]slide.Slide, error)
}

func NewParser() Parser {
	return &LineParser{}
}

type LineParser struct{}

func (p *LineParser) Parse(lines []lexline) ([]slide.Slide, error) {
	slides := []slide.Slide{}
	var curr *slide.Slide
	var lastLineType linetype

	for _, line := range lines {
		if curr == nil {
			curr = &slide.Slide{}
		}
		switch line.t {
		case header:
			ln := slide.SlideLine{Text: line.text, T: slide.Header}
			curr.Lines = append(curr.Lines, ln)
		case subheader:
			ln := slide.SlideLine{Text: line.text, T: slide.Subheader}
			curr.Lines = append(curr.Lines, ln)
		case subsubheader:
			ln := slide.SlideLine{Text: line.text, T: slide.Subsubheader}
			curr.Lines = append(curr.Lines, ln)
		case comment:
			continue // we ignore these
		case emptyLine:
			if lastLineType == emptyLine { // ignore consecutive empty lines
				continue
			}
			slides = append(slides, *curr)
			curr = nil
		case emptySlide:
			slides = append(slides, *curr)
			curr = nil
		case text:
			ln := slide.SlideLine{Text: line.text, T: slide.Text}
			curr.Lines = append(curr.Lines, ln)
		case fileTop, fileBottom, fileLeft,fileRight, fileCenter:
			var pos slide.ImgPostion
			switch line.t {
			case fileTop:
				pos = slide.Top
			case fileBottom:
				pos = slide.Bottom
			case fileLeft:
				pos = slide.Left
			case fileRight:
				pos = slide.Right
			case fileCenter:
				pos = slide.Center
			}
			img, err := slide.NewImage(line.text, pos)
			if err != nil {
				return nil, err
			}
			curr.Image = img
		default:
			return nil, fmt.Errorf("invalid state")
		}
		lastLineType = line.t
	}
	if curr != nil {
		slides = append(slides, *curr)
	}
	return slides, nil
}
