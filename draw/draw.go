package draw

import (
	"image"
	"image/color"
	"image/draw"
	"presentation/parser"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func fill(img *image.RGBA, bgc color.Color) {
	draw.Draw(img, img.Bounds(), &image.Uniform{bgc}, image.Point{}, draw.Src)
}

type ImageItem interface {
	// apply item to passed image
	Draw(img *image.RGBA, x, y int)
	// height in pixels of the item
	Height() int
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
	d.DrawString(i.Text) // TODO: beautify text
}

func (i *TextImageItem) Height() int {
	return i.Face.Metrics().Height.Ceil()
}

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
	yPad := 60
	y := 0 + yPad

	iItems := generateItems(s, d.b, d.theme)
	for _, item := range iItems {
		item.Draw(img, xPad, y)
		y += item.Height()
	}

	return img
}
