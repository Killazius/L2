package utils

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func GetDir(urlPath string) string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd + string(os.PathSeparator) + urlPath

}

func URLToFileName(baseURL *url.URL, outputDir string) string {
	var filename string
	path := baseURL.Path
	if path == "" || strings.HasSuffix(path, "/") {
		filename = "index.html"
	} else {
		segments := strings.Split(path, "/")
		filename = segments[len(segments)-1]

		if filepath.Ext(filename) == "" {
			filename += ".html"
		}
	}

	relPath := strings.TrimPrefix(path, "/")
	dirPart := filepath.Dir(relPath)
	if dirPart == "." {
		dirPart = ""
	}
	fullPath := filepath.Join(outputDir, dirPart)
	return filepath.Join(fullPath, filename)
}
