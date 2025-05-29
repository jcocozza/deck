package slide

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type ImgPostion int
const (
	Left ImgPostion = iota
	Right
	Bottom
	Top
)

type Image struct {
	Path string
	I    image.Image
	Position ImgPostion
}

func NewImage(path string) (*Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return &Image{
		Path: path,
		I:    img,
		Position: Right,
	}, nil
}

type Slide struct {
	Lines []string
	Image *Image
}

func TestSlides() []Slide {
	testImg, err := NewImage("test.png")
	if err != nil {
		panic(err)
	}
	s := []Slide{
		{Lines: []string{"some", "lines", "of text"}, Image: nil},
		{Lines: []string{"some lines of text"}, Image: testImg},
		{Lines: nil, Image: testImg},
	}
	return s
}
