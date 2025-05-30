package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"github.com/jcocozza/deck/draw"
	"github.com/jcocozza/deck/parser"
	"github.com/jcocozza/deck/render"
)

func readInput(r io.Reader) ([]string, error) {
	lines := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ln := scanner.Text()
		lines = append(lines, ln)
	}
	return lines, scanner.Err()
}

func main() {
	var lines []string
	if len(os.Args) > 1 { // read from a list of file paths
		for _, path := range os.Args[1:] {
			f, err := os.Open(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}
			defer f.Close()
			flines, err := readInput(f)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}
			lines = append(lines, flines...)
			// need to append a new line to ensure each file is separate
			lines = append(lines, "")
		}
	} else { // read from stdin
		flines, err := readInput(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		lines = flines
	}

	theme, err := draw.ReadTheme("ppt/unpacked/ppt/theme/theme1.xml")
	//theme := draw.DefaultTheme
	if err != nil {
		panic(err)
	}
	fnt, err := draw.ReadFont("/usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf")
	//fnt := draw.ReadTestFont()
	if err != nil {
		panic(err)
	}

	cnts := parser.Parse(lines)
	if err != nil {
		panic(err)
	}
	render.Render(cnts, theme, fnt)
}
