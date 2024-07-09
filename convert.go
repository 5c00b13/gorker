package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// Global variables
var modelRefs interface{}

// Settings struct to hold configuration
type Settings struct {
	CUDA               bool
	INFERENCE_RAM      int
	VRAM_PER_TASK      int
	TORCH_DEVICE       string
	TORCH_DEVICE_MODEL string
}

var settings = Settings{
	CUDA:               true,
	INFERENCE_RAM:      8000,
	VRAM_PER_TASK:      1000,
	TORCH_DEVICE:       "cuda",
	TORCH_DEVICE_MODEL: "cuda",
}

func workerInit(sharedModel interface{}) {
	if sharedModel == nil {
		sharedModel = loadAllModels()
	}
	modelRefs = sharedModel
}

func workerExit() {
	modelRefs = nil
}

func processSinglePDF(filepath, outFolder string, metadata map[string]interface{}, minLength int) {
	fname := filepath.Base(filepath)
	if markdownExists(outFolder, fname) {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error converting %s: %v\n", filepath, r)
		}
	}()

	if minLength > 0 {
		filetype := findFiletype(filepath)
		if filetype == "other" {
			return
		}

		length := getLengthOfText(filepath)
		if length < minLength {
			return
		}
	}

	fullText, images, outMetadata := convertSinglePDF(filepath, modelRefs, metadata)
	if len(fullText) > 0 {
		saveMarkdown(outFolder, fname, fullText, images, outMetadata)
	} else {
		fmt.Printf("Empty file: %s. Could not convert.\n", filepath)
	}
}

func main() {
	inFolder := flag.String("in_folder", "", "Input folder with pdfs")
	outFolder := flag.String("out_folder", "", "Output folder")
	chunkIdx := flag.Int("chunk_idx", 0, "Chunk index to convert")
	numChunks := flag.Int("num_chunks", 1, "Number of chunks being processed in parallel")
	maxFiles := flag.Int("max", 0, "Maximum number of pdfs to convert")
	workers := flag.Int("workers", 5, "Number of worker processes to use")
	metadataFile := flag.String("metadata_file", "", "Metadata json file to use for filtering")
	minLength := flag.Int("min_length", 0, "Minimum length of pdf to convert")

	flag.Parse()

	if *inFolder == "" || *outFolder == "" {
		fmt.Println("Both in_folder and out_folder arguments are required")
		flag.Usage()
		os.Exit(1)
	}

	inFolder, _ = filepath.Abs(*inFolder)
	outFolder, _ = filepath.Abs(*outFolder)
	os.MkdirAll(*outFolder, os.ModePerm)

	files, err := ioutil.ReadDir(*inFolder)
	if err != nil {
		fmt.Printf("Error reading input folder: %v\n", err)
		os.Exit(1)
	}

	var filesToConvert []string
	for _, file := range files {
		if !file.IsDir() {
			filesToConvert = append(filesToConvert, filepath.Join(inFolder, file.Name()))
		}
	}

	// Handle chunks if we're processing in parallel
	chunkSize := int(math.Ceil(float64(len(filesToConvert)) / float64(*numChunks)))
	startIdx := *chunkIdx * chunkSize
	endIdx := startIdx + chunkSize
	if endIdx > len(filesToConvert) {
		endIdx = len(filesToConvert)
	}
	filesToConvert = filesToConvert[startIdx:endIdx]

	// Limit files converted if needed
	if *maxFiles > 0 && *maxFiles < len(filesToConvert) {
		filesToConvert = filesToConvert[:*maxFiles]
	}

	metadata := make(map[string]interface{})
	if *metadataFile != "" {
		metadataBytes, err := ioutil.ReadFile(*metadataFile)
		if err != nil {
			fmt.Printf("Error reading metadata file: %v\n", err)
			os.Exit(1)
		}
		json.Unmarshal(metadataBytes, &metadata)
	}

	totalProcesses := *workers
	if settings.CUDA {
		tasksPerGPU := settings.INFERENCE_RAM / settings.VRAM_PER_TASK
		if tasksPerGPU < totalProcesses {
			totalProcesses = tasksPerGPU
		}
	}

	if totalProcesses > len(filesToConvert) {
		totalProcesses = len(filesToConvert)
	}

	if settings.TORCH_DEVICE == "mps" || settings.TORCH_DEVICE_MODEL == "mps" {
		fmt.Println("Cannot use MPS with torch multiprocessing share_memory. This will make things less memory efficient. If you want to share memory, you have to use CUDA or CPU. Set the TORCH_DEVICE environment variable to change the device.")
		modelRefs = nil
	} else {
		modelRefs = loadAllModels()
	}

	fmt.Printf("Converting %d pdfs in chunk %d/%d with %d processes, and storing in %s\n", len(filesToConvert), *chunkIdx+1, *numChunks, totalProcesses, *outFolder)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, totalProcesses)
	bar := progressbar.Default(int64(len(filesToConvert)))

	for _, file := range filesToConvert {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(file string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			processSinglePDF(file, *outFolder, metadata, *minLength)
			bar.Add(1)
		}(file)
	}

	wg.Wait()

	// Clean up
	modelRefs = nil
	runtime.GC()
}

// Placeholder functions - these would need to be implemented
func loadAllModels() interface{}                  { return nil }
func markdownExists(outFolder, fname string) bool { return false }
func findFiletype(filepath string) string         { return "" }
func getLengthOfText(filepath string) int         { return 0 }
func convertSinglePDF(filepath string, modelRefs interface{}, metadata map[string]interface{}) (string, []interface{}, map[string]interface{}) {
	return "", nil, nil
}
func saveMarkdown(outFolder, fname, fullText string, images []interface{}, outMetadata map[string]interface{}) {
}
