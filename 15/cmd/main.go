package main

import (
	"15/internal/shell"
	"os"
)

func main() {
	sh := shell.New()
	if len(os.Args) == 2 {
		sh.RunBash(os.Args[1])
	} else {
		sh.RunInteractive()
	}
}
