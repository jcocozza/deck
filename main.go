package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("not enough args")
	}
	theme, err := ReadTheme("ppt/unpacked/ppt/theme/theme1.xml")
	if err != nil {
		panic(err)
	}

	fnt, err := ReadFont("/usr/share/fonts/opentype/urw-base35/P052-Roman.otf")
	if err != nil {
		panic(err)
	}
	//fnt := readTestFont()

	path := os.Args[1]
	slides, err := ParseFile(path, theme, fnt)
	if err != nil {
		panic(err)
	}
	Render(slides)
}
