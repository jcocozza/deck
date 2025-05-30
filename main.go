package main

import (
	"fmt"
	"os"

	"github.com/jcocozza/deck/internal/format"
	"github.com/jcocozza/deck/internal/render"
	"github.com/jcocozza/deck/internal/utils"
)

func main() {
	lines, err := utils.ReadFromStdinOrFiles()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	lexer := format.NewLexer()
	parser := format.NewParser()

	lexLines := lexer.Lex(lines)
	slides, err := parser.Parse(lexLines)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	// slides := slide.TextSlides()
	render.Render(slides)
}
