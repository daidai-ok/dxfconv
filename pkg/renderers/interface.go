package renderers

// Renderer defines the interface for drawing backend
type Renderer interface {
	// Init initializes the renderer with page dimensions
	Init(width, height float64)
	// Line draws a line segment
	Line(x1, y1, x2, y2 float64)
	// Circle draws a circle
	Circle(x, y, r float64)
	// Arc draws an arc
	Arc(x, y, r, startAngle, endAngle float64)
	// Polyline draws a polyline
	Polyline(points [][]float64, closed bool)
	// Text draws text at the specified location
	Text(x, y, height float64, text string)
	// Finish finalizes the rendering and writes to output
	Finish() error
}
