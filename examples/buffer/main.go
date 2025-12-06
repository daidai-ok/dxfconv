package main

import (
	"bytes"
	"fmt"
	"os"

	dxfconverter "dxfconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input.dxf>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	// Create a buffer to write to
	var buf bytes.Buffer

	opts := dxfconverter.DefaultOptions()
	// You can choose PDF or SVG
	opts.Format = dxfconverter.FormatPDF

	// Pass the buffer as the io.Writer
	err = dxfconverter.Convert(inputFile, &buf, opts)
	if err != nil {
		fmt.Printf("Error converting to buffer: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted to buffer. Size: %d bytes\n", buf.Len())

	// Verify by writing buffer to file (optional, just to prove it works)
	err = os.WriteFile("buffer_output.pdf", buf.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Error writing buffer to file: %v\n", err)
	}
}
