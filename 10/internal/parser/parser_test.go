package parser

import (
	"10/internal/config"
	"os"
	"strings"
	"testing"
)

func TestParse_FromString(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{}
	p := New(cfg, strings.NewReader("a\nb\nc\n"))
	lines, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "a" || lines[1] != "b" || lines[2] != "c" {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestParse_Empty(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{}
	p := New(cfg, strings.NewReader(""))
	lines, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("expected 0 lines, got %d", len(lines))
	}
}

func TestParse_FileClose(t *testing.T) {
	t.Parallel()
	f, err := os.CreateTemp("", "parser_test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	_, _ = f.WriteString("foo\nbar\n")
	_, err = f.Seek(0, 0)
	if err != nil {
		return
	}
	cfg := &config.Config{}
	p := New(cfg, f)
	lines, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
	if _, err := f.Stat(); err != nil {
		if _, err := f.WriteString("test"); err == nil {
			t.Error("expected error writing to closed file, got nil")
		}
	}
	os.Remove(f.Name())
}

func TestParse_NilReader(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{}
	p := New(cfg, nil)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error when source is nil, got nil")
	}
}
