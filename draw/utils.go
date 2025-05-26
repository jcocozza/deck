package draw

import (
	"image"
	"os"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/font"
)

const tabSize = 4

// TODO: currently this does nothing
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

func LoadImageFromFile(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}
