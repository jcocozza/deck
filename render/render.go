package render

import (
	"log"

	"github.com/jcocozza/deck/draw"
	"github.com/jcocozza/deck/slide"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Slides []slide.Slide
	curr   *ebiten.Image

	current int
	width   int
	height  int

	lastWidth  int
	lastHeight int

	// keep track if certain things are pressed
	// the state updates faster then the user can handle and will skip slides
	fPressed     bool
	rightPressed bool
	leftPressed  bool

	// true if a redraw is needed
	redraw bool
}

func (g *Game) Update() error {
	fPressed := ebiten.IsKeyPressed(ebiten.KeyF)
	if fPressed && !g.fPressed {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	g.fPressed = fPressed

	// handle exits
	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyC) {
		return ebiten.Termination
	}
	rightPressed := ebiten.IsKeyPressed(ebiten.KeyRight)
	leftPressed := ebiten.IsKeyPressed(ebiten.KeyLeft)

	if rightPressed && !g.rightPressed && g.current < len(g.Slides)-1 {
		g.current++
		g.redraw = true

	}
	if leftPressed && !g.leftPressed && g.current > 0 {
		g.current--
		g.redraw = true
	}
	g.rightPressed = rightPressed
	g.leftPressed = leftPressed

	// detect resize
	if g.width != g.lastWidth || g.height != g.lastHeight {
		g.redraw = true
		g.lastWidth = g.width
		g.lastHeight = g.height
	}

	// only redraw the image if we have changed slides
	if g.redraw {
		img := draw.Draw(g.width, g.height, g.Slides[g.current], false)
		g.curr = ebiten.NewImageFromImage(img)
		g.redraw = false
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	if g.curr != nil {
		screen.DrawImage(g.curr, opts)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.width = outsideWidth
	g.height = outsideHeight
	return outsideWidth, outsideHeight
}

func TestSlides() []slide.Slide {
	img, _ := draw.ReadImage("simpsons.gif")
	s := slide.Slide{Img: slide.SlideImage{Image: img}}
	s2 := slide.Slide{
		Lines: []slide.Line{
			{T: slide.Header, Content: "a;sldkfja;sldkfj"},
			{T: slide.SubHeader, Content: "Baz"},
			{T: slide.Text, Content: "a;sldkfja;lsdkjf"},
		},
	}
	s3 := slide.Slide{
	Lines: []slide.Line{
		{T: slide.Header, Content: "Foo Bar"},
		{T: slide.SubHeader, Content: "Baz"},
		{T: slide.Text, Content: "Much longer bit of text about nothing"},
	},
}
	return []slide.Slide{
		s3, s2, s,
	}
}

func Render() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("deck")
	g := &Game{
		Slides: TestSlides(),
		curr:   ebiten.NewImage(640, 480),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
