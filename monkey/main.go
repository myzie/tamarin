package main

import (
	"os"

	"github.com/myzie/tamarin/monkey/repl"
)

func main() {
	repl.Run(os.Stdin, os.Stdout)
}
