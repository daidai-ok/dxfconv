package dxfconv

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// Minimal DXF content for testing
const validDXF = `  0
SECTION
  2
ENTITIES
  0
LINE
  8
0
 10
0.0
 20
0.0
 11
100.0
 21
100.0
  0
ENDSEC
  0
EOF
`

func TestConvert_PDF(t *testing.T) {
	r := strings.NewReader(validDXF)
	var w bytes.Buffer
	opts := DefaultOptions()
	opts.Format = FormatPDF

	err := Convert(r, &w, opts)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if w.Len() == 0 {
		t.Error("Expected output to be written, but got 0 bytes")
	}

	// Basic check for PDF header
	if !strings.HasPrefix(w.String(), "%PDF") {
		t.Error("Expected PDF output to start with %PDF")
	}
}

func TestConvert_SVG(t *testing.T) {
	r := strings.NewReader(validDXF)
	var w bytes.Buffer
	opts := DefaultOptions()
	opts.Format = FormatSVG

	err := Convert(r, &w, opts)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if w.Len() == 0 {
		t.Error("Expected output to be written, but got 0 bytes")
	}

	// Basic check for SVG tag
	if !strings.Contains(w.String(), "<svg") {
		t.Error("Expected SVG output to contain <svg tag")
	}
}

type FaultyReader struct{}

func (f *FaultyReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated read error")
}

func TestConvert_ReadError(t *testing.T) {
	r := &FaultyReader{}
	var w bytes.Buffer
	opts := DefaultOptions()

	err := Convert(r, &w, opts)
	if err == nil {
		t.Error("Expected error for faulty reader, but got nil")
	}
	if !strings.Contains(err.Error(), "simulated read error") {
		t.Errorf("Expected error measure to contain 'simulated read error', got: %v", err)
	}
}
