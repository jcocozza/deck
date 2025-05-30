package main

import (
	"fmt"
	"os"

	"github.com/jcocozza/deck/internal/cli"
)

func main() {
	err := cli.Cli()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
