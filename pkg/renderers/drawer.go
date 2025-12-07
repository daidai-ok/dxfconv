package renderers

import (
	"github.com/daidai-ok/dxfconv/pkg/dxf"
)

// DrawEntity draws a single DXF entity using the Renderer
func DrawEntity(r Renderer, e dxf.Entity, scale float64, offsetX, offsetY float64, height float64) {
	transformX := func(x float64) float64 {
		return (x * scale) + offsetX
	}

	transformY := func(y float64) float64 {
		// Flip Y
		return height - ((y * scale) + offsetY)
	}

	switch e := e.(type) {
	case *dxf.Line:
		r.Line(transformX(e.Start[0]), transformY(e.Start[1]), transformX(e.End[0]), transformY(e.End[1]))
	case *dxf.Circle:
		rad := e.Radius * scale
		r.Circle(transformX(e.Center[0]), transformY(e.Center[1]), rad)
	case *dxf.Arc:
		rad := e.Radius * scale
		// DXF angles are in degrees.
		r.Arc(transformX(e.Center[0]), transformY(e.Center[1]), rad, e.StartAngle, e.EndAngle)
	case *dxf.LwPolyline:
		if len(e.Vertices) < 2 {
			return
		}
		points := make([][]float64, len(e.Vertices))
		for i, v := range e.Vertices {
			points[i] = []float64{transformX(v.X), transformY(v.Y)}
		}
		r.Polyline(points, e.Closed)
	case *dxf.Polyline:
		if len(e.Vertices) < 2 {
			return
		}
		points := make([][]float64, len(e.Vertices))
		for i, v := range e.Vertices {
			points[i] = []float64{transformX(v.X), transformY(v.Y)}
		}
		// Closed flag is already handled in parser
		r.Polyline(points, e.Closed)
	case *dxf.Spline:
		if len(e.ControlPoints) < 2 {
			return
		}
		// Approximation by control points
		points := make([][]float64, len(e.ControlPoints))
		for i, v := range e.ControlPoints {
			points[i] = []float64{transformX(v[0]), transformY(v[1])}
		}
		r.Polyline(points, e.Closed) // Spline can be closed
	case *dxf.Point:
		// Draw as a small circle, simplistic representation
		radius := 1.0 * scale // Fixed visual size or scaled
		r.Circle(transformX(e.Coord[0]), transformY(e.Coord[1]), radius)
	case *dxf.Text:
		r.Text(transformX(e.Point[0]), transformY(e.Point[1]), e.Height*scale, e.Value)
	case *dxf.MText:
		// Handling MText similarly to Text for now
		r.Text(transformX(e.Point[0]), transformY(e.Point[1]), e.Height*scale, e.Value)
	}
}
