package pdf

import (
	"bytes"
	"fmt"
	"io"
	"math"
)

// PDF represents a simple PDF generator
type PDF struct {
	width      float64
	height     float64
	currentBuf bytes.Buffer
}

// New creates a new PDF generator
func New(width, height float64) *PDF {
	p := &PDF{
		width:  width,
		height: height,
	}
	return p
}

// AddPage adds a new page.
// In this simplified implementation, we assume a single page and reset the buffer.
func (p *PDF) AddPage() {
	// Reset buffer for new page content
	p.currentBuf.Reset()
}

// Line draws a line
func (p *PDF) Line(x1, y1, x2, y2 float64) {
	// PDF coordinates start at the bottom-left (0,0).
	// We flip the Y coordinate to match the top-left origin used by the converter.
	p.currentBuf.WriteString(fmt.Sprintf("%.2f %.2f m %.2f %.2f l S\n", x1, p.height-y1, x2, p.height-y2))
}

// Circle draws a circle
func (p *PDF) Circle(x, y, r float64) {
	// Approximation with Bezier curves using 4 segments
	k := 0.552284749831 * r
	cx, cy := x, p.height-y // center (flip Y for PDF coordinates)

	// Let's just draw standard circle around cx, cy
	// P1: (r, 0)
	// C1: (r, k)
	// C2: (k, r)
	// P2: (0, r)
	p.currentBuf.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f %.2f %.2f c\n", cx+r, cy+k, cx+k, cy+r, cx, cy+r))

	// P2: (0, r) -> P3: (-r, 0)
	p.currentBuf.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f %.2f %.2f c\n", cx-k, cy+r, cx-r, cy+k, cx-r, cy))

	// P3: (-r, 0) -> P4: (0, -r)
	p.currentBuf.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f %.2f %.2f c\n", cx-r, cy-k, cx-k, cy-r, cx, cy-r))

	// P4: (0, -r) -> P1: (r, 0)
	p.currentBuf.WriteString(fmt.Sprintf("%.2f %.2f %.2f %.2f %.2f %.2f c S\n", cx+k, cy-r, cx+r, cy-k, cx+r, cy))
}

// Arc draws an arc
func (p *PDF) Arc(x, y, r, startAngle, endAngle float64) {
	// This is complex in PDF (requires bezier approximation).
	// For now, just implement the math for a simple arc or skip if too complex for this step.
	// Approximate with small line segments.

	step := 5.0 // degrees
	cx, cy := x, p.height-y

	startRad := startAngle * math.Pi / 180
	endRad := endAngle * math.Pi / 180

	// Handling angle wrap around could be needed if start > end, but let's assume valid range or simple loop
	if endRad < startRad {
		endRad += 2 * math.Pi
	}

	// Move to start
	p.currentBuf.WriteString(fmt.Sprintf("%.2f %.2f m\n", cx+r*math.Cos(startRad), cy+r*math.Sin(startRad)))

	for a := startRad; a <= endRad; a += step * math.Pi / 180 {
		p.currentBuf.WriteString(fmt.Sprintf("%.2f %.2f l\n", cx+r*math.Cos(a), cy+r*math.Sin(a)))
	}
	// Final point
	p.currentBuf.WriteString(fmt.Sprintf("%.2f %.2f l S\n", cx+r*math.Cos(endRad), cy+r*math.Sin(endRad)))
}

// Text draws text
func (p *PDF) Text(x, y, size float64, text string) {
	// BT /F1 size Tf x y Td (text) Tj ET
	// Escape text parens
	// y needs flip
	// Hard code /F1 for now as we standardizing on Helvetica
	p.currentBuf.WriteString(fmt.Sprintf("BT /F1 %.2f Tf %.2f %.2f Td (%s) Tj ET\n", size, x, p.height-y, text))
}

// output writes the PDF to the writer
func (p *PDF) Output(w io.Writer) error {
	// Object Map:
	// 1: Catalog
	// 2: Pages
	// 3: Page
	// 4: Content Stream
	// 5: Font (Helvetica)

	// We need to build the objects first to calculate offsets
	var objects []string

	// 1. Catalog
	objects = append(objects, "<< /Type /Catalog /Pages 2 0 R >>")

	// 2. Pages
	objects = append(objects, "<< /Type /Pages /Kids [3 0 R] /Count 1 >>")

	// 3. Page
	objects = append(objects, fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %.2f %.2f] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>", p.width, p.height))

	// 4. Content Stream
	stream := p.currentBuf.String()
	objects = append(objects, fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(stream), stream))

	// 5. Font
	objects = append(objects, "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>")

	// Write Header
	n, err := w.Write([]byte("%PDF-1.4\n"))
	if err != nil {
		return err
	}
	offset := n

	// Write Objects
	var offsets []int
	for i, obj := range objects {
		offsets = append(offsets, offset)
		id := i + 1
		header := fmt.Sprintf("%d 0 obj\n", id)
		w.Write([]byte(header))
		w.Write([]byte(obj))
		w.Write([]byte("\nendobj\n"))
		offset += len(header) + len(obj) + 8 // 8 for \nendobj\n
	}

	// Write Xref
	xrefOffset := offset
	w.Write([]byte("xref\n"))
	w.Write([]byte(fmt.Sprintf("0 %d\n", len(objects)+1)))
	w.Write([]byte("0000000000 65535 f \n"))
	for _, o := range offsets {
		w.Write([]byte(fmt.Sprintf("%010d 00000 n \n", o)))
	}

	// Write Trailer
	w.Write([]byte("trailer\n"))
	w.Write([]byte(fmt.Sprintf("<< /Size %d /Root 1 0 R >>\n", len(objects)+1)))
	w.Write([]byte("startxref\n"))
	w.Write([]byte(fmt.Sprintf("%d\n", xrefOffset)))
	w.Write([]byte("%%EOF\n"))

	return nil
}
