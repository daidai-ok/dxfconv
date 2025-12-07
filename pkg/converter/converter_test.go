package converter

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/daidai-ok/dxfconv/pkg/dxf"
	"github.com/daidai-ok/dxfconv/pkg/dxfconverror"
)

func TestConvert_PDF(t *testing.T) {
	dxfData, err := os.ReadFile("../../fixtures/overall.dxf")
	if err != nil {
		t.Fatalf("Failed to read encoded DXF: %v", err)
	}
	r := bytes.NewReader(dxfData)
	var w bytes.Buffer
	opts := DefaultOptions()
	opts.Format = FormatPDF

	err = Convert(r, &w, opts)
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
	dxfData, err := os.ReadFile("../../fixtures/overall.dxf")
	if err != nil {
		t.Fatalf("Failed to read encoded DXF: %v", err)
	}
	r := bytes.NewReader(dxfData)
	var w bytes.Buffer
	opts := DefaultOptions()
	opts.Format = FormatSVG

	err = Convert(r, &w, opts)
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

func TestConvert_Arc(t *testing.T) {
	dxfData, err := os.ReadFile("../../fixtures/arc.dxf")
	if err != nil {
		t.Fatalf("Failed to read encoded DXF: %v", err)
	}
	r := bytes.NewReader(dxfData)
	var w bytes.Buffer
	opts := DefaultOptions()
	opts.Format = FormatSVG

	err = Convert(r, &w, opts)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if w.Len() == 0 {
		t.Error("Expected output to be written, but got 0 bytes")
	}

	output := w.String()
	// Basic check for SVG tag
	if !strings.Contains(output, "<svg") {
		t.Error("Expected SVG output to contain <svg tag")
	}
	// svgo Arc uses path
	if !strings.Contains(output, "<path") {
		t.Error("Expected SVG output to contain <path tag for Arc")
	}
	// We expect it to be an arc, so it should have the A command in d attribute
	if !strings.Contains(output, " d=\"M") || !strings.Contains(output, "A") {
		t.Error("Expected SVG path to contain Move (M) and Arc (A) commands")
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

	// The previous external library error wrapping might have been different.
	// Now with our parser, it should fail directly or via bufio.Scanner error handling.
	// Our Parse function returns the scanner error directly if Scan fails unexpectedly.
	// bufio.Scanner treats Read errors as Scan errors.

	// Check if wrapped error contains the message
	if !strings.Contains(err.Error(), "simulated read error") {
		t.Errorf("Expected error measure to contain 'simulated read error', got: %v", err)
	}
}

func TestConvert_BrokenDXF(t *testing.T) {
	dxfData, err := os.ReadFile("../../fixtures/broken.dxf")
	if err != nil {
		t.Fatalf("Failed to read broken DXF: %v", err)
	}
	r := bytes.NewReader(dxfData)
	var w bytes.Buffer
	opts := DefaultOptions()

	err = Convert(r, &w, opts)
	if err == nil {
		t.Fatal("Expected error for broken DXF, but got nil")
	}

	var parseErr *dxfconverror.ParseError
	if !errors.As(err, &parseErr) {
		t.Errorf("Expected ParseError, got %T: %v", err, err)
	}

	// The error message format might have changed with the new parser.
	// Verify it reports the line or something relevant.
	// Our new scanner reports "line N: invalid group code".
	if !strings.Contains(err.Error(), "line") {
		t.Errorf("Expected error measure to contain 'line', got: %v", err)
	}
}

func TestDrawing_Entities(t *testing.T) {
	dxfPath := "../../fixtures/overall.dxf"
	f, err := os.Open(dxfPath)
	if err != nil {
		t.Fatalf("Failed to open DXF from %s: %v", dxfPath, err)
	}
	defer f.Close()

	drawing, err := dxf.Parse(f)
	if err != nil {
		t.Fatalf("Failed to parse DXF: %v", err)
	}

	entities := drawing.Entities
	if len(entities) != 1 {
		t.Fatalf("Expected 1 entity, got %d", len(entities))
	}

	line, ok := entities[0].(*dxf.Line)
	if !ok {
		t.Fatalf("Expected entity type *dxf.Line, got %T", entities[0])
	}

	for i, v := range []float64{0.0, 0.0, 0.0} {
		if line.Start[i] != v {
			t.Errorf("Expected Start[%d] %v, got %v", i, v, line.Start[i])
		}
	}
	for i, v := range []float64{100.0, 100.0, 0.0} {
		if line.End[i] != v {
			t.Errorf("Expected End[%d] %v, got %v", i, v, line.End[i])
		}
	}
}

func TestDrawing_NewEntities(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		check    func(*testing.T, *dxf.Drawing)
	}{
		{
			name:     "Points",
			filename: "../../fixtures/points.dxf",
			check: func(t *testing.T, d *dxf.Drawing) {
				entities := d.Entities
				if len(entities) != 1 {
					t.Fatalf("Expected 1 entity, got %d", len(entities))
				}
				point, ok := entities[0].(*dxf.Point)
				if !ok {
					t.Fatalf("Expected entity type *dxf.Point, got %T", entities[0])
				}
				// Point.Coord is [3]float64
				if point.Coord[0] != 10.0 || point.Coord[1] != 20.0 {
					t.Errorf("Expected Point {10.0, 20.0}, got {%v, %v}", point.Coord[0], point.Coord[1])
				}
			},
		},
		{
			name:     "Polylines",
			filename: "../../fixtures/polylines.dxf",
			check: func(t *testing.T, d *dxf.Drawing) {
				entities := d.Entities
				if len(entities) != 1 {
					t.Fatalf("Expected 1 entity, got %d", len(entities))
				}
				polyline, ok := entities[0].(*dxf.LwPolyline)
				if !ok {
					t.Fatalf("Expected entity type *dxf.LwPolyline, got %T", entities[0])
				}
				if len(polyline.Vertices) != 2 {
					t.Errorf("Expected 2 vertices, got %d", len(polyline.Vertices))
				}
			},
		},
		{
			name:     "Text",
			filename: "../../fixtures/text.dxf",
			check: func(t *testing.T, d *dxf.Drawing) {
				entities := d.Entities
				if len(entities) != 1 {
					t.Fatalf("Expected 1 entity, got %d", len(entities))
				}
				text, ok := entities[0].(*dxf.Text)
				if !ok {
					t.Fatalf("Expected entity type *dxf.Text, got %T", entities[0])
				}
				if text.Value != "Hello World" {
					t.Errorf("Expected Text 'Hello World', got '%v'", text.Value)
				}
			},
		},
		{
			name:     "Arc",
			filename: "../../fixtures/arc.dxf",
			check: func(t *testing.T, d *dxf.Drawing) {
				entities := d.Entities
				if len(entities) != 1 {
					t.Fatalf("Expected 1 entity, got %d", len(entities))
				}
				arc, ok := entities[0].(*dxf.Arc)
				if !ok {
					t.Fatalf("Expected entity type *dxf.Arc, got %T", entities[0])
				}
				if arc.Radius != 100.0 {
					t.Errorf("Expected Radius 100.0, got %v", arc.Radius)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.filename)
			if err != nil {
				t.Fatalf("Failed to open DXF from %s: %v", tt.filename, err)
			}
			defer f.Close()

			drawing, err := dxf.Parse(f)
			if err != nil {
				t.Fatalf("Failed to parse DXF from %s: %v", tt.filename, err)
			}
			tt.check(t, drawing)
		})
	}
}
