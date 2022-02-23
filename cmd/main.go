package main

import (
	"os"

	"github.com/sbrki/monkey/pkg/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
