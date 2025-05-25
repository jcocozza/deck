package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{
	Slides []Slide

	curr *ebiten.Image

	current int
	width int
	height int

	lastWidth int
	lastHeight int

	// keep track if certain things are pressed
	// the state updates faster then the user can handle and will skip slides
	fPressed bool
	rightPressed bool
	leftPressed bool

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
		g.redraw= true

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
		img := Draw(g.width, g.height, g.Slides[g.current])
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

func Render(slides []Slide) {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("TODO: <Think of some cool tool name>")
	g := &Game{
		Slides: slides,
		curr: ebiten.NewImage(640,480),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
