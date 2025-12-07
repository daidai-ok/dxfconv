package renderers

import (
	"io"

	"github.com/jung-kurt/gofpdf"
)

// PDFRenderer implements the Renderer interface for PDF output
type PDFRenderer struct {
	pdf        *gofpdf.Fpdf
	writer     io.Writer
	fontFamily string
}

// NewPDFRenderer creates a new PDFRenderer
func NewPDFRenderer(w io.Writer, orientation string, width, height float64, fontPath string) *PDFRenderer {
	pdf := gofpdf.New(orientation, "mm", "A4", "")
	pdf.AddPageFormat(orientation, gofpdf.SizeType{Wd: width, Ht: height})
	pdf.SetAutoPageBreak(false, 0)

	fontFamily := "Arial"
	if fontPath != "" {
		pdf.AddUTF8Font("Custom", "", fontPath)
		fontFamily = "Custom"
	}

	return &PDFRenderer{pdf: pdf, writer: w, fontFamily: fontFamily}
}

func (r *PDFRenderer) Init(width, height float64) {
	// Already initialized
}

func (r *PDFRenderer) Line(x1, y1, x2, y2 float64) {
	r.pdf.Line(x1, y1, x2, y2)
}

func (r *PDFRenderer) Circle(x, y, radius float64) {
	r.pdf.Circle(x, y, radius, "D")
}

func (r *PDFRenderer) Arc(x, y, radius, startAngle, endAngle float64) {
	r.pdf.Arc(x, y, radius, radius, 0, startAngle, endAngle, "D")
}

func (r *PDFRenderer) Polyline(points [][]float64, closed bool) {
	if len(points) < 2 {
		return
	}
	r.pdf.MoveTo(points[0][0], points[0][1])
	for i := 1; i < len(points); i++ {
		r.pdf.LineTo(points[i][0], points[i][1])
	}
	if closed {
		r.pdf.ClosePath()
	}
	r.pdf.DrawPath("D")
}

// Text draws text at the specified location
func (r *PDFRenderer) Text(x, y, height float64, text string) {
	r.pdf.SetFont(r.fontFamily, "", height)
	r.pdf.Text(x, y, text)
}

func (r *PDFRenderer) Finish() error {
	return r.pdf.Output(r.writer)
}
