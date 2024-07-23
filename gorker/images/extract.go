package main

import (
	"fmt"
)

// Placeholder for settings
var settings struct {
	BBOX_INTERSECTION_THRESH float64
}

// Placeholder structs and types
type Bbox struct {
	// Define bbox properties
}

type Block struct {
	Lines []Line
}

type Line struct {
	Spans []Span
	Bbox  Bbox
}

type Span struct {
	Bbox       Bbox
	Text       string
	Font       string
	Rotation   int
	FontWeight int
	FontSize   int
	Image      bool
	SpanID     string
}

type Page struct {
	Blocks []Block
	Images []Image
	Bbox   Bbox
	Layout struct {
		Bboxes []struct {
			Bbox  Bbox
			Label string
		}
		ImageBbox Bbox
	}
}

type Image struct {
	// Define image properties
}

type PageObj struct {
	// Define page object properties
}

type Doc struct {
	// Define document properties
}

// Placeholder functions
func rescaleBbox(imageBbox, pageBbox, b Bbox) Bbox {
	// Implement bbox rescaling
	return Bbox{}
}

func (l Line) IntersectionPct(region Bbox) float64 {
	// Implement intersection percentage calculation
	return 0
}

func findInsertBlock(blocks []Block, region Bbox) int {
	// Implement find insert block logic
	return 0
}

func renderBboxImage(pageObj PageObj, page Page, bbox Bbox) Image {
	// Implement image rendering
	return Image{}
}

func getImageFilename(page Page, imageIdx int) string {
	// Implement image filename generation
	return fmt.Sprintf("image_%d.png", imageIdx)
}

func findImageBlocks(page Page) [][3]interface{} {
	imageBlocks := [][3]interface{}{}
	imageRegions := []Bbox{}

	for _, l := range page.Layout.Bboxes {
		if l.Label == "Figure" || l.Label == "Picture" {
			imageRegions = append(imageRegions, l.Bbox)
		}
	}

	for i := range imageRegions {
		imageRegions[i] = rescaleBbox(page.Layout.ImageBbox, page.Bbox, imageRegions[i])
	}

	insertPoints := make(map[int][2]int)

	for regionIdx, region := range imageRegions {
		for blockIdx, block := range page.Blocks {
			for lineIdx, line := range block.Lines {
				if line.IntersectionPct(region) > settings.BBOX_INTERSECTION_THRESH {
					line.Spans = []Span{} // We will remove this line from the block
					if _, exists := insertPoints[regionIdx]; !exists {
						insertPoints[regionIdx] = [2]int{blockIdx, lineIdx}
					}
				}
			}
		}
	}

	// Account for images with no detected lines
	for regionIdx, region := range imageRegions {
		if _, exists := insertPoints[regionIdx]; !exists {
			insertPoints[regionIdx] = [2]int{findInsertBlock(page.Blocks, region), 0}
		}
	}

	for regionIdx, imageRegion := range imageRegions {
		imageInsert := insertPoints[regionIdx]
		imageBlocks = append(imageBlocks, [3]interface{}{imageInsert[0], imageInsert[1], imageRegion})
	}

	return imageBlocks
}

func extractPageImages(pageObj PageObj, page *Page) {
	page.Images = []Image{}
	imageBlocks := findImageBlocks(*page)

	for imageIdx, block := range imageBlocks {
		blockIdx := block[0].(int)
		lineIdx := block[1].(int)
		bbox := block[2].(Bbox)

		if blockIdx >= len(page.Blocks) {
			blockIdx = len(page.Blocks) - 1
		}
		if blockIdx < 0 {
			continue
		}

		block := &page.Blocks[blockIdx]
		image := renderBboxImage(pageObj, *page, bbox)
		imageFilename := getImageFilename(*page, imageIdx)
		imageMarkdown := fmt.Sprintf("\n\n![%s](%s)\n\n", imageFilename, imageFilename)

		imageSpan := Span{
			Bbox:       bbox,
			Text:       imageMarkdown,
			Font:       "Image",
			Rotation:   0,
			FontWeight: 0,
			FontSize:   0,
			Image:      true,
			SpanID:     fmt.Sprintf("image_%d", imageIdx),
		}

		if len(block.Lines) > lineIdx {
			block.Lines[lineIdx].Spans = append(block.Lines[lineIdx].Spans, imageSpan)
		} else {
			line := Line{
				Bbox:  bbox,
				Spans: []Span{imageSpan},
			}
			block.Lines = append(block.Lines, line)
		}

		page.Images = append(page.Images, image)
	}
}

func extractImages(doc Doc, pages []Page) {
	for pageIdx := range pages {
		pageObj := PageObj{} // Assuming we can get the page object from the document
		extractPageImages(pageObj, &pages[pageIdx])
	}
}

func main() {
	// Example usage
	doc := Doc{}
	pages := []Page{}
	extractImages(doc, pages)
}
