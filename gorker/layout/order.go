package main

import (
	"sort"
)

type Page struct {
	Layout Layout
	Blocks []Block
	Order  Order
	Bbox   BoundingBox
}

type Layout struct {
	Bboxes []LayoutBlock
}

type LayoutBlock struct {
	Bbox BoundingBox
}

type Order struct {
	Bboxes    []OrderBox
	ImageBbox BoundingBox
}

type OrderBox struct {
	Bbox     BoundingBox
	Position int
}

type Block struct {
	Bbox BoundingBox
}

type BoundingBox struct {
	// Define BoundingBox structure
}

type Settings struct {
	OrderBatchSize   *int
	TorchDeviceModel string
	SuryaOrderDPI    int
	OrderMaxBboxes   int
}

var settings Settings

func getBatchSize() int {
	if settings.OrderBatchSize != nil {
		return *settings.OrderBatchSize
	} else if settings.TorchDeviceModel == "cuda" {
		return 6
	} else if settings.TorchDeviceModel == "mps" {
		return 6
	}
	return 6
}

func suryaOrder(doc interface{}, pages []Page, orderModel interface{}, batchMultiplier float64) {
	images := make([]Image, len(pages))
	for i := range pages {
		images[i] = renderImage(doc, i, settings.SuryaOrderDPI)
	}

	bboxes := make([][]BoundingBox, len(pages))
	for i, page := range pages {
		bbox := make([]BoundingBox, 0, settings.OrderMaxBboxes)
		for _, b := range page.Layout.Bboxes {
			if len(bbox) >= settings.OrderMaxBboxes {
				break
			}
			bbox = append(bbox, b.Bbox)
		}
		bboxes[i] = bbox
	}

	processor := getProcessor(orderModel)
	batchSize := int(float64(getBatchSize()) * batchMultiplier)
	orderResults := batchOrdering(images, bboxes, orderModel, processor, batchSize)

	for i, page := range pages {
		page.Order = orderResults[i]
	}
}

func sortBlocksInReadingOrder(pages []Page) {
	for _, page := range pages {
		order := page.Order
		blockPositions := make(map[int]struct {
			intersection float64
			position     int
		})
		maxPosition := 0

		for i, block := range page.Blocks {
			for _, orderBox := range order.Bboxes {
				orderBbox := orderBox.Bbox
				position := orderBox.Position
				orderBbox = rescaleBbox(order.ImageBbox, page.Bbox, orderBbox)
				blockIntersection := block.IntersectionPct(orderBbox)

				if _, ok := blockPositions[i]; !ok || blockIntersection > blockPositions[i].intersection {
					blockPositions[i] = struct {
						intersection float64
						position     int
					}{blockIntersection, position}
				}
				maxPosition = max(maxPosition, position)
			}
		}

		blockGroups := make(map[int][]Block)
		for i, block := range page.Blocks {
			position := maxPosition + 1
			if pos, ok := blockPositions[i]; ok {
				position = pos.position
			} else {
				maxPosition++
			}
			blockGroups[position] = append(blockGroups[position], block)
		}

		var positions []int
		for position := range blockGroups {
			positions = append(positions, position)
		}
		sort.Ints(positions)

		var newBlocks []Block
		for _, position := range positions {
			blockGroup := sortBlockGroup(blockGroups[position])
			newBlocks = append(newBlocks, blockGroup...)
		}
		page.Blocks = newBlocks
	}
}

// Helper functions (not implemented, just signatures)
func renderImage(doc interface{}, pageNum int, dpi int) Image {
	// Implementation
}

func getProcessor(orderModel interface{}) interface{} {
	// Implementation
}

func batchOrdering(images []Image, bboxes [][]BoundingBox, orderModel interface{}, processor interface{}, batchSize int) []Order {
	// Implementation
}

func rescaleBbox(imageBbox, pageBbox, layoutBbox BoundingBox) BoundingBox {
	// Implementation
}

func (b Block) IntersectionPct(other BoundingBox) float64 {
	// Implementation
}

func sortBlockGroup(blocks []Block) []Block {
	// Implementation
}

type Image struct {
	// Define Image structure
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
