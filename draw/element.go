package draw

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// raw "api" to draw onto an image
type ImageElement interface {
	// draw the object onto the passed image at x,y
	Draw(img *image.RGBA, size float64, x int, y int)
	// with of the object in pixels
	Width(size float64) int
	// height of the object in pixels
	Height(size float64) int
}

// implements the ImageElement interface
//
// represents a piece of text drawn on an image at a given color, font and font size
//
// should be created with the NewTextElement method
//
// face is recomputed at Draw time
type TextElement struct {
	Text  string
	Color color.Color
	Font  *opentype.Font
}

// create a new Text Element
//
// this will create the underlying face needed to do work
func NewTextElement(text string, color color.Color, size float64, font *opentype.Font) (*TextElement, error) {
	return &TextElement{
		Text:  text,
		Color: color,
		Font:  font,
	}, nil
}

func (e *TextElement) Draw(img *image.RGBA, size float64, x int, y int) {
	face, _ := opentype.NewFace(e.Font, &opentype.FaceOptions{Size: size, DPI: 72})
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(e.Color),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(e.Text)
}

func (e *TextElement) Width(size float64) int {
	face, _ := opentype.NewFace(e.Font, &opentype.FaceOptions{Size: size, DPI: 72})
	return font.MeasureString(face, e.Text).Ceil()
}

func (e *TextElement) Height(size float64) int {
	face, _ := opentype.NewFace(e.Font, &opentype.FaceOptions{Size: size, DPI: 72})
	return face.Metrics().Height.Ceil()
}

// implements the ImageElement interface
//
// represents an image file
//
// should be created with the NewImageElement method
type EmbeddedImageElement struct {
	Path string
	Img  image.Image
}

func (i *EmbeddedImageElement) Draw(img *image.RGBA, size float64, x, y int) {
	offset := image.Pt(x, y)
	bounds := i.Img.Bounds().Add(offset)
	draw.Draw(img, bounds, i.Img, image.Point{}, draw.Over)
}
func (e *EmbeddedImageElement) Height(size float64) int {
	return e.Img.Bounds().Dy()
}
func (e *EmbeddedImageElement) Width(size float64) int {
	return e.Img.Bounds().Dx()
}

// implements the ImageElement interface
//
// represents a text file drawn at a given font and font size
//
// should be created with the NewEmbededTextFileElement method
type EmbededTextFileElement struct {
	Path     string
	Size     float64
	Font     *opentype.Font
	contents []string
}

func NewEmbededTextFileElement(path string, size float64, fnt *opentype.Font) (*EmbededTextFileElement, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return &EmbededTextFileElement{
		Path:     path,
		Size:     size,
		Font:     fnt,
		contents: lines,
	}, scanner.Err()
}
func (e *EmbededTextFileElement) Draw(img *image.RGBA, size float64, x int, y int) {
	face, _ := opentype.NewFace(e.Font, &opentype.FaceOptions{Size: size, DPI: 72})
	startX := x
	startY := y
	for _, line := range e.contents {
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(color.Black),
			Face: face,
			Dot:  fixed.P(startX, startY),
		}
		d.DrawString(line)
		startY += e.Height(size)
	}
}
func (e *EmbededTextFileElement) Height(size float64) int {
	face, _ := opentype.NewFace(e.Font, &opentype.FaceOptions{Size: size, DPI: 72})
	return face.Metrics().Height.Ceil() * len(e.contents)
}
func (e *EmbededTextFileElement) Width(size float64) int {
	face, _ := opentype.NewFace(e.Font, &opentype.FaceOptions{Size: size, DPI: 72})
	maxWidth := 0
	for _, line := range e.contents {
		w := font.MeasureString(face, line).Ceil()
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}
