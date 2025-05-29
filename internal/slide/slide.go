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

func NewImage(path string, postion ImgPostion) (*Image, error) {
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
		Position: postion,
	}, nil
}

type Slide struct {
	Lines []string
	Image *Image
}

func TestSlides() []Slide {
	testImgLeft, err := NewImage("test.png", Left)
	if err != nil {
		panic(err)
	}
	testImgRight, err := NewImage("test.png", Right)
	if err != nil {
		panic(err)
	}
	testImgTop, err := NewImage("test.png", Top)
	if err != nil {
		panic(err)
	}
	testImgBottom, err := NewImage("test.png", Bottom)
	if err != nil {
		panic(err)
	}
	s := []Slide{
		{Lines: []string{"some", "lines", "of text"}, Image: nil},
		{Lines: []string{"some lines of text"}, Image: testImgLeft},
		{Lines: []string{"left"}, Image: testImgLeft},
		{Lines: []string{"right"}, Image: testImgRight},
		{Lines: []string{"bottom"}, Image: testImgBottom},
		{Lines: []string{"top"}, Image: testImgTop},
	}
	return s
}
