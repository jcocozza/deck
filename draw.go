package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func fill(img *image.RGBA, bgc color.Color) {
	draw.Draw(img, img.Bounds(), &image.Uniform{bgc}, image.Point{}, draw.Src)
}

type Element interface {
	// draw the element on the passed image
	Draw(img *image.RGBA, x, y int)
	// Height of the element in pixels
	Height() int
}

type TextElement struct {
	Content string
	lines   []string
	color color.Color

	face font.Face
}

func NewTextElement(content string, face font.Face, color color.Color, maxWidth int) *TextElement {
	lines := WrapText(content, face, maxWidth)
	return &TextElement{
		Content: content,
		lines:   lines,
		face:    face,
		color: color,
	}
}

func (e *TextElement) Draw(img *image.RGBA, x, y int) {
	adjY := y
	for _, line := range e.lines {
		fmt.Println("drawing line: ", line)
		fmt.Println("drawing line: ", e.face.Metrics())
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(e.color),
			Face: e.face,
			Dot:  fixed.P(x, adjY),
		}
		d.DrawString(line)
		adjY += e.face.Metrics().Height.Ceil()
	}
}

func (e *TextElement) Height() int {
	m := e.face.Metrics()
	return m.Height.Ceil() * len(e.lines)
}


type ImageElement struct {
	Img image.Image
	//ScaledWidth int
	//ScaledHeight int
}

func (e *ImageElement) Draw(img *image.RGBA,  x, y int) {
	offset := image.Pt(x, y)
	bounds := e.Img.Bounds().Add(offset)
	draw.Draw(img, bounds, e.Img, image.Point{}, draw.Over)
}

func (e *ImageElement) Height() int {
	return e.Img.Bounds().Dy()
}


func Draw(width, height int, slide Slide) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	fill(img, slide.Theme.Background)

	xPadding := 60
	yPadding := 60
	y := 0 + yPadding

	// set up content
	titleff, err := FontFace(slide.Font, 18)
	if err != nil {
		return nil, err
	}
	defer titleff.Close()
	textff, err := FontFace(slide.Font, 14)
	if err != nil {
		return nil, err
	}
	defer textff.Close()
	elms := []Element{
		NewTextElement(slide.Title, titleff, slide.Theme.Heading, width-(2*xPadding)),
	}
	for _, line := range slide.Content {
		elms = append(elms, NewTextElement(line, textff, slide.Theme.Text, width-(2*xPadding)))
	}
	if slide.Image != nil {
		elms = append(elms, &ImageElement{Img: slide.Image})
	}
	// actually draw stuff
	for _, elm := range elms {
		var xDraw int
		switch e := elm.(type) {
		case *TextElement:
			xDraw = xPadding
		case *ImageElement:
			xDraw = (width - e.Img.Bounds().Dx()) / 2
		}
		elm.Draw(img, xDraw, y)
		y += elm.Height()
	}
	return img, nil
}
