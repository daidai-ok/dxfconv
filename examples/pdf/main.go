package main

import (
	"fmt"
	"os"

	dxfconv "github.com/daidai-ok/dxfconv/pkg/converter"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <sample_input.dxf> <sample_output.pdf>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = inputFile.Close() }()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = outputFile.Close() }()

	opts := dxfconv.DefaultOptions()
	opts.Format = dxfconv.FormatPDF
	// Example: Customize options if needed
	// opts.Orientation = dxfconv.OrientationLandscape

	err = dxfconv.Convert(inputFile, outputFile, opts)
	if err != nil {
		fmt.Printf("Error converting DXF to PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputPath, outputPath)
}
