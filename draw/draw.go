package draw

import (
	"image"
	"image/color"
	"image/draw"
	"github.com/jcocozza/deck/parser"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func fill(img *image.RGBA, bgc color.Color) {
	draw.Draw(img, img.Bounds(), &image.Uniform{bgc}, image.Point{}, draw.Src)
}

// an ImageItem is a chunk of data that is to be drawn on the passed image
// it is _not_ necessarily 1 one line
//  e.g an image or a colored piece of text
type ImageItem interface {
	// apply item to passed image
	Draw(img *image.RGBA, x, y int)
	// height in pixels of the item
	Height() int
	Width() int
}

// a TextImageItem is a single single piece of text to be drawn
// this is not necessarily 1 line, it can be as small as a single char
type TextImageItem struct {
	Text  string
	Color color.Color
	Face  font.Face
}

func (i *TextImageItem) Draw(img *image.RGBA, x, y int) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(i.Color),
		Face: i.Face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(i.Text)
}

func (i *TextImageItem) Height() int {
	return i.Face.Metrics().Height.Ceil()
}

func (i *TextImageItem) Width() int {
	return font.MeasureString(i.Face, i.Text).Ceil()
}

type NewLineItem struct{ Face font.Face }

func (i *NewLineItem) Draw(img *image.RGBA, x, y int) {}
func (i *NewLineItem) Height() int {
	return i.Face.Metrics().Height.Ceil()
}
func (i *NewLineItem) Width() int { return 0 }

type EmbededImageItem struct {
	Img image.Image
}

func (i *EmbededImageItem) Draw(img *image.RGBA, x, y int) {
	offset := image.Pt(x, y)
	bounds := i.Img.Bounds().Add(offset)
	draw.Draw(img, bounds, i.Img, image.Point{}, draw.Over)
}
func (i *EmbededImageItem) Height() int {
	return i.Img.Bounds().Dy()
}

func (i *EmbededImageItem) Width() int {
	return i.Img.Bounds().Dx()
}

type Drawer interface {
	DrawSlide(width int, height int, slide *parser.Content) image.Image
}

func NewDrawer(theme Theme, font *opentype.Font) Drawer {
	return &DrawerImpl{
		theme: theme,
		font:  font,
		b:     NewBeutifier(font),
	}
}

type DrawerImpl struct {
	b     Beautifier
	theme Theme
	font  *opentype.Font
}

func generateItems(c *parser.Content, b Beautifier, t Theme) []ImageItem {
	items := []ImageItem{}
	if c.Img != nil {
		items = append(items, &EmbededImageItem{Img: c.Img})
		// new to add a newline after image to leave room for things
		items = append(items, &NewLineItem{Face: b.Face(12)}) // TODO: font size should not be hard coded
	}
	if len(c.Text) > 0 {
		i := b.Beautify(c.Text, c.T, c.Level, t)
		items = append(items, i...)
	}

	for _, child := range c.Children {
		items = append(items, generateItems(child, b, t)...)
	}
	return items
}

func (d *DrawerImpl) DrawSlide(width int, height int, s *parser.Content) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	fill(img, d.theme.Background)

	xPad := 60
	xShift := 0
	yPad := 60
	yshift := 0

	iItems := generateItems(s, d.b, d.theme)
	for _, item := range iItems {
		var xDraw int
		switch e := item.(type) {
		case *NewLineItem:
			yshift += item.Height()
			xShift = 0
			continue // we don't draw newline items
		case *TextImageItem:
			xDraw = xPad
			item.Draw(img, xDraw+xShift, yshift+yPad)
			xShift += item.Width()
		case *EmbededImageItem:
			xDraw = (width - e.Img.Bounds().Dx()) / 2 //center the image
			item.Draw(img, xDraw+xShift, yshift+yPad)
			yshift += item.Height()
			xShift += item.Width()
		}
	}

	return img
}
