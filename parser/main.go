package parser

import (
	"bufio"
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

//func ReadAndParse(path string) ([]SlideTokens, error) {
//	lexer := NewLexer()
//	parser := NewParser()
//	beautifier := NewBeautifier(DefaultTheme)
//	lines, err := readlines(path)
//	if err != nil {
//		return nil, err
//	}
//	ll := lexer.Lex(lines)
//	cnts := parser.Parse(ll)
//	//for _, c := range cnts {
//	//	fmt.Println(c.String())
//	//}
//	tokens := beautifier.Beautify(cnts)
//	fmt.Println(tokens)
//	return tokens, nil
//}
