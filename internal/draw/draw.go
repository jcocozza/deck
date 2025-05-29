package draw

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/jcocozza/deck/internal/slide"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func fillImg(img *image.RGBA, c color.Color) {
	draw.Draw(img, img.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
}

func newImage(width int, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return img
}

func readTestFont() *opentype.Font {
	f, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		panic(err)
	}
	return f
}

func GenerateSlideImage(s slide.Slide, screenWidth int, screenHeight int, paddingX int, paddingY int, fnt *opentype.Font) (image.Image, error) {
	// THIS IS TMP
	fnt = readTestFont()
	canvas := newImage(screenWidth, screenHeight)
	fillImg(canvas, color.White)

	imageOffsetX, imageOffsetY := 0, 0
	if s.Image != nil {
		pt := s.Image.I.Bounds().Size()
		switch s.Image.Position {
		case slide.Left, slide.Right:
			imageOffsetX = pt.X
		case slide.Bottom, slide.Top:
			imageOffsetY = pt.Y
		default:
			panic("invalid state")
		}
	}
	maxWidth := screenWidth - 2*paddingX - imageOffsetX
	maxHeight := screenHeight - 2*paddingY - imageOffsetY
	fntSize, err := scaleText(maxWidth, maxHeight, fnt, s.Lines)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{Size: fntSize, DPI: 72})
	if err != nil {
		return nil, err
	}
	defer face.Close()

	txtHeight := face.Metrics().Height.Ceil()
	ascent := face.Metrics().Ascent.Ceil()
	totalTextHeight := txtHeight * len(s.Lines)
	textStartX, textStartY := paddingX, paddingY

	if s.Image != nil {
		switch s.Image.Position {
		case slide.Left:
			// Image on left, text on right
			imageX := 0
			imageY := (screenHeight - s.Image.I.Bounds().Dy()) / 2
			drawImage(canvas, s.Image.I, imageX, imageY)
			textStartX = paddingX + s.Image.I.Bounds().Dx()
			textStartY = paddingY + (screenHeight-2*paddingY-totalTextHeight)/2
		case slide.Right:
			// Image on right, text on left
			imageX := screenWidth - s.Image.I.Bounds().Dx()
			imageY := (screenHeight - s.Image.I.Bounds().Dy()) / 2
			drawImage(canvas, s.Image.I, imageX, imageY)
			textStartX = paddingX
			textStartY = paddingY + (screenHeight-2*paddingY-totalTextHeight)/2
		case slide.Top:
			// Image on top, text on bottom
			imageX := (screenWidth - s.Image.I.Bounds().Dx()) / 2
			imageY := 0
			drawImage(canvas, s.Image.I, imageX, imageY)
			textStartY = paddingY + s.Image.I.Bounds().Dy() + (screenHeight-2*paddingY-s.Image.I.Bounds().Dy()-totalTextHeight)/2
		case slide.Bottom:
			// Image on bottom, text on top
			imageX := (screenWidth - s.Image.I.Bounds().Dx()) / 2
			imageY := screenHeight - s.Image.I.Bounds().Dy()
			drawImage(canvas, s.Image.I, imageX, imageY)
			textStartY = paddingY + (screenHeight-2*paddingY-s.Image.I.Bounds().Dy()-totalTextHeight)/2
		default:
			panic("invalid state")
		}
	} else {
		// No image: center text in full canvas
		textStartY = paddingY + (screenHeight-2*paddingY-totalTextHeight)/2
	}
	baseline := textStartY + ascent
	for _, line := range s.Lines {
		lineWidth := font.MeasureString(face, line).Ceil()
		textAreaWidth := screenWidth - 2*paddingX - imageOffsetX
		x := textStartX + (textAreaWidth-lineWidth)/2
		drawText(canvas, face, color.Black, x, baseline, line)
		baseline += txtHeight
	}
	return canvas, nil
}
