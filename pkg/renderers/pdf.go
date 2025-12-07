package renderers

import (
	"io"

	"github.com/daidai-ok/dxfconv/pkg/pdf"
)

// PDFRenderer implements the Renderer interface for PDF output
type PDFRenderer struct {
	pdf    *pdf.PDF
	writer io.Writer
}

// NewPDFRenderer creates a new PDFRenderer
func NewPDFRenderer(w io.Writer, orientation string, width, height float64) *PDFRenderer {
	// Simple validation of orientation - our simple PDF assumes dimensions are enough.
	// We ignore orientation flag if width/height are correct, or swap them if needed?
	// The caller (Convert) calculates width/height.

	p := pdf.New(width, height)
	p.AddPage()

	return &PDFRenderer{pdf: p, writer: w}
}

func (r *PDFRenderer) Init(width, height float64) {
	// Already initialized
}

func (r *PDFRenderer) Line(x1, y1, x2, y2 float64) {
	r.pdf.Line(x1, y1, x2, y2)
}

func (r *PDFRenderer) Circle(x, y, radius float64) {
	r.pdf.Circle(x, y, radius)
}

func (r *PDFRenderer) Arc(x, y, radius, startAngle, endAngle float64) {
	r.pdf.Arc(x, y, radius, startAngle, endAngle)
}

func (r *PDFRenderer) Polyline(points [][]float64, closed bool) {
	if len(points) < 2 {
		return
	}

	for i := 0; i < len(points)-1; i++ {
		r.pdf.Line(points[i][0], points[i][1], points[i+1][0], points[i+1][1])
	}
	if closed {
		r.pdf.Line(points[len(points)-1][0], points[len(points)-1][1], points[0][0], points[0][1])
	}
}

// Text draws text at the specified location
func (r *PDFRenderer) Text(x, y, height float64, text string) {
	r.pdf.Text(x, y, height, text)
}

func (r *PDFRenderer) Finish() error {
	return r.pdf.Output(r.writer)
}
