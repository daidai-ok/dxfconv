package converter

import (
	"fmt"
	"io"
	"os"

	"github.com/yofu/dxf"
	"github.com/yofu/dxf/entity"

	"dxfconv/pkg/renderers"
)

// Convert reads DXF data from r and writes PDF data to w
func Convert(r io.Reader, w io.Writer, opts *Options) error {
	if opts == nil {
		opts = DefaultOptions()
	}

	// Since github.com/yofu/dxf mainly works with files, we might need a workaround if it doesn't support io.Reader directly.
	// Checking the library... usually Open takes a filename.
	// If the library only supports file paths, we might need to change the API or write to a temp file.
	// TODO: In the future, we aim to remove the dependency on github.com/yofu/dxf to allow for more flexible inclusion of metadata such as error location in parse errors.
	tempFile, err := os.CreateTemp("", "dxf-*.dxf")
	if err != nil {
		return &InternalError{Err: fmt.Errorf("failed to create temp file: %w", err)}
	}
	defer os.Remove(tempFile.Name()) // clean up

	_, err = io.Copy(tempFile, r)
	if err != nil {
		return &InternalError{Err: fmt.Errorf("failed to write to temp file: %w", err)}
	}
	tempFile.Close()

	dxfDrawing, err := dxf.Open(tempFile.Name())
	if err != nil {
		return &ParseError{Err: fmt.Errorf("failed to parse DXF: %w", err)}
	}

	// Calculate Bounding Box
	bb := NewBoundingBox()
	for _, e := range dxfDrawing.Entities() {
		switch e := e.(type) {
		case *entity.Line:
			bb.Update(e.Start[0], e.Start[1])
			bb.Update(e.End[0], e.End[1])
		case *entity.Circle:
			bb.Update(e.Center[0]-e.Radius, e.Center[1]-e.Radius)
			bb.Update(e.Center[0]+e.Radius, e.Center[1]+e.Radius)
		case *entity.LwPolyline:
			for _, v := range e.Vertices {
				bb.Update(v[0], v[1])
			}
		case *entity.Polyline:
			for _, v := range e.Vertices {
				bb.Update(v.Coord[0], v.Coord[1])
			}
		case *entity.Spline:
			for _, v := range e.Controls {
				bb.Update(v[0], v[1])
			}
		case *entity.Point:
			bb.Update(e.Coord[0], e.Coord[1])
		case *entity.Text:
			bb.Update(e.Coord1[0], e.Coord1[1])
			// Text width/height is complicated to calculate without font metrics, so we only include the anchor point.
			// Or we could approximate:
			// bb.Update(e.Coord1[0], e.Coord1[1] + e.Height)
			// For now, let's Stick to the anchor point to avoid incorrect expansion.
		}
	}

	// Setup Renderer
	var renderer Renderer
	pageW, pageH := opts.PageSize.Width, opts.PageSize.Height
	if opts.Orientation == OrientationLandscape {
		pageW, pageH = pageH, pageW
	}

	switch opts.Format {
	case FormatSVG:
		renderer = renderers.NewSVGRenderer(w, pageW, pageH, opts.Font)
	case FormatPDF:
		fallthrough
	default:
		renderer = renderers.NewPDFRenderer(w, string(opts.Orientation), pageW, pageH, opts.Font)
	}

	renderer.Init(pageW, pageH)

	// Calculate Scale
	availW := pageW - (2 * opts.Margin)
	availH := pageH - (2 * opts.Margin)

	scale := opts.Scale
	if scale == 0 {
		scaleX := availW / bb.Width()
		scaleY := availH / bb.Height()
		if scaleX < scaleY {
			scale = scaleX
		} else {
			scale = scaleY
		}
	}

	// Centering:
	// Wait, in drawEntity we did: PDF_Y = Height - ((y * scale) + offsetY)
	// If we want the drawing to be centered, we need to be careful with the flip.

	// Let's adjust drawEntity logic to be simpler:
	// PDF_X = (DXF_X - bb.MinX) * scale + Margin + (availW - bb.Width()*scale)/2
	// PDF_Y = (pageH - Margin - (availH - bb.Height()*scale)/2) - (DXF_Y - bb.MinY) * scale

	// Let's recalculate the offsets to pass to drawEntity to make the previous formula work.
	// Previous formula:
	// X = x * scale + offsetX
	// Y = height - (y * scale + offsetY)

	// We want:
	// X = (x - bb.MinX) * scale + Margin + CenteringX
	// X = x * scale + (-bb.MinX * scale + Margin + CenteringX)
	// So passed OffsetX = -bb.MinX * scale + Margin + (availW - bb.Width()*scale)/2

	// We want:
	// Y = PageH - (Margin + CenteringY + (y - bb.MinY) * scale)
	// Y = PageH - (Margin + CenteringY - bb.MinY * scale + y * scale)
	// Y = PageH - ( (Margin + CenteringY - bb.MinY * scale) + y * scale )
	// So passed OffsetY = Margin + (availH - bb.Height()*scale)/2 - bb.MinY * scale

	realOffsetX := -bb.MinX*scale + opts.Margin + (availW-bb.Width()*scale)/2
	realOffsetY := -bb.MinY*scale + opts.Margin + (availH-bb.Height()*scale)/2

	// Draw Entities
	for _, e := range dxfDrawing.Entities() {
		drawEntity(renderer, e, scale, realOffsetX, realOffsetY, pageH)
	}

	if err := renderer.Finish(); err != nil {
		return &RenderingError{Err: err}
	}
	return nil
}
