package main

import (
	"regexp"
)

func replaceBullets(text string) string {
	// Replace bullet characters with a -
	bulletPattern := `(^|[\n ])[•●○■▪▫–—]( )`
	re := regexp.MustCompile(bulletPattern)
	replacedString := re.ReplaceAllString(text, "$1-$2")
	return replacedString
}
