package main

import (
	"strings"
	"unicode"
)

func alphanumRatio(text string) float64 {
	// Remove spaces and newlines
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\n", "")

	if len(text) == 0 {
		return 1
	}

	alphanumericCount := 0
	for _, c := range text {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			alphanumericCount++
		}
	}

	ratio := float64(alphanumericCount) / float64(len(text))
	return ratio
}
