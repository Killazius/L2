package config

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *Config
	}{
		{
			name: "default values",
			args: []string{"cmd"},
			expected: &Config{
				Fields:    "",
				Delimiter: "\t",
				Separated: false,
			},
		},
		{
			name: "all flags set",
			args: []string{"cmd", "-f", "1,2,3", "-d", ",", "-s"},
			expected: &Config{
				Fields:    "1,2,3",
				Delimiter: ",",
				Separated: true,
			},
		},
		{
			name: "only separated flag",
			args: []string{"cmd", "-s"},
			expected: &Config{
				Fields:    "",
				Delimiter: "\t",
				Separated: true,
			},
		},
		{
			name: "custom delimiter",
			args: []string{"cmd", "-d", ";"},
			expected: &Config{
				Fields:    "",
				Delimiter: ";",
				Separated: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.args[0], flag.ExitOnError)
			os.Args = tt.args

			cfg := New()

			if !reflect.DeepEqual(cfg, tt.expected) {
				t.Errorf("New() = %v, want %v", cfg, tt.expected)
			}
		})
	}
}
