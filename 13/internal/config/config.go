package config

import (
	"flag"
)

type Config struct {
	Fields    string
	Delimiter string
	Separated bool
}

func New() *Config {
	fields := flag.String("f", "", "fields to output (e.g. 1,3-5)")
	delim := flag.String("d", "\t", "delimiter")
	sep := flag.Bool("s", false, "only lines with delimiter")
	flag.Parse()
	return &Config{
		Fields:    *fields,
		Delimiter: *delim,
		Separated: *sep,
	}
}
