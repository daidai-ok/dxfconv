package dxfconv

import (
	"math"
)

// BoundingBox represents the min and max coordinates of the drawing
type BoundingBox struct {
	MinX, MinY, MaxX, MaxY float64
}

// NewBoundingBox initializes a BoundingBox with extreme values
func NewBoundingBox() *BoundingBox {
	return &BoundingBox{
		MinX: math.MaxFloat64,
		MinY: math.MaxFloat64,
		MaxX: -math.MaxFloat64,
		MaxY: -math.MaxFloat64,
	}
}

// Update expands the bounding box to include the given point
func (bb *BoundingBox) Update(x, y float64) {
	if x < bb.MinX {
		bb.MinX = x
	}
	if y < bb.MinY {
		bb.MinY = y
	}
	if x > bb.MaxX {
		bb.MaxX = x
	}
	if y > bb.MaxY {
		bb.MaxY = y
	}
}

// Width returns the width of the bounding box
func (bb *BoundingBox) Width() float64 {
	return bb.MaxX - bb.MinX
}

// Height returns the height of the bounding box
func (bb *BoundingBox) Height() float64 {
	return bb.MaxY - bb.MinY
}
