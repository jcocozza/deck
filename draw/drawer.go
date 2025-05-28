package draw

import (
	"image"
	"image/color"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/jcocozza/deck/slide"
	"golang.org/x/image/font/opentype"
)

// based on a particular font, determine the max font size the text can can be
func scaleText(maxWidth int, maxHeight int, lines []ImageElement) (float64, error) {
	fontSize := 1.0

	if len(lines) == 1 {
		switch lines[0].(type) {
		case *EmbeddedImageElement:
			return 10, nil
		}
	}

	for {
		largestWidth := 0
		totalHeight := 0
		for _, line := range lines {
			width := line.Width(fontSize)
			if width > largestWidth {
				largestWidth = width
			}
			totalHeight += line.Height(fontSize)
		}
		if largestWidth > maxWidth || totalHeight > maxHeight {
			break
		}
		fontSize++
	}
	return fontSize, nil
}

// TODO: this is probably where we inject the theme
func slideLineToElement(line slide.Line) (ImageElement, error) {
	c := color.Black
	size := 12.0
	fnt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		return nil, err
	}
	switch line.T {
	case slide.Header:
		return NewTextElement(line.Content, c, size, fnt)
	case slide.SubHeader:
		return NewTextElement(line.Content, c, size, fnt)
	case slide.SubSubHeader:
		return NewTextElement(line.Content, c, size, fnt)
	case slide.Text:
		return NewTextElement(line.Content, c, size, fnt)
	case slide.File:
		switch filepath.Ext(line.Content) {
		default:
			fnt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
			if err != nil {
				return nil, err
			}
			return NewEmbededTextFileElement(line.Content, 12, fnt)
		}
	}
	panic("unexpected slide line type")
}

type SlideDrawer interface {
	Draw(img *image.RGBA, paddingX int, paddingY int, s slide.Slide) error
}

type Drawer struct {
	fnt *opentype.Font
}

func (d *Drawer) Draw(img *image.RGBA, paddingX int, paddingY int, s slide.Slide) error {
	elms := []ImageElement{}
	for _, line := range s.Lines {
		elm, err := slideLineToElement(line)
		if err != nil {
			return err
		}
		elms = append(elms, elm)
	}

	imgSpaceY := 0
	if s.HasImg() {
		imgSpaceY += s.Img.Image.Bounds().Size().Y
	}
	pt := img.Bounds().Size()
	fntSize, err := scaleText(pt.X-2*paddingX, pt.Y-2*paddingY-imgSpaceY, elms)
	if err != nil {
		return err
	}

	face, _ := opentype.NewFace(d.fnt, &opentype.FaceOptions{Size: fntSize, DPI: 72})
	txtHeight := face.Metrics().Height.Ceil()
	ascent := face.Metrics().Ascent.Ceil()
	totalTextHeight := txtHeight * len(elms)
	startY := paddingY + (pt.Y-2*paddingY-totalTextHeight-imgSpaceY)/2
	baseline := startY + ascent

	for _, elm := range elms {
		// handle resizing
		width := elm.Width(fntSize)
		x := paddingX + (pt.X-2*paddingX-width)/2
		elm.Draw(img, fntSize, x, baseline)
		baseline += elm.Height(fntSize)
	}

	if s.HasImg() {
		e := EmbeddedImageElement{Img: s.Img.Image}
		e.Draw(img, fntSize, 0,baseline)
	}
	return nil
}

type PrettyDrawer struct {
	fnts *opentype.Collection
	fnt  *opentype.Font
}

func (d *PrettyDrawer) Draw(img *image.RGBA, paddingX int, paddingY int, s slide.Slide) error {
	return nil
}

func DrawerFactory(pretty bool, fnt *opentype.Font) SlideDrawer {
	if pretty {
		return &PrettyDrawer{fnt: fnt}
	}
	return &Drawer{fnt: fnt}
}
