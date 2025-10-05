package main

import (
	"10/internal/config"
	"10/internal/parser"
	"10/internal/sorter"
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	cfg := config.New()
	p := parser.New(cfg, cfg.GetSource())
	s := sorter.New(p, cfg)
	sorted, err := s.Sort()
	if err != nil {
		panic(err)
	}
	var dst io.Writer = os.Stdout
	if err = writeLines(sorted, dst); err != nil {
		panic(err)
	}
}

func writeLines(lines []string, dst io.Writer) error {
	w := bufio.NewWriter(dst)
	for _, line := range lines {
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return w.Flush()
}
