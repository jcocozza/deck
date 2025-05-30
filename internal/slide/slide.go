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
	// NOTE: Any text with the center tag will not be shown
	// this is enforced at the parser level
	Center
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

type SlideLineType int

const (
	Header SlideLineType = iota
	Subheader
	Subsubheader
	Text
)

type SlideLine struct {
	Text string
	T SlideLineType
}

type Slide struct {
	Lines []SlideLine
	Image *Image // TODO: generalize to the an aribrary embedded filetype
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
	testImgCenter, err := NewImage("test.png", Center)
	if err != nil {
		panic(err)
	}


	someLines := []SlideLine{
		SlideLine{Text: "some", T: Text},
		SlideLine{Text: "lines", T: Text},
		SlideLine{Text: "of text", T: Text},
	}

	left := []SlideLine{
		SlideLine{Text: "left", T: Text},
	}
	right := []SlideLine{
		SlideLine{Text: "right", T: Text},
	}
	bottom := []SlideLine{
		SlideLine{Text: "bottom", T: Text},
	}
	top := []SlideLine{
		SlideLine{Text: "top", T: Text},
	}
	center := []SlideLine{
		SlideLine{Text: "center", T: Text},
	}
	list := []SlideLine{
		SlideLine{Text: "list", T: Text},
		SlideLine{Text: "1. foo", T: Text},
		SlideLine{Text: "2. bar", T: Text},
		SlideLine{Text: "3. baz", T: Text},
	}

	s := []Slide{
		{Lines: someLines, Image: nil},
		{Lines: left, Image: testImgLeft},
		{Lines: right, Image: testImgRight},
		{Lines: bottom, Image: testImgBottom},
		{Lines: top, Image: testImgTop},
		{Lines: center, Image: testImgCenter},
		{Lines: nil, Image: testImgCenter},
		{Lines: list, Image: nil},
	}
	return s
}
