package sorter

import (
	"10/internal/config"
	"10/internal/parser"
	"strings"
	"testing"
)

func newTestSorter(input string, cfg *config.Config) *Sorter {
	p := parser.New(cfg, strings.NewReader(input))
	return New(p, cfg)
}

func TestSort_Default(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{K: 1, Delim: "\t"}
	s := newTestSorter("b\na\nc\n", cfg)
	res, err := s.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"a", "b", "c"}
	for i := range want {
		if res[i] != want[i] {
			t.Errorf("expected %q at %d, got %q", want[i], i, res[i])
		}
	}
}

func TestSort_Column(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{K: 2, Delim: ","}
	input := "a,2\nb,1\nc,3\n"
	s := newTestSorter(input, cfg)
	res, err := s.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"b,1", "a,2", "c,3"}
	for i := range want {
		if res[i] != want[i] {
			t.Errorf("expected %q at %d, got %q", want[i], i, res[i])
		}
	}
}

func TestSort_Numeric(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{K: 2, Delim: ",", N: true}
	input := "a,10\nb,2\nc,1\n"
	s := newTestSorter(input, cfg)
	res, err := s.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"c,1", "b,2", "a,10"}
	for i := range want {
		if res[i] != want[i] {
			t.Errorf("expected %q at %d, got %q", want[i], i, res[i])
		}
	}
}

func TestSort_Reverse(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{K: 1, Delim: "\t", R: true}
	s := newTestSorter("a\nb\nc\n", cfg)
	res, err := s.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"c", "b", "a"}
	for i := range want {
		if res[i] != want[i] {
			t.Errorf("expected %q at %d, got %q", want[i], i, res[i])
		}
	}
}

func TestSort_Unique(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{K: 1, Delim: "\t", U: true}
	s := newTestSorter("a\na\nb\nb\nc\n", cfg)
	res, err := s.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"a", "b", "c"}
	for i := range want {
		if res[i] != want[i] {
			t.Errorf("expected %q at %d, got %q", want[i], i, res[i])
		}
	}
}

func TestSort_Empty(t *testing.T) {
	t.Parallel()
	cfg := &config.Config{K: 1, Delim: "\t"}
	s := newTestSorter("", cfg)
	res, err := s.Sort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected 0 lines, got %d", len(res))
	}
}
