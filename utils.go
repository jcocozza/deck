package main

import (
	"image"
	"image/color"
	"image/draw"
	"strings"

	"golang.org/x/image/font"
)

const tabSize = 4

func WrapText(text string, face font.Face, maxWidth int) []string {
	// TODO: this is just an approximation
	tabReplacement := strings.Repeat(" ", tabSize)
	text = strings.ReplaceAll(text, "\t", tabReplacement)
	words := strings.Split(text, " ")

	spaceWidth := font.MeasureString(face, " ").Ceil()
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var builder strings.Builder
	lineLen := 0

	for i, word := range words {
		wordWidth := font.MeasureString(face, word).Ceil()
		if wordWidth+lineLen > maxWidth {
			lines = append(lines, builder.String())
			builder.Reset()
			lineLen = 0
		}
		builder.WriteString(word)
		lineLen += wordWidth

		if i != len(words)-1 {
			builder.WriteString(" ")
			lineLen += wordWidth + spaceWidth
		}
	}
	if builder.Len() != 0 {
		lines = append(lines, builder.String())
	}
	return lines
}

// GenerateSampleImage returns a 200x200 image with a red background and a blue circle
func GenerateSampleImage() image.Image {
	width, height := 200, 200
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill red background
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{red}, image.Point{}, draw.Src)

	// Draw a blue circle in the center
	blue := color.RGBA{0, 0, 255, 255}
	cx, cy, r := width/2, height/2, 60
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				img.Set(cx+x, cy+y, blue)
			}
		}
	}
	return img
}
