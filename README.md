<p align="center">
  <img alt="dxfconv logo" src="docs/icon.png" height="140" />
  <h3 align="center">dxfconv</h3>
  <p align="center">A lightweight CAD DXF converter</p>
</p>

---

`dxfconv` is a lightweight Go library for converting DXF (Drawing Exchange Format) files into PDF or SVG. It is designed to be simple to use and easy to integrate into your Go applications.

<p align="center">
  <img alt="summary" src="docs/summary.png" height="380" />
</p>

## Features

-   **DXF to PDF**: Convert CAD drawings to standard PDF documents.
-   **DXF to SVG**: Convert CAD drawings to Scalable Vector Graphics.
-   **Flexible Output**: Write to files or directly to `io.Writer` (e.g., `bytes.Buffer`, HTTP response).
-   **Entity Support**: Supports common DXF entities:
    -   LINES
    -   CIRCLES
    -   ARCS
    -   LWPOLYLINES
    -   POLYLINES
    -   SPLINES
    -   TEXT
-   **Customization**: Control page size (A4, A3, etc.), orientation (Portrait, Landscape), and scaling.

## Installation

```bash
go get github.com/daidai-ok/dxfconv
```

## Usage

### Basic Conversion (DXF to PDF)

```go
package main

import (
	"os"
	dxfconv "dxfconv/pkg/converter"
)

func main() {
	f, _ := os.Open("drawing.dxf")
	defer f.Close()

	out, _ := os.Create("drawing.pdf")
	defer out.Close()

	// Convert with default options (PDF, A4, Portrait)
	dxfconv.Convert(f, out, nil)
}
```

### DXF to SVG

```go
opts := dxfconv.DefaultOptions()
opts.Format = dxfconv.FormatSVG

dxfconv.Convert(f, out, opts)
```

### Writing to Buffer

```go
var buf bytes.Buffer
dxfconv.Convert(f, &buf, nil)
// Use buf.Bytes() ...
```

### Custom Options

```go
opts := &dxfconv.Options{
	PageSize:    dxfconv.PageSizeA3,
	Orientation: dxfconv.OrientationLandscape,
	Format:      dxfconv.FormatPDF,
	Scale:       1.0, // 1:1 scale
	Margin:      20.0,
}
dxfconv.Convert(f, out, opts)
```

## License

MIT License

The [Go gopher](https://go.dev/blog/gopher) was designed by [Renée French](https://reneefrench.blogspot.com/).
The project icon is based on [gophers](https://github.com/egonelbre/gophers) by [Egon Elbre](https://github.com/egonelbre), licensed under [CC0 1.0 Universal](https://creativecommons.org/publicdomain/zero/1.0/).
