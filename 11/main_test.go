package main

import (
	"reflect"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		name  string
		words []string
		want  map[string][]string
	}{
		{
			name:  "Basic",
			words: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"},
			want: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			name:  "Empty",
			words: []string{},
			want:  map[string][]string{},
		},
		{
			name:  "NoAnagrams",
			words: []string{"кот", "собака", "мышь"},
			want:  map[string][]string{},
		},
		{
			name:  "CaseAndSpaces",
			words: []string{"Пятак", "  пятка ", "ТЯПКА", "листок", "слиток", "столик"},
			want: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			name:  "Duplicates",
			words: []string{"пятак", "пятак", "пятка", "тяпка"},
			want: map[string][]string{
				"пятак": {"пятак", "пятак", "пятка", "тяпка"},
			},
		},
	}
	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindAnagrams(tt.words)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s: expected %v, got %v", tt.name, tt.want, got)
			}
		})
	}
}
