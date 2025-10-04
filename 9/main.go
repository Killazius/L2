package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func UnpackString(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	r := []rune(s)
	length := len(r)
	if allDigits(r) {
		return "", errors.New("string cannot consist of digits only")
	}
	builder := strings.Builder{}

	builder.Grow(length)
	for l := 0; l < length; l++ {
		curr := r[l]
		if unicode.IsDigit(curr) {
			if l == 0 || unicode.IsDigit(r[l-1]) {
				return "", errors.New("invalid string format")
			}
			continue
		}
		count := 1
		if l+1 < length && unicode.IsDigit(r[l+1]) {
			n, err := strconv.Atoi(string(r[l+1]))
			if err != nil {
				return "", err
			}
			count = n
			l++
		}
		for j := 0; j < count; j++ {
			builder.WriteRune(curr)
		}
	}
	return builder.String(), nil
}

func allDigits(r []rune) bool {
	for _, v := range r {
		if !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}
