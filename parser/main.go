package parser

import (
	"bufio"
	"fmt"
	"os"
)

func readlines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func Parse(lines []string) []*Content {
	lexer := NewLexer()
	parser := NewParser()
	ll := lexer.Lex(lines)
	for _, l := range ll {
		fmt.Println(l.String())
	}
	cnts := parser.Parse(ll)
	return cnts
}

func ReadAndParse(path string) ([]*Content, error) {
	lexer := NewLexer()
	parser := NewParser()
	lines, err := readlines(path)
	if err != nil {
		return nil, err
	}
	ll := lexer.Lex(lines)
	cnts := parser.Parse(ll)
	return cnts, nil
}
