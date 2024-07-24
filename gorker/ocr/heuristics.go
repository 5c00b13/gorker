package main

import (
	"math"
	"regexp"
	"strings"
	"unicode"

	"github.com/your-package/marker/schema"
	"github.com/your-package/marker/settings"
)

func shouldOCRPage(page *schema.Page, noText bool) bool {
	detectedLinesFound, totalLines := detectedLineCoverage(page)

	// No reason to OCR page if it has no text lines
	if totalLines == 0 {
		return false
	}

	// OCR page if we got minimal text, or if we got too many spaces
	conditions := []bool{
		noText, // Full doc has no text, and needs full OCR
		len(page.PrelimText) > 0 && detectBadOCR(page.PrelimText), // Bad OCR
		!detectedLinesFound, // didn't extract text for all detected lines
	}

	for _, condition := range conditions {
		if condition {
			return true
		}
	}

	return settings.OCRAllPages
}

func detectBadOCR(text string, spaceThreshold, newlineThreshold, alphanumThreshold float64) bool {
	if len(text) == 0 {
		// Assume OCR failed if we have no text
		return true
	}

	spaceRegex := regexp.MustCompile(`\s+`)
	spaces := len(spaceRegex.FindAllString(text, -1))
	alphaChars := len(spaceRegex.ReplaceAllString(text, ""))
	if float64(spaces)/float64(alphaChars+spaces) > spaceThreshold {
		return true
	}

	newlineRegex := regexp.MustCompile(`\n+`)
	newlines := len(newlineRegex.FindAllString(text, -1))
	nonNewlines := len(newlineRegex.ReplaceAllString(text, ""))
	if float64(newlines)/float64(newlines+nonNewlines) > newlineThreshold {
		return true
	}

	if alphanumRatio(text) < alphanumThreshold { // Garbled text
		return true
	}

	invalidChars := 0
	for _, c := range text {
		if strings.ContainsRune(settings.InvalidChars, c) {
			invalidChars++
		}
	}
	if float64(invalidChars) > math.Max(6.0, float64(len(text))*0.03) {
		return true
	}

	return false
}

func noTextFound(pages []*schema.Page) bool {
	var fullText strings.Builder
	for _, page := range pages {
		fullText.WriteString(page.PrelimText)
	}
	return len(strings.TrimSpace(fullText.String())) == 0
}

func detectedLineCoverage(page *schema.Page, intersectThresh, detectionThresh float64) (bool, int) {
	foundLines := 0
	for _, detectedLine := range page.TextLines.Bboxes {
		// Get bbox and rescale to match dimensions of original page
		detectedBbox := detectedLine.Bbox
		detectedBbox = rescaleBbox(page.TextLines.ImageBbox, page.Bbox, detectedBbox)
		totalIntersection := 0.0
		for _, block := range page.Blocks {
			for _, line := range block.Lines {
				intersectionPct := boxIntersectionPct(detectedBbox, line.Bbox)
				totalIntersection += intersectionPct
			}
		}
		if totalIntersection > intersectThresh {
			foundLines++
		}
	}
	totalLines := len(page.TextLines.Bboxes)
	if totalLines == 0 {
		return true, 0
	}
	return float64(foundLines)/float64(totalLines) > detectionThresh, totalLines
}

func alphanumRatio(text string) float64 {
	alphanumCount := 0
	for _, c := range text {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			alphanumCount++
		}
	}
	return float64(alphanumCount) / float64(len(text))
}

// These functions are assumed to be defined elsewhere in your Go codebase:
// rescaleBbox
// boxIntersectionPct
