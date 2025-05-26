package parser

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"
)


type ContentType int

const (
	Paragraph ContentType = iota
	Header
	Text
	List
	File
)

var cts = [...]string{
	Paragraph: "paragraph",
	Header: "header",
	Text: "text",
	List: "list",
	File: "file",
}

type Content struct {
	T      ContentType
	Text   []string
	Img    image.Image // will be null mostly
	Children []*Content
	Level  int
}

func (c *Content) String() string {
	tabs := strings.Repeat("\t", c.Level)
	txt := strings.Join(c.Text, ", ")
	s := fmt.Sprintf("%s[%s(%d)] %s\n", tabs, cts[c.T], c.Level, "["+txt+"]")
	for _, child := range c.Children {
		s += child.String()
	}
	return s
}

type Parser interface {
	Parse(lines []lexline) []*Content
}

func NewParser() Parser {
	return &ParserImpl{}
}

type ParserImpl struct{}

func (p *ParserImpl) Parse(lines []lexline) []*Content {
	var root []*Content
	var currRoot *Content
	var curr *Content

	var lastType linetype

	for _, line := range lines {
		if currRoot == nil {
			currRoot = &Content{
				T: Paragraph,
			}
			curr = currRoot
		}
		switch line.t {
		case title:
			h := &Content{
				T: Header,
				Text: []string{line.text},
				Level: 1,
			}
			currRoot.Children = append(currRoot.Children, h)
			curr = h
			lastType = title
		case subtitle:
			h := &Content{
				T: Header,
				Text: []string{line.text},
				Level: 2,
			}
			curr.Children = append(curr.Children, h)
			curr = h
			lastType = subtitle
		case comment:
			continue
		case file: // a file is always a leaf node
			path := strings.TrimPrefix(line.text, prefixes[file])
			var img image.Image
			var lines []string
			var err error

			switch filepath.Ext(path) {
			case ".png", ".jpg", ".gif":
				img, err = ReadImage(path)
				if err != nil {
					lines = []string{fmt.Sprintf("ERR: unable to read %s. Err is: %s", path, err.Error())}
				}

			case ".txt", ".py":
				lines, err = ReadTxt(path)
				if err != nil {
					lines = []string{fmt.Sprintf("ERR: unable to read %s. Err is: %s", path, err.Error())}
				}
			default:
				lines = []string{fmt.Sprintf("ERR: unsupported file: %s", path)}
			}
			h := &Content{
				T: File,
				Text: lines,
				Img: img,
				Level: curr.Level+1,
			}
			curr.Children = append(curr.Children, h)
			lastType = file
		case emptyLine:
			curr = nil
			root = append(root, currRoot)
			currRoot = nil
			lastType = emptyLine
		case text: // text is always a leaf node
			h := &Content{
				T: Text,
				Text: []string{line.text},
				Level: curr.Level+1,
			}
			curr.Children = append(curr.Children, h)
			lastType = text
		case list:
			if lastType == list {
				curr.Text = append(curr.Text, line.text)
				lastType = list // for clarity
				continue
			}
			// otherwise start a new list
			h := &Content{
				T: List,
				Text: []string{line.text},
				Level: curr.Level+1,
			}
			curr.Children = append(curr.Children, h)
			curr = h
			lastType = list
		}
	}
	// handle the last one
	if currRoot != nil {
		root = append(root, currRoot)
	}
	return root
}

