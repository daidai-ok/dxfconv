package dxf

import (
	"os"
	"testing"
)

func TestParse_Line(t *testing.T) {
	dxfPath := "../../fixtures/line_simple.dxf"
	f, err := os.Open(dxfPath)
	if err != nil {
		t.Fatalf("Failed to open DXF from %s: %v", dxfPath, err)
	}
	defer f.Close()

	d, err := Parse(f)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(d.Entities) != 1 {
		t.Fatalf("Expected 1 entity, got %d", len(d.Entities))
	}
	line, ok := d.Entities[0].(*Line)
	if !ok {
		t.Fatalf("Expected Line, got %T", d.Entities[0])
	}
	if line.Start[0] != 0 || line.End[0] != 100 {
		t.Errorf("Unexpected coordinate")
	}
}

func TestParse_MText(t *testing.T) {
	dxfPath := "../../fixtures/mtext.dxf"
	f, err := os.Open(dxfPath)
	if err != nil {
		t.Fatalf("Failed to open DXF from %s: %v", dxfPath, err)
	}
	defer f.Close()

	d, err := Parse(f)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(d.Entities) != 1 {
		t.Fatalf("Expected 1 entity, got %d", len(d.Entities))
	}
	mtext, ok := d.Entities[0].(*MText)
	if !ok {
		t.Fatalf("Expected MText, got %T", d.Entities[0])
	}
	if mtext.Value != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", mtext.Value)
	}
	if mtext.Point[0] != 10.0 {
		t.Errorf("Expected X=10.0, got %f", mtext.Point[0])
	}
}

func TestParse_MText_MultiTags(t *testing.T) {
	dxfPath := "../../fixtures/mtext_multi.dxf"
	f, err := os.Open(dxfPath)
	if err != nil {
		t.Fatalf("Failed to open DXF from %s: %v", dxfPath, err)
	}
	defer f.Close()

	d, err := Parse(f)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if len(d.Entities) != 1 {
		t.Fatalf("Expected 1 entity, got %d", len(d.Entities))
	}
	mtext, ok := d.Entities[0].(*MText)
	if !ok {
		t.Fatalf("Expected MText, got %T", d.Entities[0])
	}
	if mtext.Value != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", mtext.Value)
	}
}
