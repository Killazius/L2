package main

import (
	"13/internal/config"
	"13/internal/cut"
	"fmt"
	"io"
	"os"
)

func main() {
	cfg := config.New()
	src := os.Stdin
	c := cut.New(cfg, src)
	var dst io.Writer = os.Stdout
	if err := c.Run(dst); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
