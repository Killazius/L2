package cut

import (
	"13/internal/config"
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Cut struct {
	cfg    *config.Config
	src    io.Reader
	fields []int
}

func New(cfg *config.Config, src io.Reader) *Cut {
	return &Cut{cfg: cfg, src: src, fields: parseFields(cfg.Fields)}
}

func (c *Cut) Run(dst io.Writer) error {
	scanner := bufio.NewScanner(c.src)
	for scanner.Scan() {
		line := scanner.Text()
		if c.cfg.Separated && !strings.Contains(line, c.cfg.Delimiter) {
			continue
		}
		parts := strings.Split(line, c.cfg.Delimiter)
		var out []string
		for _, idx := range c.fields {
			if idx-1 < len(parts) && idx-1 >= 0 {
				out = append(out, parts[idx-1])
			}
		}
		if len(out) > 0 {
			_, err := dst.Write([]byte(strings.Join(out, c.cfg.Delimiter) + "\n"))
			if err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

func parseFields(s string) []int {
	var res []int
	for _, part := range strings.Split(s, ",") {
		if strings.Contains(part, "-") {
			bounds := strings.Split(part, "-")
			if len(bounds) == 2 {
				start, err1 := strconv.Atoi(bounds[0])
				end, err2 := strconv.Atoi(bounds[1])
				if err1 == nil && err2 == nil && start > 0 && end >= start {
					for i := start; i <= end; i++ {
						res = append(res, i)
					}
				}
			}
		} else {
			idx, err := strconv.Atoi(part)
			if err == nil && idx > 0 {
				res = append(res, idx)
			}
		}
	}
	return res
}
