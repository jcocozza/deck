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

// if len(paths) == 0 or the only arg is "-", read from stdin
func ReadFromStdinOrFiles(paths []string) ([]string, error) {
	var lines []string

	readStdin := len(paths) == 1 && paths[0] == "-"
	if len(paths) > 0 && !readStdin {
		for _, path := range paths {
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
