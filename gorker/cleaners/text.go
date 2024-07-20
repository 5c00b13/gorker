package main

import (
	"regexp"
	"strings"
)

func cleanupText(fullText string) string {
	// Replace 3 or more newlines with 2 newlines
	re := regexp.MustCompile(`\n{3,}`)
	fullText = re.ReplaceAllString(fullText, "\n\n")

	// Replace 3 or more occurrences of newline followed by whitespace with 2 newlines
	re = regexp.MustCompile(`(\n\s){3,}`)
	fullText = re.ReplaceAllString(fullText, "\n\n")

	// Replace non-breaking spaces with regular spaces
	fullText = strings.ReplaceAll(fullText, "\u00A0", " ")

	return fullText
}

func main() {
	// Example usage
	text := "This is a\n\n\ntest\n \n \ntext with\u00A0non-breaking\u00A0spaces."
	cleanedText := cleanupText(text)
	println(cleanedText)
}
