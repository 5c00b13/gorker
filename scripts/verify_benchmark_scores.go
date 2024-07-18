package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type MarkerData struct {
	Marker struct {
		Files struct {
			MulticolCNN struct {
				Score float64 `json:"score"`
			} `json:"multicolcnn.pdf"`
			SwitchTrans struct {
				Score float64 `json:"score"`
			} `json:"switch_trans.pdf"`
		} `json:"files"`
	} `json:"marker"`
}

func verifyScores(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var markerData MarkerData
	if err := json.Unmarshal(data, &markerData); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	multicolcnnScore := markerData.Marker.Files.MulticolCNN.Score
	switchTransScore := markerData.Marker.Files.SwitchTrans.Score

	if multicolcnnScore <= 0.39 || switchTransScore <= 0.4 {
		return fmt.Errorf("one or more scores are below the required threshold of 0.4")
	}

	return nil
}

func main() {
	filePath := flag.String("file", "", "Path to the JSON file")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Please provide a file path using the -file flag")
		os.Exit(1)
	}

	if err := verifyScores(*filePath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Scores verified successfully")
}
