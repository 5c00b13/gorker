package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type FileStats struct {
	Time  float64 `json:"time"`
	Score float64 `json:"score"`
	Pages int     `json:"pages"`
}

type MethodData struct {
	Files       map[string]FileStats `json:"files"`
	AvgScore    float64              `json:"avg_score"`
	TimePerPage float64              `json:"time_per_page"`
	TimePerDoc  float64              `json:"time_per_doc"`
}

func main() {
	inFolder := flag.String("in_folder", "", "Input PDF files")
	referenceFolder := flag.String("reference_folder", "", "Reference folder with reference markdown files")
	outFile := flag.String("out_file", "", "Output filename")
	nougat := flag.Bool("nougat", false, "Run nougat and compare")
	markerBatchMultiplier := flag.Int("marker_batch_multiplier", 1, "Batch size multiplier to use for marker when making predictions")
	nougatBatchSize := flag.Int("nougat_batch_size", 1, "Batch size to use for nougat when making predictions")
	mdOutPath := flag.String("md_out_path", "", "Output path for generated markdown files")
	profileMemory := flag.Bool("profile_memory", false, "Profile memory usage")

	flag.Parse()

	methods := []string{"marker"}
	if *nougat {
		methods = append(methods, "nougat")
	}

	if *profileMemory {
		startMemoryProfiling()
	}

	modelLst := loadAllModels()

	if *profileMemory {
		stopMemoryProfiling("model_load.pickle")
	}

	scores := make(map[string]map[string]float64)
	times := make(map[string]map[string]float64)
	pages := make(map[string]int)

	benchmarkFiles, err := filepath.Glob(filepath.Join(*inFolder, "*.pdf"))
	if err != nil {
		fmt.Println("Error reading benchmark files:", err)
		return
	}

	for idx, fname := range benchmarkFiles {
		mdFilename := strings.TrimSuffix(filepath.Base(fname), ".pdf") + ".md"
		referenceFilename := filepath.Join(*referenceFolder, mdFilename)

		reference, err := ioutil.ReadFile(referenceFilename)
		if err != nil {
			fmt.Printf("Error reading reference file %s: %v\n", referenceFilename, err)
			continue
		}

		doc, err := openPdfDocument(fname)
		if err != nil {
			fmt.Printf("Error opening PDF %s: %v\n", fname, err)
			continue
		}
		pages[fname] = getPageCount(doc)

		for _, method := range methods {
			start := time.Now()
			var fullText string

			switch method {
			case "marker":
				if *profileMemory {
					startMemoryProfiling()
				}
				fullText, _, _ = convertSinglePdf(fname, modelLst, *markerBatchMultiplier)
				if *profileMemory {
					stopMemoryProfiling(fmt.Sprintf("marker_memory_%d.pickle", idx))
				}
			case "nougat":
				fullText = nougatPrediction(fname, *nougatBatchSize)
			default:
				fmt.Printf("Unknown method: %s\n", method)
				continue
			}

			elapsed := time.Since(start).Seconds()
			if times[method] == nil {
				times[method] = make(map[string]float64)
			}
			times[method][fname] = elapsed

			score := scoreText(fullText, string(reference))
			if scores[method] == nil {
				scores[method] = make(map[string]float64)
			}
			scores[method][fname] = score

			if *mdOutPath != "" {
				mdOutFilename := filepath.Join(*mdOutPath, fmt.Sprintf("%s_%s", method, mdFilename))
				err := ioutil.WriteFile(mdOutFilename, []byte(fullText), 0644)
				if err != nil {
					fmt.Printf("Error writing markdown output file %s: %v\n", mdOutFilename, err)
				}
			}
		}
	}

	writeData := make(map[string]MethodData)
	for _, method := range methods {
		fileStats := make(map[string]FileStats)
		var totalScore, totalTime float64
		var totalPages int

		for _, fname := range benchmarkFiles {
			fileStats[fname] = FileStats{
				Time:  times[method][fname],
				Score: scores[method][fname],
				Pages: pages[fname],
			}
			totalScore += scores[method][fname]
			totalTime += times[method][fname]
			totalPages += pages[fname]
		}

		writeData[method] = MethodData{
			Files:       fileStats,
			AvgScore:    totalScore / float64(len(scores[method])),
			TimePerPage: totalTime / float64(totalPages),
			TimePerDoc:  totalTime / float64(len(scores[method])),
		}
	}

	jsonData, err := json.MarshalIndent(writeData, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = ioutil.WriteFile(*outFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}

	printSummaryTable(writeData, methods)
	printScoreTable(writeData, methods, benchmarkFiles)
}

func printSummaryTable(data map[string]MethodData, methods []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Method", "Average Score", "Time per page", "Time per document"})

	for _, method := range methods {
		methodData := data[method]
		table.Append([]string{
			method,
			fmt.Sprintf("%.4f", methodData.AvgScore),
			fmt.Sprintf("%.4f", methodData.TimePerPage),
			fmt.Sprintf("%.4f", methodData.TimePerDoc),
		})
	}

	table.Render()
}

func printScoreTable(data map[string]MethodData, methods, files []string) {
	table := tablewriter.NewWriter(os.Stdout)
	header := append([]string{"Method"}, files...)
	table.SetHeader(header)

	for _, method := range methods {
		row := []string{method}
		for _, file := range files {
			score := data[method].Files[file].Score
			row = append(row, fmt.Sprintf("%.4f", score))
		}
		table.Append(row)
	}

	fmt.Println("\nScores by file")
	table.Render()
}

// Placeholder functions - these would need to be implemented
func startMemoryProfiling()                       {}
func stopMemoryProfiling(string)                  {}
func loadAllModels() []interface{}                { return nil }
func openPdfDocument(string) (interface{}, error) { return nil, nil }
func getPageCount(interface{}) int                { return 0 }
func convertSinglePdf(string, []interface{}, int) (string, interface{}, interface{}) {
	return "", nil, nil
}
func nougatPrediction(string, int) string { return "" }
func scoreText(string, string) float64    { return 0 }
