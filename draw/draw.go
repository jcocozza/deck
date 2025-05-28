package draw

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/jcocozza/deck/slide"
)

func fillImg(img *image.RGBA, c color.Color) {
	draw.Draw(img, img.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
}

func NewImage(width int, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return img
}

func Draw(width int, height int, s slide.Slide, pretty bool) image.Image {
	fnt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		panic(err)
	}
	img := NewImage(width, height)
	fillImg(img, color.White)
	drawer := DrawerFactory(pretty, fnt)
	fmt.Println("drawing slide;", s)
	err = drawer.Draw(img, 30, 30, s)
	if err != nil {
		panic(err)
	}
	return img
}
