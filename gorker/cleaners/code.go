package main

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
)

// Structs to represent the document structure
type Span struct {
	Text       string
	BBox       []float64
	SpanID     string
	Font       string
	FontWeight string
	FontSize   float64
}

type Line struct {
	Spans      []Span
	BBox       []float64
	PrelimText string
}

type Block struct {
	BlockType string
	Lines     []Line
	BBox      []float64
}

type Page struct {
	Blocks []Block
}

// Helper functions
func mean(nums []float64) float64 {
	sum := 0.0
	for _, num := range nums {
		sum += num
	}
	return sum / float64(len(nums))
}

func median(nums []float64) float64 {
	sort.Float64s(nums)
	if len(nums)%2 == 0 {
		return (nums[len(nums)/2-1] + nums[len(nums)/2]) / 2
	}
	return nums[len(nums)/2]
}

func isCodeLinelen(lines []Line, thresh float64) bool {
	re := regexp.MustCompile(`\w`)
	totalAlnumChars := 0
	for _, line := range lines {
		totalAlnumChars += len(re.FindAllString(line.PrelimText, -1))
	}
	totalNewlines := math.Max(float64(len(lines)-1), 1)

	if totalAlnumChars == 0 {
		return false
	}

	ratio := float64(totalAlnumChars) / totalNewlines
	return ratio < thresh
}

func commentCount(lines []Line) int {
	pattern := regexp.MustCompile(`^(//|#|'|--|/\*|'''|"""|--\[\[|<!--|%|%{|\(\*)`)
	count := 0
	for _, line := range lines {
		if pattern.MatchString(line.PrelimText) {
			count++
		}
	}
	return count
}

func identifyCodeBlocks(pages []Page) int {
	codeBlockCount := 0
	var fontSizes, lineHeights []float64

	for _, page := range pages {
		fontSizes = append(fontSizes, page.getFontSizes()...)
		lineHeights = append(lineHeights, page.getLineHeights()...)
	}

	var avgFontSize, avgLineHeight float64
	if len(fontSizes) > 0 {
		avgLineHeight = median(lineHeights)
		avgFontSize = mean(fontSizes)
	}

	for _, page := range pages {
		for i := range page.Blocks {
			block := &page.Blocks[i]
			if block.BlockType != "Text" {
				continue
			}

			if len(block.Lines) == 0 {
				continue
			}

			minStart := block.getMinLineStart()

			var isIndent []bool
			var lineFonts []string
			var lineFontSizes, blockLineHeights []float64

			for _, line := range block.Lines {
				for _, span := range line.Spans {
					lineFonts = append(lineFonts, span.Font)
					lineFontSizes = append(lineFontSizes, span.FontSize)
				}
				blockLineHeights = append(blockLineHeights, line.BBox[3]-line.BBox[1])

				isIndent = append(isIndent, line.BBox[0] > minStart)
			}

			commentLines := commentCount(block.Lines)
			isCode := []bool{
				len(block.Lines) > 3,
				isCodeLinelen(block.Lines, 80),
				float64(sum(isIndent)+commentLines) > float64(len(block.Lines))*0.7,
			}

			if avgFontSize != 0 {
				fontChecks := []bool{
					mean(lineFontSizes) <= avgFontSize*0.8,
					mean(blockLineHeights) < avgLineHeight*0.8,
				}
				isCode = append(isCode, fontChecks...)
			}

			if all(isCode) {
				codeBlockCount++
				block.BlockType = "Code"
			}
		}
	}

	return codeBlockCount
}

func indentBlocks(pages []Page) {
	spanCounter := 0
	for _, page := range pages {
		for i := range page.Blocks {
			block := &page.Blocks[i]
			if block.BlockType != "Code" {
				continue
			}

			var lines []struct {
				BBox []float64
				Text string
			}
			minLeft := 1000.0
			colWidth := 0.0

			for _, line := range block.Lines {
				text := ""
				minLeft = math.Min(line.BBox[0], minLeft)
				for _, span := range line.Spans {
					if colWidth == 0 && len(span.Text) > 0 {
						colWidth = (span.BBox[2] - span.BBox[0]) / float64(len(span.Text))
					}
					text += span.Text
				}
				lines = append(lines, struct {
					BBox []float64
					Text string
				}{line.BBox, text})
			}

			blockText := ""
			blankLine := false
			for _, line := range lines {
				text := line.Text
				var prefix string
				if colWidth == 0 {
					prefix = ""
				} else {
					prefix = strings.Repeat(" ", int((line.BBox[0]-minLeft)/colWidth))
				}
				currentLineBlank := len(strings.TrimSpace(text)) == 0
				if blankLine && currentLineBlank {
					continue
				}

				blockText += prefix + text + "\n"
				blankLine = currentLineBlank
			}

			newSpan := Span{
				Text:       blockText,
				BBox:       block.BBox,
				SpanID:     fmt.Sprintf("%d_fix_code", spanCounter),
				Font:       block.Lines[0].Spans[0].Font,
				FontWeight: block.Lines[0].Spans[0].FontWeight,
				FontSize:   block.Lines[0].Spans[0].FontSize,
			}
			spanCounter++
			block.Lines = []Line{{Spans: []Span{newSpan}, BBox: block.BBox}}
		}
	}
}

// Helper functions
func sum(bools []bool) int {
	count := 0
	for _, b := range bools {
		if b {
			count++
		}
	}
	return count
}

func all(bools []bool) bool {
	for _, b := range bools {
		if !b {
			return false
		}
	}
	return true
}

// Methods for Page struct
func (p Page) getFontSizes() []float64 {
	var sizes []float64
	for _, block := range p.Blocks {
		for _, line := range block.Lines {
			for _, span := range line.Spans {
				sizes = append(sizes, span.FontSize)
			}
		}
	}
	return sizes
}

func (p Page) getLineHeights() []float64 {
	var heights []float64
	for _, block := range p.Blocks {
		for _, line := range block.Lines {
			heights = append(heights, line.BBox[3]-line.BBox[1])
		}
	}
	return heights
}

// Method for Block struct
func (b Block) getMinLineStart() float64 {
	minStart := math.Inf(1)
	for _, line := range b.Lines {
		minStart = math.Min(minStart, line.BBox[0])
	}
	return minStart
}

func main() {
	// Your main code here
}
