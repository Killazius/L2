package config

import (
	"flag"
	"time"
)

type Config struct {
	Depth      uint
	Timeout    uint
	Retries    uint
	NumWorkers uint
}

func (c Config) TimeoutSeconds() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

func New() *Config {
	duration := flag.Uint("t", 10, "timeout seconds per page (default 10)")
	depth := flag.Uint("d", 0, "depth (default 0 for only root page)")
	retries := flag.Uint("r", 1, "num of tries (default 1)")
	numWorkers := flag.Uint("w", 10, "num of workers(default 10)")
	flag.Parse()
	if *numWorkers == 0 {
		*numWorkers = 10
	}
	if *retries == 0 {
		*retries = 1
	}

	return &Config{
		Depth:      *depth,
		Timeout:    *duration,
		Retries:    *retries,
		NumWorkers: *numWorkers,
	}
}
