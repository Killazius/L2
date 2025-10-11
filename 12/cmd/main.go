package main

import (
	"12/internal/config"
	"12/internal/grep"
	"fmt"
	"io"
	"os"
)

func main() {
	cfg := config.New()
	src, err := cfg.GetSource()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer src.Close()
	g := grep.New(cfg, src)
	var dst io.Writer = os.Stdout
	if err := g.Run(dst); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
}
