package pdf

import (
	"bytes"
	"strings"
	"testing"
)

func TestPDF_Line(t *testing.T) {
	p := New(100, 100) // 100x100 canvas
	p.Line(10, 10, 90, 90)

	got := p.currentBuf.String()
	// Y is flipped: 100 - 10 = 90, 100 - 90 = 10
	want := "10.00 90.00 m 90.00 10.00 l S\n"

	if got != want {
		t.Errorf("Line() got = %q, want %q", got, want)
	}
}

func TestPDF_Circle(t *testing.T) {
	p := New(100, 100)
	p.Circle(50, 50, 10)

	got := p.currentBuf.String()
	// Just check if it contains expected commands as the exact float string might be long
	// Center is (50, 100-50=50)
	// Radius 10
	// We expect multiple 'c' (curve) commands and 'S' (stroke) at the end.
	if !strings.Contains(got, "c") {
		t.Errorf("Circle() output should contain bezier curves 'c', got %q", got)
	}
	if !strings.HasSuffix(got, "S\n") {
		t.Errorf("Circle() output should end with stroke 'S', got %q", got)
	}
}

func TestPDF_Arc(t *testing.T) {
	p := New(100, 100)
	// Arc from 0 to 90 degrees, radius 10 at (50,50)
	p.Arc(50, 50, 10, 0, 90)

	got := p.currentBuf.String()
	// Should contain move 'm', line 'l', and stroke 'S'
	if !strings.Contains(got, "m") {
		t.Errorf("Arc() output should contain move 'm', got %q", got)
	}
	if !strings.Contains(got, "l") {
		t.Errorf("Arc() output should contain line 'l', got %q", got)
	}
	if !strings.HasSuffix(got, "S\n") {
		t.Errorf("Arc() output should end with stroke 'S', got %q", got)
	}

	// Verify start point
	// 0 degrees: x=50+10=60, y=50 (flipped y=100-50=50)
	// PDF coords: 60.00 50.00 m
	wantStart := "60.00 50.00 m"
	if !strings.Contains(got, wantStart) {
		t.Errorf("Arc() output should start at %q, got %q", wantStart, got)
	}

	// Verify end point
	// 90 degrees: x=50, y=50+10=60 (flipped y=100-60=40)
	// PDF coords: 50.00 40.00 l S
	wantEnd := "50.00 40.00 l S"
	if !strings.Contains(got, wantEnd) {
		t.Errorf("Arc() output should end at %q, got %q", wantEnd, got)
	}
}

func TestPDF_Text(t *testing.T) {
	p := New(100, 100)
	p.Text(10, 20, 12, "Hello World")

	got := p.currentBuf.String()
	// x=10, y=20 -> flipped y=80
	// BT /F1 12.00 Tf 10.00 80.00 Td (Hello World) Tj ET
	want := "BT /F1 12.00 Tf 10.00 80.00 Td (Hello World) Tj ET\n"

	if got != want {
		t.Errorf("Text() got = %q, want %q", got, want)
	}
}

func TestPDF_Output(t *testing.T) {
	p := New(100, 200)
	p.Line(0, 0, 100, 200)

	var buf bytes.Buffer
	err := p.Output(&buf)
	if err != nil {
		t.Fatalf("Output() error = %v", err)
	}

	output := buf.String()

	// Basic structural checks
	checks := []string{
		"%PDF-1.4",
		"1 0 obj", // Catalog
		"/Type /Catalog",
		"2 0 obj", // Pages
		"/Type /Pages",
		"3 0 obj", // Page
		"/MediaBox [0 0 100.00 200.00]",
		"4 0 obj", // Content stream
		"stream",
		"endstream",
		"slider", // Wait, slider is not expected
		"trailer",
		"%%EOF",
	}

	for _, check := range checks {
		if check == "slider" {
			continue
		} // Just a placeholder in my loop
		if !strings.Contains(output, check) {
			t.Errorf("Output() missing %q", check)
		}
	}
}

func TestPDF_AddPage(t *testing.T) {
	p := New(100, 100)
	p.Line(0, 0, 10, 10)
	if p.currentBuf.Len() == 0 {
		t.Error("Buffer should not be empty after drawing")
	}

	p.AddPage()
	if p.currentBuf.Len() != 0 {
		t.Error("Buffer should be empty after AddPage")
	}
}
