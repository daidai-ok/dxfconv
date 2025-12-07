package renderers

import (
	"fmt"
	"io"
	"math"

	svg "github.com/ajstarks/svgo"
)

// SVGRenderer implements the Renderer interface for SVG output
type SVGRenderer struct {
	canvas     *svg.SVG
	width      float64
	height     float64
	fontFamily string
}

// NewSVGRenderer creates a new SVGRenderer
func NewSVGRenderer(w io.Writer, width, height float64, font string) *SVGRenderer {
	canvas := svg.New(w)
	fontFamily := "Arial"
	if font != "" {
		fontFamily = font
	}
	return &SVGRenderer{canvas: canvas, width: width, height: height, fontFamily: fontFamily}
}

func (r *SVGRenderer) Init(width, height float64) {
	r.canvas.Start(int(width), int(height))
	r.canvas.Rect(0, 0, int(width), int(height), "fill:none;stroke:none") // Optional background
}

func (r *SVGRenderer) Line(x1, y1, x2, y2 float64) {
	r.canvas.Line(int(x1), int(y1), int(x2), int(y2), "stroke:black;stroke-width:1")
}

func (r *SVGRenderer) Circle(x, y, radius float64) {
	r.canvas.Circle(int(x), int(y), int(radius), "fill:none;stroke:black;stroke-width:1")
}

func (r *SVGRenderer) Arc(x, y, radius, startAngle, endAngle float64) {
	startRad := startAngle * math.Pi / 180
	endRad := endAngle * math.Pi / 180

	sx := int(x + radius*math.Cos(startRad))
	sy := int(y + radius*math.Sin(startRad))
	ex := int(x + radius*math.Cos(endRad))
	ey := int(y + radius*math.Sin(endRad))

	large := false
	diff := endAngle - startAngle
	if diff < 0 {
		diff += 360
	}
	if diff > 180 {
		large = true
	}

	r.canvas.Arc(sx, sy, int(radius), int(radius), 0, large, true, ex, ey, "fill:none;stroke:black")
}

func (r *SVGRenderer) Polyline(points [][]float64, closed bool) {
	x := make([]int, len(points))
	y := make([]int, len(points))
	for i, p := range points {
		x[i] = int(p[0])
		y[i] = int(p[1])
	}
	style := "fill:none;stroke:black;stroke-width:1"
	if closed {
		r.canvas.Polygon(x, y, style)
	} else {
		r.canvas.Polyline(x, y, style)
	}
}

func (r *SVGRenderer) Text(x, y, height float64, text string) {
	r.canvas.Text(int(x), int(y), text, "font-family:"+r.fontFamily+";font-size:"+fmt.Sprintf("%d", int(height)))
}

func (r *SVGRenderer) Finish() error {
	r.canvas.End()
	return nil
}
