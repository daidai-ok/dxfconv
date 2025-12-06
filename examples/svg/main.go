package main

import (
	"fmt"
	"os"

	dxfconverter "dxfconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <input.dxf> <output.svg>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	opts := dxfconverter.DefaultOptions()
	opts.Format = dxfconverter.FormatSVG

	err = dxfconverter.Convert(inputFile, outputFile, opts)
	if err != nil {
		fmt.Printf("Error converting DXF to SVG: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputPath, outputPath)
}
