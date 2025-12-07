package renderers

import (
	"github.com/yofu/dxf/entity"
)

// drawEntity draws a single DXF entity using the Renderer
func DrawEntity(r Renderer, e entity.Entity, scale float64, offsetX, offsetY float64, height float64) {
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
	case *entity.Polyline:
		if len(e.Vertices) < 2 {
			return
		}
		points := make([][]float64, len(e.Vertices))
		for i, v := range e.Vertices {
			points[i] = []float64{transformX(v.Coord[0]), transformY(v.Coord[1])}
		}
		// Flag & 1 indicates closed polyline
		closed := e.Flag&1 == 1
		r.Polyline(points, closed)
	case *entity.Spline:
		if len(e.Controls) < 2 {
			return
		}
		// Approximation by control points
		points := make([][]float64, len(e.Controls))
		for i, v := range e.Controls {
			points[i] = []float64{transformX(v[0]), transformY(v[1])}
		}
		r.Polyline(points, false)
	case *entity.Point:
		// Draw as a small circle, simplistic representation
		radius := 1.0 * scale // Fixed visual size or scaled
		r.Circle(transformX(e.Coord[0]), transformY(e.Coord[1]), radius)
	case *entity.Text:
		r.Text(transformX(e.Coord1[0]), transformY(e.Coord1[1]), e.Height*scale, e.Value)
	}
}
