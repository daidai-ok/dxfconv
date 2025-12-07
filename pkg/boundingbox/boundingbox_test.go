package boundingbox

import (
	"math"
	"testing"
)

func TestNewBoundingBox(t *testing.T) {
	bb := NewBoundingBox()

	if bb.MinX != math.MaxFloat64 {
		t.Errorf("Expected MinX to be MaxFloat64, got %f", bb.MinX)
	}
	if bb.MinY != math.MaxFloat64 {
		t.Errorf("Expected MinY to be MaxFloat64, got %f", bb.MinY)
	}
	if bb.MaxX != -math.MaxFloat64 {
		t.Errorf("Expected MaxX to be -MaxFloat64, got %f", bb.MaxX)
	}
	if bb.MaxY != -math.MaxFloat64 {
		t.Errorf("Expected MaxY to be -MaxFloat64, got %f", bb.MaxY)
	}
}

func TestBoundingBox_Update(t *testing.T) {
	tests := []struct {
		name     string
		points   [][2]float64
		wantMinX float64
		wantMinY float64
		wantMaxX float64
		wantMaxY float64
	}{
		{
			name:     "Single point positive",
			points:   [][2]float64{{10.0, 20.0}},
			wantMinX: 10.0,
			wantMinY: 20.0,
			wantMaxX: 10.0,
			wantMaxY: 20.0,
		},
		{
			name:     "Single point negative",
			points:   [][2]float64{{-5.0, -15.0}},
			wantMinX: -5.0,
			wantMinY: -15.0,
			wantMaxX: -5.0,
			wantMaxY: -15.0,
		},
		{
			name:     "Two points expanding range",
			points:   [][2]float64{{0.0, 0.0}, {100.0, 50.0}},
			wantMinX: 0.0,
			wantMinY: 0.0,
			wantMaxX: 100.0,
			wantMaxY: 50.0,
		},
		{
			name:     "Mixed positive and negative",
			points:   [][2]float64{{-10.0, 10.0}, {20.0, -20.0}},
			wantMinX: -10.0,
			wantMinY: -20.0,
			wantMaxX: 20.0,
			wantMaxY: 10.0,
		},
		{
			name:     "Points within existing bounds",
			points:   [][2]float64{{0.0, 0.0}, {10.0, 10.0}, {5.0, 5.0}},
			wantMinX: 0.0,
			wantMinY: 0.0,
			wantMaxX: 10.0,
			wantMaxY: 10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := NewBoundingBox()
			for _, p := range tt.points {
				bb.Update(p[0], p[1])
			}

			if bb.MinX != tt.wantMinX {
				t.Errorf("Update() MinX = %v, want %v", bb.MinX, tt.wantMinX)
			}
			if bb.MinY != tt.wantMinY {
				t.Errorf("Update() MinY = %v, want %v", bb.MinY, tt.wantMinY)
			}
			if bb.MaxX != tt.wantMaxX {
				t.Errorf("Update() MaxX = %v, want %v", bb.MaxX, tt.wantMaxX)
			}
			if bb.MaxY != tt.wantMaxY {
				t.Errorf("Update() MaxY = %v, want %v", bb.MaxY, tt.wantMaxY)
			}
		})
	}
}

func TestBoundingBox_Dimensions(t *testing.T) {
	bb := NewBoundingBox()
	bb.Update(0, 0)
	bb.Update(100, 50)

	if w := bb.Width(); w != 100 {
		t.Errorf("Width() = %v, want %v", w, 100)
	}
	if h := bb.Height(); h != 50 {
		t.Errorf("Height() = %v, want %v", h, 50)
	}
}
