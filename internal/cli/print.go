package cli

import (
	"fmt"
	"os"
)

func ExitError(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func Warn(msg string) {
	fmt.Fprintf(os.Stderr, "[WARN] %s\n", msg)
}

func Info(msg string) {
	fmt.Fprintln(os.Stdout, msg)
}
