package dxfconverror

import "fmt"

// InternalError represents an internal error such as I/O failure during temp file creation.
type InternalError struct {
	Err error
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("internal error: %v", e.Err)
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

// ParseError represents an error during DXF parsing.
type ParseError struct {
	Err error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error: %v", e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// RenderingError represents an error during the rendering phase.
type RenderingError struct {
	Err error
}

func (e *RenderingError) Error() string {
	return fmt.Sprintf("rendering error: %v", e.Err)
}

func (e *RenderingError) Unwrap() error {
	return e.Err
}
