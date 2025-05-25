package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("not enough args")
	}
	path := os.Args[1]
	slides, err := ParseFile(path)
	if err != nil {panic(err)}
	Render(slides)
}
