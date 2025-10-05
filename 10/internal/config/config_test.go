package config

import (
	"io"
	"os"
	"sync"
	"testing"

	"github.com/spf13/pflag"
)

var testConfigMu sync.Mutex

func TestNew_Defaults(t *testing.T) {
	t.Parallel()
	testConfigMu.Lock()
	origArgs := os.Args
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	os.Args = []string{os.Args[0]}
	testConfigMu.Unlock()
	defer func() {
		testConfigMu.Lock()
		os.Args = origArgs
		testConfigMu.Unlock()
	}()
	cfg := New()
	if cfg.K != 1 {
		t.Errorf("expected K=1, got %d", cfg.K)
	}
	if cfg.N {
		t.Error("expected N=false")
	}
	if cfg.R {
		t.Error("expected R=false")
	}
	if cfg.U {
		t.Error("expected U=false")
	}
	if cfg.Delim != "\t" {
		t.Errorf("expected Delim=\\t, got %q", cfg.Delim)
	}
	if cfg.GetSource() != os.Stdin {
		t.Error("expected GetSource to return os.Stdin")
	}
}

func TestNew_WithFilePath(t *testing.T) {
	t.Parallel()
	testConfigMu.Lock()
	origArgs := os.Args
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	tmpFile, err := os.CreateTemp("", "testfile")
	testConfigMu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	testConfigMu.Lock()
	os.Args = []string{os.Args[0], tmpFile.Name()}
	testConfigMu.Unlock()
	defer func() {
		testConfigMu.Lock()
		os.Args = origArgs
		testConfigMu.Unlock()
	}()
	cfg := New()
	if cfg.GetSource() == os.Stdin {
		t.Error("expected GetSource to return file, not os.Stdin")
	}
}

func TestGetSource_Stdin(t *testing.T) {
	t.Parallel()
	cfg := &Config{}
	r := cfg.GetSource()
	if r != os.Stdin {
		t.Error("expected os.Stdin")
	}
}

func TestGetSource_File(t *testing.T) {
	t.Parallel()
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	cfg := &Config{filePath: tmpFile.Name()}
	r := cfg.GetSource()
	if r == os.Stdin {
		t.Error("expected file, not os.Stdin")
	}
	buf := make([]byte, 1)
	_, err = r.Read(buf)
	if err != nil && err != io.EOF {
		t.Errorf("file read error: %v", err)
	}
}

func TestGetSource_FileOpenError(t *testing.T) {
	t.Parallel()
	cfg := &Config{filePath: "test"}
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when file cannot be opened, but no panic occurred")
		}
	}()
	_ = cfg.GetSource()
}
