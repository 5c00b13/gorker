package main

import (
	"image"
	"math"
	"strings"

	"github.com/gen2brain/go-fitz"
	"github.com/otiai10/gosseract/v2"
)

type Page struct {
	TextLines []string
}

type Settings struct {
	DetectorBatchSize *int
	TorchDeviceModel  string
	SuryaDetectorDPI  float64
}

var settings Settings

func getBatchSize() int {
	if settings.DetectorBatchSize != nil {
		return *settings.DetectorBatchSize
	} else if settings.TorchDeviceModel == "cuda" {
		return 4
	}
	return 4
}

func renderImage(doc *fitz.Document, pageNum int, dpi float64) (image.Image, error) {
	return doc.Image(pageNum, dpi, dpi, 0)
}

func batchTextDetection(images []image.Image, client *gosseract.Client, batchSize int) ([][]string, error) {
	var predictions [][]string
	for i := 0; i < len(images); i += batchSize {
		end := int(math.Min(float64(i+batchSize), float64(len(images))))
		batch := images[i:end]

		for _, img := range batch {
			client.SetImageFromBytes(img.Pix)
			text, err := client.Text()
			if err != nil {
				return nil, err
			}
			predictions = append(predictions, strings.Split(text, "\n"))
		}
	}
	return predictions, nil
}

func suryaDetection(doc *fitz.Document, pages []Page, client *gosseract.Client, batchMultiplier float64) error {
	maxLen := int(math.Min(float64(len(pages)), float64(doc.NumPage())))
	var images []image.Image

	for i := 0; i < maxLen; i++ {
		img, err := renderImage(doc, i, settings.SuryaDetectorDPI)
		if err != nil {
			return err
		}
		images = append(images, img)
	}

	batchSize := int(float64(getBatchSize()) * batchMultiplier)
	predictions, err := batchTextDetection(images, client, batchSize)
	if err != nil {
		return err
	}

	for i, pred := range predictions {
		pages[i].TextLines = pred
	}

	return nil
}

func main() {
	// Initialize settings and other necessary setup
	// Use the functions as needed
}
