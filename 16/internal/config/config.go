package config

import (
	"flag"
	"time"
)

type Config struct {
	Depth      uint
	Timeout    int64
	Retries    uint
	NumWorkers int64
}

func (c Config) TimeoutSeconds() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

func New() *Config {
	duration := flag.Int64("t", 10, "timeout seconds per page (default 10)")
	depth := flag.Uint("d", 0, "depth (default 0 for only root page)")
	retries := flag.Uint("r", 1, "num of tries (default 1)")
	numWorkers := flag.Int64("w", 10, "num of workers(default 10)")
	flag.Parse()

	if *duration <= 0 {
		*duration = 10
	}
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
