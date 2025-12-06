package dxfconv

import (
	"github.com/yofu/dxf/entity"
)

// drawEntity draws a single DXF entity using the Renderer
func drawEntity(r Renderer, e entity.Entity, scale float64, offsetX, offsetY float64, height float64) {
	transformX := func(x float64) float64 {
		return (x * scale) + offsetX
	}

	transformY := func(y float64) float64 {
		// Flip Y
		return height - ((y * scale) + offsetY)
	}

	switch e := e.(type) {
	case *entity.Line:
		r.Line(transformX(e.Start[0]), transformY(e.Start[1]), transformX(e.End[0]), transformY(e.End[1]))
	case *entity.Circle:
		rad := e.Radius * scale
		r.Circle(transformX(e.Center[0]), transformY(e.Center[1]), rad)
	case *entity.Arc:
		rad := e.Radius * scale
		// DXF angles are in degrees.
		r.Arc(transformX(e.Center[0]), transformY(e.Center[1]), rad, e.Angle[0], e.Angle[1])
	case *entity.LwPolyline:
		if len(e.Vertices) < 2 {
			return
		}
		points := make([][]float64, len(e.Vertices))
		for i, v := range e.Vertices {
			points[i] = []float64{transformX(v[0]), transformY(v[1])}
		}
		r.Polyline(points, e.Closed)
	}
}
