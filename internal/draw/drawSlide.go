package draw

import (
	"image"
	"image/draw"

	"github.com/jcocozza/deck/internal/slide"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type DrawerType int

const (
	Auto DrawerType = iota
	Pretty
)

type Drawer interface {
	DrawSlide(s slide.Slide, screenWidth int, screenHeight int, paddingX int, paddingY int, fnt *opentype.Font) (image.Image, error)
}

func NewDrawer(t DrawerType, theme Theme) Drawer {
	switch t {
	case Pretty:
		return &PrettyDrawer{theme: theme}
	default:
		return &AutoDrawer{theme: theme}
	}
}

// returns adjusted textStartX, textStartY
//
// scaled for text height
func drawImageOnCanvasAutoScaled(
	canvas draw.Image,
	img image.Image,
	pos slide.ImgPostion,
	screenWidth int,
	screenHeight int,
	textStartX int,
	textStartY int,
	totalTextHeight int,
	paddingX int,
	paddingY int,
) (int, int) {
	switch pos {
	case slide.Center:
		imageX := (screenWidth - img.Bounds().Dx()) / 2
		imageY := (screenHeight - img.Bounds().Dy()) / 2
		drawImage(canvas, img, imageX, imageY)
	case slide.Left:
		// Image on left, text on right
		imageX := 0
		imageY := (screenHeight - img.Bounds().Dy()) / 2
		drawImage(canvas, img, imageX, imageY)
		textStartX = paddingX + img.Bounds().Dx()
		textStartY = paddingY + (screenHeight-2*paddingY-totalTextHeight)/2
	case slide.Right:
		// Image on right, text on left
		imageX := screenWidth - img.Bounds().Dx()
		imageY := (screenHeight - img.Bounds().Dy()) / 2
		drawImage(canvas, img, imageX, imageY)
		textStartX = paddingX
		textStartY = paddingY + (screenHeight-2*paddingY-totalTextHeight)/2
	case slide.Top:
		// Image on top, text on bottom
		imageX := (screenWidth - img.Bounds().Dx()) / 2
		imageY := 0
		drawImage(canvas, img, imageX, imageY)
		textStartY = paddingY + img.Bounds().Dy() + (screenHeight-2*paddingY-img.Bounds().Dy()-totalTextHeight)/2
	case slide.Bottom:
		// Image on bottom, text on top
		imageX := (screenWidth - img.Bounds().Dx()) / 2
		imageY := screenHeight - img.Bounds().Dy()
		drawImage(canvas, img, imageX, imageY)
		textStartY = paddingY + (screenHeight-2*paddingY-img.Bounds().Dy()-totalTextHeight)/2
	}
	return textStartX, textStartY
}

// non scaled for text height
func drawImageOnCanvas(
	canvas draw.Image,
	img image.Image,
	pos slide.ImgPostion,
	screenWidth int,
	screenHeight int,
	textStartX int,
	textStartY int,
	paddingX int,
	paddingY int,
) (int, int) {
	switch pos {
	case slide.Center:
		imageX := (screenWidth - img.Bounds().Dx()) / 2
		imageY := (screenHeight - img.Bounds().Dy()) / 2
		drawImage(canvas, img, imageX, imageY)
	case slide.Left:
		// Image on left, text on right
		imageX := 0
		imageY := (screenHeight - img.Bounds().Dy()) / 2
		drawImage(canvas, img, imageX, imageY)
		textStartX = paddingX + img.Bounds().Dx()
		textStartY = paddingY + (screenHeight-2*paddingY)/2
	case slide.Right:
		// Image on right, text on left
		imageX := screenWidth - img.Bounds().Dx()
		imageY := (screenHeight - img.Bounds().Dy()) / 2
		drawImage(canvas, img, imageX, imageY)
		textStartX = paddingX
		textStartY = paddingY + (screenHeight-2*paddingY)/2
	case slide.Top:
		// Image on top, text on bottom
		imageX := (screenWidth - img.Bounds().Dx()) / 2
		imageY := 0
		drawImage(canvas, img, imageX, imageY)
		textStartY = paddingY + img.Bounds().Dy() + (screenHeight-2*paddingY-img.Bounds().Dy())/2
	case slide.Bottom:
		// Image on bottom, text on top
		imageX := (screenWidth - img.Bounds().Dx()) / 2
		imageY := screenHeight - img.Bounds().Dy()
		drawImage(canvas, img, imageX, imageY)
		textStartY = paddingY + (screenHeight-2*paddingY-img.Bounds().Dy())/2
	}
	return textStartX, textStartY
}

// the auto drawer will consider a theme, but only colors.
//
// all text is scaled to 1 size fits all
type AutoDrawer struct{ theme Theme }

func (d *AutoDrawer) DrawSlide(s slide.Slide, screenWidth int, screenHeight int, paddingX int, paddingY int, fnt *opentype.Font) (image.Image, error) {
	// THIS IS TMP
	fnt = readTestFont()
	canvas := newImage(screenWidth, screenHeight)
	fillImg(canvas, d.theme.Background)
	imageOffsetX, imageOffsetY := 0, 0
	if s.Image != nil {
		pt := s.Image.I.Bounds().Size()
		switch s.Image.Position {
		case slide.Left, slide.Right:
			imageOffsetX = pt.X
		case slide.Bottom, slide.Top:
			imageOffsetY = pt.Y
		case slide.Center:
			imageOffsetX = pt.X
			imageOffsetY = pt.Y
		default:
			panic("invalid state")
		}
	}
	maxWidth := screenWidth - 2*paddingX - imageOffsetX
	maxHeight := screenHeight - 2*paddingY - imageOffsetY
	strLines := []string{}
	for _, ln := range s.Lines {
		strLines = append(strLines, ln.Text)
	}
	fntSize, err := scaleText(maxWidth, maxHeight, fnt, strLines)
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
		textStartX, textStartY = drawImageOnCanvasAutoScaled(
			canvas,
			s.Image.I,
			s.Image.Position,
			screenWidth,
			screenHeight,
			textStartX,
			textStartY,
			totalTextHeight,
			paddingX,
			paddingY,
		)
	} else {
		// No image: center text in full canvas
		textStartY = paddingY + (screenHeight-2*paddingY-totalTextHeight)/2
	}
	baseline := textStartY + ascent
	for _, line := range s.Lines {
		lineWidth := font.MeasureString(face, line.Text).Ceil()
		textAreaWidth := screenWidth - 2*paddingX - imageOffsetX
		var x int
		if line.T == slide.ListItem {
			x = textStartX
		} else {
			x = textStartX + (textAreaWidth-lineWidth)/2
		}

		pretties := MakePretty(line, d.theme)
		for _, sub := range pretties {
			drawText(canvas, face, sub.T.Color, x, baseline, sub.Text)
			x += font.MeasureString(face, sub.Text).Ceil()
		}
		baseline += txtHeight
	}
	return canvas, nil
}

