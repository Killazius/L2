package main

import (
	"16/internal/config"
	"16/internal/downloader"
	"16/internal/parser"
	"16/internal/utils"
	"flag"
	"log/slog"
	"net/url"
	"os"
)

func main() {
	cfg := config.New()
	args := flag.Args()
	if len(args) != 1 {
		slog.Error("Please provide a single URL as an argument")
		return
	}
	rawURL := args[0]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		slog.Error("Invalid URL provided", "error", err)
		return
	}

	dir := utils.GetDir(parsedURL.Host)
	err = os.Mkdir(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		slog.Error("Failed to create directory", "error", err)
		return
	}

	loader := downloader.New(
		cfg,
		parser.New(dir),
	)
	if err := loader.Start(parsedURL); err != nil {
		slog.Error("Download failed", "error", err)
		return
	}

}
