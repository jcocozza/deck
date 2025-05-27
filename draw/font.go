package draw

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func ReadFont(path string) (*opentype.Font, error) {
	fbytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return opentype.Parse(fbytes)
}

func FontFace(font *opentype.Font, size int) (font.Face, error) {
	return opentype.NewFace(
		font,
		&opentype.FaceOptions{Size: float64(size), DPI: 72},
	)
}

func ReadTestFont() *opentype.Font {
	f, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		panic(err)
	}
	return f
}
