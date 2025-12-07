package converter

import (
	"fmt"
	"io"

	"github.com/daidai-ok/dxfconv/pkg/boundingbox"
	"github.com/daidai-ok/dxfconv/pkg/dxf"
	"github.com/daidai-ok/dxfconv/pkg/dxfconverror"
	"github.com/daidai-ok/dxfconv/pkg/renderers"
)

// Convert reads DXF data from r and writes PDF data to w
func Convert(r io.Reader, w io.Writer, opts *Options) error {
	if opts == nil {
		opts = DefaultOptions()
	}

	dxfDrawing, err := dxf.Parse(r)
	if err != nil {
		return &dxfconverror.ParseError{Err: fmt.Errorf("failed to parse DXF: %w", err)}
	}

	// Calculate Bounding Box
	bb := calculateBoundingBox(dxfDrawing)

	// Setup Renderer
	var renderer renderers.Renderer
	pageW, pageH := opts.PageSize.Width, opts.PageSize.Height
	if opts.Orientation == OrientationLandscape {
		pageW, pageH = pageH, pageW
	}

	switch opts.Format {
	case FormatSVG:
		renderer = renderers.NewSVGRenderer(w, pageW, pageH)
	case FormatPDF:
		fallthrough
	default:
		renderer = renderers.NewPDFRenderer(w, string(opts.Orientation), pageW, pageH)
	}

	renderer.Init(pageW, pageH)

	// Calculate Scale
	availW := pageW - (2 * opts.Margin)
	availH := pageH - (2 * opts.Margin)

	scale := opts.Scale
	if scale == 0 {
		var scaleX, scaleY float64
		if bb.Width() > 0 {
			scaleX = availW / bb.Width()
		} else {
			scaleX = 1.0 // Default scale if width is 0 (e.g. single point or vertical line)
		}
		if bb.Height() > 0 {
			scaleY = availH / bb.Height()
		} else {
			scaleY = 1.0 // Default scale if height is 0 (e.g. single point or horizontal line)
		}

		if scaleX < scaleY {
			scale = scaleX
		} else {
			scale = scaleY
		}
	}

	// Centering Logic
	realOffsetX := -bb.MinX*scale + opts.Margin + (availW-bb.Width()*scale)/2
	realOffsetY := -bb.MinY*scale + opts.Margin + (availH-bb.Height()*scale)/2

	// Draw Entities
	for _, e := range dxfDrawing.Entities {
		renderers.DrawEntity(renderer, e, scale, realOffsetX, realOffsetY, pageH)
	}

	if err := renderer.Finish(); err != nil {
		return &dxfconverror.RenderingError{Err: err}
	}
	return nil
}

func calculateBoundingBox(dxfDrawing *dxf.Drawing) *boundingbox.BoundingBox {
	bb := boundingbox.NewBoundingBox()
	for _, e := range dxfDrawing.Entities {
		switch e := e.(type) {
		case *dxf.Line:
			bb.Update(e.Start[0], e.Start[1])
			bb.Update(e.End[0], e.End[1])
		case *dxf.Circle:
			bb.Update(e.Center[0]-e.Radius, e.Center[1]-e.Radius)
			bb.Update(e.Center[0]+e.Radius, e.Center[1]+e.Radius)
		case *dxf.Arc:
			// Arc bounding box is tricky, approximate with full circle for now or centers/endpoints
			// Better to just update center +/- radius
			bb.Update(e.Center[0]-e.Radius, e.Center[1]-e.Radius)
			bb.Update(e.Center[0]+e.Radius, e.Center[1]+e.Radius)
		case *dxf.LwPolyline:
			for _, v := range e.Vertices {
				bb.Update(v.X, v.Y)
			}
		case *dxf.Polyline:
			for _, v := range e.Vertices {
				bb.Update(v.X, v.Y)
			}
		case *dxf.Spline:
			for _, v := range e.ControlPoints {
				bb.Update(v[0], v[1])
			}
		case *dxf.Point:
			bb.Update(e.Coord[0], e.Coord[1])
		case *dxf.Text:
			bb.Update(e.Point[0], e.Point[1])
		case *dxf.MText:
			bb.Update(e.Point[0], e.Point[1])
		}
	}
	return bb
}
