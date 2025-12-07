package converter

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/yofu/dxf"
	"github.com/yofu/dxf/drawing"
	"github.com/yofu/dxf/entity"
)

func TestConvert_PDF(t *testing.T) {
	dxfData, err := os.ReadFile("testdata/overall.dxf")
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
	dxfData, err := os.ReadFile("testdata/overall.dxf")
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

func TestDrawing_Entities(t *testing.T) {
	dxfPath := "testdata/overall.dxf"
	drawing, err := dxf.Open(dxfPath)
	if err != nil {
		t.Fatalf("Failed to open DXF from %s: %v", dxfPath, err)
	}

	entities := drawing.Entities()
	if len(entities) != 1 {
		t.Fatalf("Expected 1 entity, got %d", len(entities))
	}

	line, ok := entities[0].(*entity.Line)
	if !ok {
		t.Fatalf("Expected entity type *entity.Line, got %T", entities[0])
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
		check    func(*testing.T, *drawing.Drawing)
	}{
		{
			name:     "Points",
			filename: "testdata/points.dxf",
			check: func(t *testing.T, d *drawing.Drawing) {
				entities := d.Entities()
				if len(entities) != 1 {
					t.Fatalf("Expected 1 entity, got %d", len(entities))
				}
				point, ok := entities[0].(*entity.Point)
				if !ok {
					t.Fatalf("Expected entity type *entity.Point, got %T", entities[0])
				}
				if len(point.Coord) < 3 {
					t.Fatalf("Expected at least 3 coordinates, got %d", len(point.Coord))
				}
				if point.Coord[0] != 10.0 || point.Coord[1] != 20.0 {
					t.Errorf("Expected Point {10.0, 20.0}, got {%v, %v}", point.Coord[0], point.Coord[1])
				}
			},
		},
		{
			name:     "Polylines",
			filename: "testdata/polylines.dxf",
			check: func(t *testing.T, d *drawing.Drawing) {
				entities := d.Entities()
				if len(entities) != 1 {
					t.Fatalf("Expected 1 entity, got %d", len(entities))
				}
				polyline, ok := entities[0].(*entity.LwPolyline)
				if !ok {
					t.Fatalf("Expected entity type *entity.LwPolyline, got %T", entities[0])
				}
				if len(polyline.Vertices) != 2 {
					t.Errorf("Expected 2 vertices, got %d", len(polyline.Vertices))
				}
			},
		},
		{
			name:     "Text",
			filename: "testdata/text.dxf",
			check: func(t *testing.T, d *drawing.Drawing) {
				entities := d.Entities()
				if len(entities) != 1 {
					t.Fatalf("Expected 1 entity, got %d", len(entities))
				}
				text, ok := entities[0].(*entity.Text)
				if !ok {
					t.Fatalf("Expected entity type *entity.Text, got %T", entities[0])
				}
				if text.Value != "Hello World" {
					t.Errorf("Expected Text 'Hello World', got '%v'", text.Value)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			drawing, err := dxf.Open(tt.filename)
			if err != nil {
				t.Fatalf("Failed to open DXF from %s: %v", tt.filename, err)
			}
			tt.check(t, drawing)
		})
	}
}
