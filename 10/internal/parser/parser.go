package parser

import (
	"10/internal/config"
	"bufio"
	"errors"
	"io"
	"os"
)

type Parser struct {
	cfg    *config.Config
	source io.Reader
}

func New(cfg *config.Config, source io.Reader) *Parser {
	return &Parser{
		cfg:    cfg,
		source: source,
	}
}

func (p *Parser) Parse() ([]string, error) {
	if p.source == nil {
		return nil, errors.New("source is nil")
	}
	scanner := bufio.NewScanner(p.source)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if file, ok := p.source.(*os.File); ok && file != os.Stdin {
		err := file.Close()
		if err != nil {
			return nil, err
		}
	}
	return lines, nil
}
