package main

// Assuming these types are defined elsewhere in your project
type BBox struct {
	// Define BBox fields
}

type Line struct {
	// Define Line fields
}

type Block struct {
	BlockType string
	Lines     []Line
	BBox      BBox
}

type Page struct {
	Blocks []Block
	Layout Layout
	BBox   BBox
}

type Layout struct {
	BBoxes    []LayoutBBox
	ImageBBox BBox
}

type LayoutBBox struct {
	BBox  BBox
	Label string
}

// Assuming these functions are defined elsewhere
func rescaleBBox(imageBBox, pageBBox, bbox BBox) BBox {
	// Implementation of rescaleBBox
}

func bboxFromLines(lines []Line) BBox {
	// Implementation of bboxFromLines
}

func (l Line) intersectionPct(bbox BBox) float64 {
	// Implementation of intersection percentage calculation
}

func (b Block) copy() Block {
	// Implementation of block copy
}

const BBoxIntersectionThresh = 0.5 // Assuming this is defined in settings

func splitHeadingBlocks(pages []Page) {
	for i := range pages {
		page := &pages[i]
		var pageHeadingBoxes []struct {
			bbox  BBox
			label string
		}

		for _, b := range page.Layout.BBoxes {
			if b.Label == "Title" || b.Label == "Section-header" {
				rescaledBBox := rescaleBBox(page.Layout.ImageBBox, page.BBox, b.BBox)
				pageHeadingBoxes = append(pageHeadingBoxes, struct {
					bbox  BBox
					label string
				}{rescaledBBox, b.Label})
			}
		}

		var newBlocks []Block
		for _, block := range page.Blocks {
			if block.BlockType != "Text" {
				newBlocks = append(newBlocks, block)
				continue
			}

			var headingLines []struct {
				index int
				label string
			}

			for lineIdx, line := range block.Lines {
				for _, headingBox := range pageHeadingBoxes {
					if line.intersectionPct(headingBox.bbox) > BBoxIntersectionThresh {
						headingLines = append(headingLines, struct {
							index int
							label string
						}{lineIdx, headingBox.label})
						break
					}
				}
			}

			if len(headingLines) == 0 {
				newBlocks = append(newBlocks, block)
				continue
			}

			start := 0
			for _, headingLine := range headingLines {
				if start < headingLine.index {
					copiedBlock := block.copy()
					copiedBlock.Lines = block.Lines[start:headingLine.index]
					copiedBlock.BBox = bboxFromLines(copiedBlock.Lines)
					newBlocks = append(newBlocks, copiedBlock)
				}

				copiedBlock := block.copy()
				copiedBlock.Lines = block.Lines[headingLine.index : headingLine.index+1]
				copiedBlock.BlockType = headingLine.label
				copiedBlock.BBox = bboxFromLines(copiedBlock.Lines)
				newBlocks = append(newBlocks, copiedBlock)

				start = headingLine.index + 1
				if start >= len(block.Lines) {
					break
				}
			}

			if start < len(block.Lines) {
				copiedBlock := block.copy()
				copiedBlock.Lines = block.Lines[start:]
				copiedBlock.BBox = bboxFromLines(copiedBlock.Lines)
				newBlocks = append(newBlocks, copiedBlock)
			}
		}

		page.Blocks = newBlocks
	}
}

func main() {
	// Example usage
	var pages []Page
	// Initialize pages
	splitHeadingBlocks(pages)
}
