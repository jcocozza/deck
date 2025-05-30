package utils

import (
	"bufio"
	"io"
	"os"
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

func readFileLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readInput(f)
}

func ReadFromStdinOrFiles() ([]string, error) {
	var lines []string
	if len(os.Args) > 1 {
		for _, path := range os.Args[1:] {
			lns, err := readFileLines(path)
			if err != nil {
				return nil, err
			}
			lines = append(lines, lns...)
		}
	} else {
		flines, err := readInput(os.Stdin)
		if err != nil {
			return nil, err
		}
		lines = flines
	}
	return lines, nil
}
