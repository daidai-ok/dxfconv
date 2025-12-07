package converter

// PageSize represents the dimensions of the PDF page
type PageSize struct {
	Width  float64
	Height float64
}

// Common page sizes (in mm)
var (
	PageSizeA4 = PageSize{Width: 210, Height: 297}
	PageSizeA3 = PageSize{Width: 297, Height: 420}
)

// Orientation represents the page orientation
type Orientation string

const (
	OrientationPortrait  Orientation = "P"
	OrientationLandscape Orientation = "L"
)

// Format represents the output format
type Format string

const (
	FormatPDF Format = "pdf"
	FormatSVG Format = "svg"
)

// Options configuration for the conversion
type Options struct {
	PageSize    PageSize
	Orientation Orientation
	// Format specifies the output format (pdf or svg)
	Format Format
	// Scale allows manual scaling. If 0, auto-scaling is used.
	Scale float64
	// Margin in mm
	Margin float64
}

// DefaultOptions returns the default configuration
func DefaultOptions() *Options {
	return &Options{
		PageSize:    PageSizeA4,
		Orientation: OrientationPortrait,
		Format:      FormatPDF,
		Scale:       0.0,
		Margin:      10.0,
	}
}
