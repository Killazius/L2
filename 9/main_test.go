package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{"aaa0b", "aab", false},
		{"a10b", "", true},
		{"a1b1c1", "abc", false},
		{"a0b0c0", "", false},
		{"a3b2c1d0e5", "aaabbceeeee", false},
		{"1abc", "", true},
		{"abc12", "", true},
	}

	for _, tt := range tests {
		got, err := UnpackString(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("UnpackString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.expected {
			t.Errorf("UnpackString(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestAllDigits(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"12345", true},
		{"abc", false},
		{"1a2b3", false},
		{"", true},
		{"0", true},
	}

	for _, c := range cases {
		r := []rune(c.input)
		got := allDigits(r)
		if got != c.expected {
			t.Errorf("allDigits(%q) = %v, want %v", c.input, got, c.expected)
		}
	}
}
