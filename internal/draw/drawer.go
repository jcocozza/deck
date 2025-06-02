package draw

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// based on a particular font and size constraints, determine the max font size the text can can be
//
// TODO: we could probably make this a binary search and make it way faster
func scaleText(maxWidth int, maxHeight int, fnt *opentype.Font, lines []string) (float64, error) {
	if len(lines) == 0 {
		return 10, nil
	}
	fontSize := 1.0
	for {
		face, err := opentype.NewFace(fnt, &opentype.FaceOptions{Size: fontSize, DPI: 72})
		if err != nil {
			return -1, err
		}
		largestWidth := 0
		totalHeight := 0

		lineHeight := face.Metrics().Height.Ceil()
		for _, line := range lines {
			width := font.MeasureString(face, line).Ceil()
			if width > largestWidth {
				largestWidth = width
			}
			totalHeight += lineHeight
		}
		err = face.Close()
		if err != nil {
			return -1 , err
		}
		if largestWidth > maxWidth || totalHeight > maxHeight {
			break
		}
		fontSize++
	}
	return fontSize-1, nil
}

// draw text on the passed image
func drawText(img *image.RGBA, face font.Face, c color.Color, x int, y int, text string) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(text)
}

func drawImage(canvas draw.Image, img image.Image, x int, y int) {
	offset := image.Pt(x, y)
	bounds := img.Bounds().Add(offset)
	draw.Draw(canvas, bounds, img, image.Point{}, draw.Over)
}
