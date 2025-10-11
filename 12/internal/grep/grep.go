package grep

import (
	"12/internal/config"
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Grep struct {
	cfg *config.Config
	src io.Reader
}

func New(cfg *config.Config, src io.Reader) *Grep {
	return &Grep{cfg: cfg, src: src}
}

func (g *Grep) Run(dst io.Writer) error {
	scanner := bufio.NewScanner(g.src)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	cfg := g.cfg
	matchCount := 0
	matched := make([]bool, len(lines))

	var re *regexp.Regexp
	pattern := cfg.Pattern
	if cfg.IgnoreCase {
		pattern = strings.ToLower(pattern)
	}

	if cfg.Fixed {
		for i, line := range lines {
			cmpLine := line
			if cfg.IgnoreCase {
				cmpLine = strings.ToLower(line)
			}
			if strings.Contains(cmpLine, pattern) != cfg.Invert {
				matched[i] = true
				matchCount++
			}
		}
	} else {
		var err error
		flags := ""
		if cfg.IgnoreCase {
			flags = "(?i)"
		}
		re, err = regexp.Compile(flags + pattern)
		if err != nil {
			return err
		}
		for i, line := range lines {
			if re.MatchString(line) != cfg.Invert {
				matched[i] = true
				matchCount++
			}
		}
	}

	if cfg.CountOnly {
		fmt.Fprintf(dst, "%d\n", matchCount)
		return nil
	}

	context := cfg.After
	before := cfg.Before
	if cfg.Context > 0 {
		context = cfg.Context
		before = cfg.Context
	}

	toPrint := make(map[int]struct{})
	for i := range matched {
		if matched[i] {
			for j := i - before; j < i; j++ {
				if j >= 0 {
					toPrint[j] = struct{}{}
				}
			}
			toPrint[i] = struct{}{}
			for j := i + 1; j <= i+context; j++ {
				if j < len(lines) {
					toPrint[j] = struct{}{}
				}
			}
		}
	}

	for i := 0; i < len(lines); i++ {
		if _, ok := toPrint[i]; ok {
			if cfg.LineNum {
				fmt.Fprintf(dst, "%d:", i+1)
			}
			fmt.Fprintln(dst, lines[i])
		}
	}
	return nil
}
