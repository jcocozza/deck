package slide

import "image"

type LineType string

const (
	Header       LineType = "header"
	SubHeader    LineType = "sub-header"
	SubSubHeader LineType = "sub-sub-header"
	Text         LineType = "text"
	File         LineType = "file"
)

type Line struct {
	T LineType
	// Raw content read in from user
	Content string
}

type SlideImage struct {
	Path string
	Image image.Image
}

// a list is just a(n ordered) list of lines
type Slide struct {
	Lines []Line
	Img SlideImage
}

func (s *Slide) HasImg() bool {
	return s.Img.Path != ""
}
