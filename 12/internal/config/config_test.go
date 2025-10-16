package config

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
)

func resetFlags() {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}

func TestNew_NoArgs_Panic(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs; resetFlags() }()

	os.Args = []string{"cmd"}
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic when no pattern provided")
		}
	}()

	_ = New()
}

func TestNew_ParseFlagsAndArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs; resetFlags() }()

	os.Args = []string{"cmd", "-A", "2", "-B", "3", "-c", "-i", "-v", "-F", "-n", "pattern", "somefile.txt"}

	cfg := New()

	if cfg.Pattern != "pattern" {
		t.Fatalf("expected pattern 'pattern', got %q", cfg.Pattern)
	}
	if cfg.filePath != "somefile.txt" {
		t.Fatalf("expected filePath 'somefile.txt', got %q", cfg.filePath)
	}
	if cfg.After != 2 || cfg.Before != 3 {
		t.Fatalf("expected After=2 Before=3, got After=%d Before=%d", cfg.After, cfg.Before)
	}
	if !cfg.CountOnly || !cfg.IgnoreCase || !cfg.Invert || !cfg.Fixed || !cfg.LineNum {
		t.Fatalf("expected boolean flags to be true: CountOnly=%v IgnoreCase=%v Invert=%v Fixed=%v LineNum=%v",
			cfg.CountOnly, cfg.IgnoreCase, cfg.Invert, cfg.Fixed, cfg.LineNum)
	}
}

func TestNew_ContextOverridesAB(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs; resetFlags() }()

	os.Args = []string{"cmd", "-C", "4", "pattern"}

	cfg := New()

	if cfg.Context != 4 {
		t.Fatalf("expected Context=4, got %d", cfg.Context)
	}
	if cfg.After != 4 || cfg.Before != 4 {
		t.Fatalf("expected After and Before to equal Context (4), got After=%d Before=%d", cfg.After, cfg.Before)
	}
}

func TestGetSource_FileAndStdin(t *testing.T) {
	tmpDir := t.TempDir()
	fpath := filepath.Join(tmpDir, "testfile.txt")
	content := "hello\nworld\n"
	if err := os.WriteFile(fpath, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	cfg := &Config{filePath: fpath}
	rc, err := cfg.GetSource()
	if err != nil {
		t.Fatalf("GetSource opened file: %v", err)
	}
	defer func() { _ = rc.Close() }()
	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read file source: %v", err)
	}
	if string(b) != content {
		t.Fatalf("file content mismatch: expected %q got %q", content, string(b))
	}

	cfg2 := &Config{filePath: ""}
	rc2, err := cfg2.GetSource()
	if err != nil {
		t.Fatalf("GetSource stdin returned error: %v", err)
	}
	if rc2 != os.Stdin {
		if rc2 == nil {
			t.Fatalf("expected non-nil ReadCloser for stdin")
		}
	}
}
