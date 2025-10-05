package main

import (
	"sort"
	"strings"
)

func FindAnagrams(words []string) map[string][]string {
	anagramGroups := make(map[string][]string)
	wordToKey := make(map[string]string)

	for _, word := range words {
		normalized := strings.ToLower(strings.TrimSpace(word))
		runes := []rune(normalized)
		sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
		key := string(runes)
		anagramGroups[key] = append(anagramGroups[key], normalized)
		if _, exists := wordToKey[key]; !exists {
			wordToKey[key] = normalized
		}
	}
	result := make(map[string][]string)
	for key, group := range anagramGroups {
		if len(group) < 2 {
			continue
		}
		sort.Strings(group)
		result[wordToKey[key]] = group
	}
	return result
}
