package parser

import (
	"16/internal/utils"
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type Parser struct {
	dir string
}

func New(dirPath string) *Parser {
	return &Parser{
		dir: dirPath,
	}
}

type ParseResult struct {
	Data     []byte
	SubLinks []string
	Filename string
}

func (p *Parser) Parse(resp *http.Response, u *url.URL) (*ParseResult, error) {
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			slog.Warn("Failed to close response body", "error", closeErr)
		}
	}(resp.Body)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sublinks []string
	var processedData = data

	if strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		links, err := p.extractLinks(bytes.NewReader(data), u)
		if err != nil {
			slog.Warn("Failed to extract links", "url", u.String(), "error", err)
		}
		sublinks = links

		processedData, err = p.replaceLinksInHTML(data, u)
		if err != nil {
			slog.Warn("Failed to replace links in HTML", "url", u.String(), "error", err)
			processedData = data
		}
	}

	return &ParseResult{
		Data:     processedData,
		SubLinks: sublinks,
		Filename: utils.URLToFileName(u, p.dir),
	}, nil
}

func (p *Parser) extractLinks(body io.Reader, baseURL *url.URL) ([]string, error) {
	var links []string
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return links, nil
			}
			return nil, tokenizer.Err()
		}
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			if token.Data == "a" || token.Data == "link" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						linkURL, err := baseURL.Parse(attr.Val)
						if err != nil {
							slog.Debug("Failed to parse link", "link", attr.Val, "error", err)
							continue
						}
						if linkURL.Host == baseURL.Host {
							links = append(links, linkURL.String())
						}
					}
				}
			}
		}
	}
}

func (p *Parser) replaceLinksInHTML(data []byte, baseURL *url.URL) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	p.replaceLinksInNode(doc, baseURL)

	var buf bytes.Buffer
	err = html.Render(&buf, doc)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *Parser) replaceLinksInNode(n *html.Node, baseURL *url.URL) {
	if n.Type == html.ElementNode {
		for i, attr := range n.Attr {
			if (attr.Key == "href" && (n.Data == "a" || n.Data == "link")) ||
				(attr.Key == "src" && (n.Data == "img" || n.Data == "script")) {

				linkURL, err := baseURL.Parse(attr.Val)
				if err != nil {
					continue
				}

				if linkURL.Host == baseURL.Host {
					localPath := p.getRelativePath(baseURL, linkURL)
					n.Attr[i].Val = localPath
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.replaceLinksInNode(c, baseURL)
	}
}

func (p *Parser) getRelativePath(currentURL, targetURL *url.URL) string {
	currentFile := utils.URLToFileName(currentURL, p.dir)
	targetFile := utils.URLToFileName(targetURL, p.dir)

	currentDir := filepath.Dir(currentFile)
	relPath, err := filepath.Rel(currentDir, targetFile)
	if err != nil {
		return filepath.Base(targetFile)
	}

	return filepath.ToSlash(relPath)
}
