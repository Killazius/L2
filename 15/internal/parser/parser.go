package parser

import (
	"fmt"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func Parse(input string) ([]*Command, error) {
	parts := splitPipeline(input)
	var commands []*Command

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		fields := strings.Fields(part)
		if len(fields) == 0 {
			continue
		}

		cmd := &Command{
			Name: fields[0],
			Args: fields[1:],
		}
		commands = append(commands, cmd)
	}

	if len(commands) == 0 {
		return nil, fmt.Errorf("no commands found")
	}

	return commands, nil
}

func splitPipeline(input string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	quoteChar := rune(0)

	for _, char := range input {
		switch {
		case char == '"' || char == '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = rune(0)
			}
			current.WriteRune(char)
		case char == '|' && !inQuotes:
			parts = append(parts, current.String())
			current.Reset()
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}
