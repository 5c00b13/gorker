package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Define command-line flags
	inFolder := flag.String("in_folder", "", "Input folder with pdfs")
	outFolder := flag.String("out_folder", "", "Output folder")

	// Parse the flags
	flag.Parse()

	// Check if required flags are provided
	if *inFolder == "" || *outFolder == "" {
		fmt.Println("Both in_folder and out_folder arguments are required")
		flag.Usage()
		os.Exit(1)
	}

	// Get the path of the current executable
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		os.Exit(1)
	}

	// Construct the path to the shell script
	scriptPath := filepath.Join(filepath.Dir(exePath), "chunk_convert.sh")

	// Construct the command
	cmd := exec.Command(scriptPath, *inFolder, *outFolder)

	// Set the command to run in a shell
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the shell script
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error executing shell script: %v\n", err)
		os.Exit(1)
	}
}