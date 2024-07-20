package main

import (
	"math"
	"regexp"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy" // for fuzzy string matching
)

type Span struct {
	Text   string
	SpanID string
}

type Line struct {
	Spans []Span
}

type Block struct {
	BlockType string
	Text      string
}

type Page struct {
	Blocks []Block
}

func filterCommonElements(lines []Line, pageCount int, threshold float64) []string {
	if pageCount < 3 {
		return []string{}
	}

	text := []string{}
	for _, line := range lines {
		for _, span := range line.Spans {
			if len(span.Text) > 4 {
				text = append(text, span.Text)
			}
		}
	}

	counter := make(map[string]int)
	for _, t := range text {
		counter[t]++
	}

	common := []string{}
	for k, v := range counter {
		if float64(v) > float64(pageCount)*threshold {
			common = append(common, k)
		}
	}

	badSpanIDs := []string{}
	for _, line := range lines {
		for _, span := range line.Spans {
			for _, c := range common {
				if span.Text == c {
					badSpanIDs = append(badSpanIDs, span.SpanID)
				}
			}
		}
	}

	return badSpanIDs
}

func filterHeaderFooter(allPageBlocks []Page, maxSelectedLines int) []string {
	var firstLines, lastLines []Line

	for _, page := range allPageBlocks {
		nonblankLines := getNonblankLines(page)
		firstLines = append(firstLines, nonblankLines[:int(math.Min(float64(maxSelectedLines), float64(len(nonblankLines))))]...)
		lastLines = append(lastLines, nonblankLines[len(nonblankLines)-int(math.Min(float64(maxSelectedLines), float64(len(nonblankLines)))):])
	}

	badSpanIDs := filterCommonElements(firstLines, len(allPageBlocks), 0.6)
	badSpanIDs = append(badSpanIDs, filterCommonElements(lastLines, len(allPageBlocks), 0.6)...)

	return badSpanIDs
}

func replaceLeadingTrailingDigits(s, replacement string) string {
	re := regexp.MustCompile(`^\d+`)
	s = re.ReplaceAllString(s, replacement)

	re = regexp.MustCompile(`\d+$`)
	s = re.ReplaceAllString(s, replacement)

	return s
}

func findOverlapElements(lst []struct {
	str string
	id  int
}, stringMatchThresh, minOverlap float64) []int {
	result := []int{}
	titles := make([]string, len(lst))
	for i, item := range lst {
		titles[i] = item.str
	}

	for i, item := range lst {
		overlapCount := 0
		for j, str2 := range titles {
			if i != j && fuzzy.RatioNormalized(item.str, str2) >= stringMatchThresh {
				overlapCount++
			}
		}
		if float64(overlapCount) >= math.Max(3.0, float64(len(lst))*minOverlap) {
			result = append(result, item.id)
		}
	}

	return result
}

func filterCommonTitles(mergedBlocks []Block) []Block {
	titles := []struct {
		str string
		id  int
	}{}

	for i, block := range mergedBlocks {
		if block.BlockType == "Title" || block.BlockType == "Section-header" {
			text := block.Text
			if strings.TrimSpace(text)[0] == '#' {
				text = regexp.MustCompile(`^#+`).ReplaceAllString(text, "")
			}
			text = strings.TrimSpace(text)
			text = replaceLeadingTrailingDigits(text, "")
			text = strings.TrimSpace(text)
			titles = append(titles, struct {
				str string
				id  int
			}{text, i})
		}
	}

	badBlockIDs := findOverlapElements(titles, 0.9, 0.05)

	newBlocks := []Block{}
	for i, block := range mergedBlocks {
		if !contains(badBlockIDs, i) {
			newBlocks = append(newBlocks, block)
		}
	}

	return newBlocks
}

// Helper functions

func getNonblankLines(page Page) []Line {
	// Implement this function based on your Page structure
	return []Line{}
}

func contains(slice []int, item int) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func main() {
	// Example usage
	// Initialize your data structures and call the functions here
}
