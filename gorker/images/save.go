package main

import (
	"fmt"
)

// Page struct to represent the Page type
type Page struct {
	Pnum   int
	Images []Image
}

// Image struct to represent the Image type
type Image struct {
	// Add necessary fields for the Image type
}

// getImageFilename function
func getImageFilename(page Page, imageIdx int) string {
	return fmt.Sprintf("%d_image_%d.png", page.Pnum, imageIdx)
}

// imagesToDict function
func imagesToDict(pages []Page) map[string]Image {
	images := make(map[string]Image)

	for _, page := range pages {
		if page.Images == nil {
			continue
		}

		for imageIdx, image := range page.Images {
			imageFilename := getImageFilename(page, imageIdx)
			images[imageFilename] = image
		}
	}

	return images
}

func main() {
	// Example usage
	pages := []Page{
		{Pnum: 1, Images: []Image{{}, {}}},
		{Pnum: 2, Images: []Image{{}}},
		{Pnum: 3, Images: nil},
	}

	result := imagesToDict(pages)

	// Print the result (just for demonstration)
	for filename, _ := range result {
		fmt.Println(filename)
	}
}
