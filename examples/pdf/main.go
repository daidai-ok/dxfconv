package main

import (
	"fmt"
	"os"

	dxfconverter "dxfconv/pkg/converter"
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

	opts := dxfconverter.DefaultOptions()
	opts.Format = dxfconverter.FormatPDF
	// Example: Customize options if needed
	// opts.Orientation = dxfconverter.OrientationLandscape

	err = dxfconverter.Convert(inputFile, outputFile, opts)
	if err != nil {
		fmt.Printf("Error converting DXF to PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputPath, outputPath)
}
