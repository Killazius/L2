package sorter

import (
	"10/internal/config"
	"10/internal/parser"
	"sort"
	"strconv"
	"strings"
)

type Sorter struct {
	parser *parser.Parser
	cfg    *config.Config
}

func New(p *parser.Parser, cfg *config.Config) *Sorter {
	return &Sorter{
		parser: p,
		cfg:    cfg,
	}
}

func (s *Sorter) Sort() ([]string, error) {
	lines, err := s.parser.Parse()
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return lines, nil
	}

	k := s.cfg.K - 1
	compare := func(a, b string) int {
		acol := ""
		bcol := ""
		as := strings.Split(a, s.cfg.Delim)
		bs := strings.Split(b, s.cfg.Delim)
		if k < len(as) {
			acol = as[k]
		}
		if k < len(bs) {
			bcol = bs[k]
		}
		if s.cfg.N {
			af, aerr := strconv.ParseFloat(acol, 64)
			bf, berr := strconv.ParseFloat(bcol, 64)
			if aerr == nil && berr == nil {
				if af < bf {
					return -1
				} else if af > bf {
					return 1
				}
				return 0
			}
		}
		if acol < bcol {
			return -1
		} else if acol > bcol {
			return 1
		}
		return 0
	}

	sort.SliceStable(lines, func(i, j int) bool {
		cmp := compare(lines[i], lines[j])
		if s.cfg.R {
			return cmp > 0
		}
		return cmp < 0
	})

	if s.cfg.U {
		uniq := make([]string, 0, len(lines))
		prev := ""
		for i, line := range lines {
			if i == 0 || line != prev {
				uniq = append(uniq, line)
				prev = line
			}
		}
		lines = uniq
	}

	return lines, nil
}
