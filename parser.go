package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strings"
)

const (
	PREFIX_Title   = "# "
	PREFIX_Comment = "// "
	PREFIX_File    = "@"
)

func NewEmptySlide() Slide {
	return Slide{
		//Theme: DefaultTheme,
	}
}

func ReadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func ReadTxt(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		ln := scanner.Text()
		lines = append(lines, ln)
	}
	return lines, scanner.Err()
}

// Parsing Rules:
// - each paragraph is a new slide
// - titles are denoted via "# "
// - lines starting with // are ignored
// - use @ to import and display a file
//   - images are displayed as images
//   - text will be imported into the slide
func ParseFile(fname string) ([]Slide, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var slide Slide
	var slides []Slide
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ln := scanner.Text()
		ln = strings.TrimSpace(ln)

		switch {
		case strings.HasPrefix(ln, PREFIX_Comment):
			fmt.Println("skipping comment")
			continue
		case ln == "":
			slides = append(slides, slide)
			slide = NewEmptySlide()
			fmt.Println("adding new slide")
		case strings.HasPrefix(ln, PREFIX_Title):
			slide.Title = strings.TrimPrefix(ln, PREFIX_Title)
			fmt.Printf("added title: %s\n", slide.Title)
		case strings.HasPrefix(ln, PREFIX_File):
			path := strings.TrimPrefix(ln, PREFIX_File)
			switch filepath.Ext(path) {
			case ".png", ".jpeg", ".gif": // TODO: actual gif support one day?
				img, err := ReadImage(path)
				if err != nil {
					slide.Content = append(slide.Content, fmt.Sprintf("ERR: unable to read %s. Err is: %s", path, err.Error()))
					continue
				}
				slide.Image = img
			case ".txt", ".py":
				lines, err := ReadTxt(path)
				if err != nil {
					slide.Content = append(slide.Content, fmt.Sprintf("ERR: unable to read %s. Err is: %s", path, err.Error()))
					continue
				}
				slide.Content = append(slide.Content, "```")
				slide.Content = append(slide.Content, lines...)
				slide.Content = append(slide.Content, "```")
			default:
				slide.Content = append(slide.Content, fmt.Sprintf("ERR: unsupported file: %s", path))
				continue

			}
		default:
			fmt.Printf("[slide %s] adding to slide body %s\n", slide.Title, ln)
			slide.Content = append(slide.Content, ln)
		}
	}
	slides = append(slides, slide)
	return slides, nil
}
