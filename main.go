package main

import (
	"os"
	"presentation/draw"
	"presentation/parser"
	"presentation/render"
)

//func main() {
//	if len(os.Args) < 2 {
//		panic("not enough args")
//	}
//
//	path := os.Args[1]
//	slides, err := ParseFile(path, theme, fnt)
//	if err != nil {
//		panic(err)
//	}
//	Render(slides)
//}

func main() {
	if len(os.Args) < 2 {
		panic("not enough args")
	}

	theme, err := draw.ReadTheme("ppt/unpacked/ppt/theme/theme1.xml")
	if err != nil {
		panic(err)
	}
	//theme := draw.DefaultTheme

	fnt, err := ReadFont("/usr/share/fonts/opentype/urw-base35/P052-Roman.otf")
	if err != nil {
		panic(err)
	}
	//fnt := readTestFont()

	path := os.Args[1]
	lines, err := parser.ReadAndParse(path)
	if err != nil {
		panic(err)
	}
	//for _, line := range lines {
	//	fmt.Println(line.String())
	//}
	render.Render(lines, theme, fnt)
	//render.Render(lines, parser.DefaultTheme, readTestFont())
}
