package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

type Settings struct {
	DebugDataFolder string
	DebugLevel      int
	TexifyDPI       float64
}

type ConvertedSpan struct {
	Text string
	BBox []float64
}

type Page struct {
	// Define your Page structure here
}

type Document struct {
	Name  string
	Pages []Page
}

var settings Settings

func dumpEquationDebugData(doc Document, images []image.Image, convertedSpans []ConvertedSpan) {
	if settings.DebugDataFolder == "" || settings.DebugLevel == 0 {
		return
	}
	if len(images) == 0 {
		return
	}
	// We attempted one conversion per image
	if len(convertedSpans) != len(images) {
		panic("Number of images and converted spans do not match")
	}

	var dataLines []map[string]interface{}

	for idx, img := range images {
		if convertedSpans[idx].Text == "" {
			continue
		}

		// Convert image to WebP
		buf := new(bytes.Buffer)
		err := webp.Encode(buf, img, &encoder.Options{Lossless: true})
		if err != nil {
			fmt.Println("Error encoding image:", err)
			continue
		}

		b64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

		dataLines = append(dataLines, map[string]interface{}{
			"image": b64Image,
			"text":  convertedSpans[idx].Text,
			"bbox":  convertedSpans[idx].BBox,
		})
	}

	// Remove extension from doc name
	docBase := strings.TrimSuffix(filepath.Base(doc.Name), filepath.Ext(doc.Name))
	debugFile := filepath.Join(settings.DebugDataFolder, fmt.Sprintf("%s_equations.json", docBase))

	jsonData, err := json.Marshal(dataLines)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = ioutil.WriteFile(debugFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing debug file:", err)
	}
}

func dumpBBoxDebugData(doc Document, fname string, blocks []Page) {
	if settings.DebugDataFolder == "" || settings.DebugLevel < 2 {
		return
	}

	// Remove extension from doc name
	docBase := strings.TrimSuffix(filepath.Base(fname), filepath.Ext(fname))
	debugFile := filepath.Join(settings.DebugDataFolder, fmt.Sprintf("%s_bbox.json", docBase))

	var debugData []map[string]interface{}

	for idx, pageBlocks := range blocks {
		page := doc.Pages[idx]
		pngImage := renderImage(page, settings.TexifyDPI)
		width, height := pngImage.Bounds().Max.X, pngImage.Bounds().Max.Y

		maxDimension := 6000
		if width > maxDimension || height > maxDimension {
			scalingFactor := math.Min(float64(maxDimension)/float64(width), float64(maxDimension)/float64(height))
			pngImage = imaging.Resize(pngImage, int(float64(width)*scalingFactor), int(float64(height)*scalingFactor), imaging.Lanczos)
		}

		buf := new(bytes.Buffer)
		err := webp.Encode(buf, pngImage, &encoder.Options{Lossless: true, Quality: 100})
		if err != nil {
			fmt.Println("Error encoding image:", err)
			continue
		}

		b64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

		pageData := modelDump(pageBlocks)
		pageData["image"] = b64Image
		debugData = append(debugData, pageData)
	}

	jsonData, err := json.Marshal(debugData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = ioutil.WriteFile(debugFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing debug file:", err)
	}
}

// Placeholder functions - you'll need to implement these
func renderImage(page Page, dpi float64) image.Image {
	// Implement this function based on your needs
	return nil
}

func modelDump(page Page) map[string]interface{} {
	// Implement this function based on your needs
	return nil
}

func main() {
	// Example usage
	settings = Settings{
		DebugDataFolder: "/path/to/debug/folder",
		DebugLevel:      2,
		TexifyDPI:       300,
	}

	// Use the functions here
}
