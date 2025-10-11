package config

import (
	"io"
	"os"

	"github.com/spf13/pflag"
)

type Config struct {
	filePath string
	Pattern  string

	After      int  // -A N
	Before     int  // -B N
	Context    int  // -C N
	CountOnly  bool // -c
	IgnoreCase bool // -i
	Invert     bool // -v
	Fixed      bool // -F
	LineNum    bool // -n
}

func New() *Config {
	cfg := &Config{}

	pflag.IntVarP(&cfg.After, "after", "A", 0, "print N lines after match")
	pflag.IntVarP(&cfg.Before, "before", "B", 0, "print N lines before match")
	pflag.IntVarP(&cfg.Context, "context", "C", 0, "print N lines of context around match")
	pflag.BoolVarP(&cfg.CountOnly, "count", "c", false, "print only count of matching lines")
	pflag.BoolVarP(&cfg.IgnoreCase, "ignore-case", "i", false, "ignore case")
	pflag.BoolVarP(&cfg.Invert, "invert-match", "v", false, "invert match")
	pflag.BoolVarP(&cfg.Fixed, "fixed-strings", "F", false, "use fixed string matching")
	pflag.BoolVarP(&cfg.LineNum, "line-number", "n", false, "print line number")

	pflag.Parse()

	if pflag.NArg() == 0 {
		panic("grep [OPTION]... PATTERNS [FILE]...")
	}
	cfg.Pattern = pflag.Arg(0)
	if pflag.NArg() > 1 {
		cfg.filePath = pflag.Arg(1)
	}
	if cfg.After < 0 || cfg.Before < 0 || cfg.Context < 0 {
		panic("Значения для -A, -B, -C не могут быть отрицательными")
	}

	if cfg.Context > 0 {
		cfg.After = cfg.Context
		cfg.Before = cfg.Context
	}
	return cfg
}

func (cfg *Config) GetSource() (io.ReadCloser, error) {
	if cfg.filePath != "" {
		file, err := os.Open(cfg.filePath)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	return os.Stdin, nil
}