// the pretty drawer takes themes into account and will recompute fonts based on them
//
// Like a traditional ppt, the user is responsible for ensuring the font sizes are set correctly
// the pretty drawer will not auto size.
type PrettyDrawer struct{ theme Theme }

func (d *PrettyDrawer) DrawSlide(s slide.Slide, screenWidth int, screenHeight int, paddingX int, paddingY int, fnt *opentype.Font) (image.Image, error) {
	// THIS IS TMP
	fnt = readTestFont()
	canvas := newImage(screenWidth, screenHeight)
	fillImg(canvas, d.theme.Background)

	textStartX, textStartY := paddingX, paddingY
	if s.Image != nil {
		textStartX, textStartY = drawImageOnCanvas(
			canvas,
			s.Image.I,
			s.Image.Position,
			screenWidth,
			screenHeight,
			textStartX,
			textStartY,
			paddingX,
			paddingY,
		)
	}
	y := textStartY
	for _, line := range s.Lines {
		pretties := MakePretty(line, d.theme)
		x := textStartX
		maxTxtHeight := 0
		for _, sub := range pretties {
			face, err := opentype.NewFace(fnt, &opentype.FaceOptions{Size: float64(sub.T.Size), DPI: 72})
			if err != nil {
				return nil, err
			}
			txtHeight := face.Metrics().Height.Ceil()
			drawText(canvas, face, sub.T.Color, x, y, sub.Text)
			x += font.MeasureString(face, sub.Text).Ceil()
			if txtHeight > maxTxtHeight {
				maxTxtHeight = txtHeight
			}
			err = face.Close()
			if err != nil {
				return nil, err
			}
		}
		y += maxTxtHeight
	}
	return canvas, nil
}
