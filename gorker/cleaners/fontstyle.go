package main

import (
	"strings"
)

// Assuming these structs are defined elsewhere in your project
type Span struct {
	Font       string
	FontWeight int
	Bold       bool
	Italic     bool
}

type Line struct {
	Spans []Span
}

type Block struct {
	BlockType string
	Lines     []Line
}

type Page struct {
	Blocks []Block
}

func findBoldItalic(pages []Page, boldMinWeight int) {
	var fontWeights []int

	// First pass: collect font weights and set bold/italic based on font name
	for _, page := range pages {
		for _, block := range page.Blocks {
			// We don't want to bias our font stats
			if block.BlockType == "Title" || block.BlockType == "Section-header" {
				continue
			}
			for i := range block.Lines {
				for j := range block.Lines[i].Spans {
					span := &block.Lines[i].Spans[j]
					if strings.Contains(strings.ToLower(span.Font), "bold") {
						span.Bold = true
					}
					if strings.Contains(strings.ToLower(span.Font), "ital") {
						span.Italic = true
					}
					fontWeights = append(fontWeights, span.FontWeight)
				}
			}
		}
	}

	if len(fontWeights) == 0 {
		return
	}

	// Second pass: set bold based on font weight
	for i := range pages {
		for j := range pages[i].Blocks {
			for k := range pages[i].Blocks[j].Lines {
				for l := range pages[i].Blocks[j].Lines[k].Spans {
					span := &pages[i].Blocks[j].Lines[k].Spans[l]
					if span.FontWeight >= boldMinWeight {
						span.Bold = true
					}
				}
			}
		}
	}
}

func main() {
	// Example usage
	pages := []Page{
		// Initialize your pages here
	}
	findBoldItalic(pages, 600)
}
