package config

import (
	"io"
	"os"

	"github.com/spf13/pflag"
)

type Config struct {
	K int
	N bool
	R bool
	U bool

	filePath string
	Delim    string
}

func New() *Config {
	cfg := &Config{}

	pflag.IntVarP(&cfg.K, "column", "k", 1, "sort by column N")
	pflag.BoolVarP(&cfg.N, "numeric", "n", false, "sort by numeric value")
	pflag.BoolVarP(&cfg.R, "reverse", "r", false, "sort in reverse order")
	pflag.BoolVarP(&cfg.U, "unique", "u", false, "output only unique lines")
	pflag.StringVarP(&cfg.Delim, "delimiter", "t", "\t", "field delimiter")

	pflag.Parse()

	if pflag.NArg() > 0 {
		cfg.filePath = pflag.Arg(0)
	}

	return cfg
}

func (cfg *Config) GetSource() io.Reader {
	if cfg.filePath != "" {
		file, err := os.Open(cfg.filePath)
		if err != nil {
			panic(err)
		}
		return file
	}
	return os.Stdin
}
