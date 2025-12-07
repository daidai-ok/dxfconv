package dxf

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Tag represents a single DXF tag (group code and value).
type Tag struct {
	Code  int
	Value string
	Line  int
}

// Scanner scans a DXF stream for tags.
type Scanner struct {
	scanner       *bufio.Scanner
	Line          int
	NextTag       *Tag
	pushedBackTag *Tag
	Err           error
}

// NewScanner creates a new scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		scanner: bufio.NewScanner(r),
		Line:    0,
	}
}

// Scan advances to the next tag.
func (s *Scanner) Scan() bool {
	if s.Err != nil {
		return false
	}

	if s.pushedBackTag != nil {
		s.NextTag = s.pushedBackTag
		s.pushedBackTag = nil
		return true
	}

	// Read Code
	if !s.scanner.Scan() {
		s.Err = s.scanner.Err()
		return false
	}
	s.Line++
	codeLine := s.Line
	codeStr := strings.TrimSpace(s.scanner.Text())

	// Read Value
	if !s.scanner.Scan() {
		s.Err = s.scanner.Err()
		if s.Err == nil {
			s.Err = io.ErrUnexpectedEOF
		}
		return false
	}
	s.Line++
	// We only trim leading whitespace for values to preserve significant trailing spaces (e.g. in Text).
	// Newline characters are already stripped by Scanner.Scan().
	valStr := strings.TrimLeft(s.scanner.Text(), " \t")

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		s.Err = fmt.Errorf("line %d: invalid group code '%s': %w", codeLine, codeStr, err)
		return false
	}

	s.NextTag = &Tag{Code: code, Value: valStr, Line: codeLine}
	return true
}

func (s *Scanner) PushBack() {
	s.pushedBackTag = s.NextTag
}

// Int returns the value as an int.
func (t *Tag) Int() (int, error) {
	v, err := strconv.Atoi(t.Value)
	if err != nil {
		return 0, fmt.Errorf("line %d: invalid integer '%s': %w", t.Line+1, t.Value, err)
	}
	return v, nil
}

// Float returns the value as a float64.
func (t *Tag) Float() (float64, error) {
	v, err := strconv.ParseFloat(t.Value, 64)
	if err != nil {
		return 0, fmt.Errorf("line %d: invalid float '%s': %w", t.Line+1, t.Value, err)
	}
	return v, nil
}
