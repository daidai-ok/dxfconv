package dxf

import (
	"io"
)

// Reference: https://help.autodesk.com/view/OARX/2021/ENU/?guid=GUID-235B22E0-A567-4CF6-92D3-38A2306D73F3

// Parse reads a DXF file and returns a Drawing.
func Parse(r io.Reader) (*Drawing, error) {
	scanner := NewScanner(r)
	drawing := &Drawing{}

	for scanner.Scan() {
		tag := scanner.NextTag
		if tag.Code == 0 && tag.Value == "SECTION" {
			if err := parseSection(scanner, drawing); err != nil {
				return nil, err
			}
		} else if tag.Code == 0 && tag.Value == "EOF" {
			break
		}
	}

	if scanner.Err != nil {
		return nil, scanner.Err
	}

	return drawing, nil
}

func parseSection(s *Scanner, d *Drawing) error {
	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 2 {
			if tag.Value == "ENTITIES" {
				return parseEntities(s, d)
			}
			// Skip other sections
			return skipSection(s)
		}
		if tag.Code == 0 && tag.Value == "ENDSEC" {
			return nil
		}
	}
	return nil
}

func skipSection(s *Scanner) error {
	for s.Scan() {
		if s.NextTag.Code == 0 && s.NextTag.Value == "ENDSEC" {
			return nil
		}
	}
	return s.Err
}

func parseEntities(s *Scanner, d *Drawing) error {
	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 0 {
			if tag.Value == "ENDSEC" {
				return nil
			}
			entity, err := parseEntity(tag.Value, s)
			if err != nil {
				return err
			}
			if entity != nil {
				d.Entities = append(d.Entities, entity)
			}
		}
	}
	return s.Err
}

func parseEntity(typeStr string, s *Scanner) (Entity, error) {
	var e Entity
	var err error

	switch typeStr {
	case "LINE":
		e, err = parseLine(s)
	case "CIRCLE":
		e, err = parseCircle(s)
	case "ARC":
		e, err = parseArc(s)
	case "LWPOLYLINE":
		e, err = parseLwPolyline(s)
	case "POLYLINE":
		e, err = parsePolyline(s)
	case "SPLINE":
		e, err = parseSpline(s)
	case "POINT":
		e, err = parsePoint(s)
	case "TEXT":
		e, err = parseText(s)
	case "MTEXT":
		e, err = parseMText(s)
	default:
		return nil, skipEntity(s)
	}

	return e, err
}

func skipEntity(s *Scanner) error {
	for s.Scan() {
		if s.NextTag.Code == 0 {
			s.PushBack()
			return nil
		}
	}
	return s.Err
}

func parseCommon(s *Scanner, e *BaseEntity) bool {
	tag := s.NextTag
	if tag.Code == 8 {
		e.LayerName = tag.Value
		return true
	}
	return false
}

func parseLine(s *Scanner) (*Line, error) {
	l := &Line{BaseEntity: BaseEntity{EntityType: LineType}}
	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 0 {
			s.PushBack()
			return l, nil
		}
		if parseCommon(s, &l.BaseEntity) {
			continue
		}
		val, err := tag.Float()
		if err != nil {
			return nil, err
		}
		switch tag.Code {
		case 10:
			l.Start[0] = val
		case 20:
			l.Start[1] = val
		case 11:
			l.End[0] = val
		case 21:
			l.End[1] = val
		}
	}
	return l, s.Err
}

func parseCircle(s *Scanner) (*Circle, error) {
	c := &Circle{BaseEntity: BaseEntity{EntityType: CircleType}}
	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 0 {
			s.PushBack()
			return c, nil
		}
		if parseCommon(s, &c.BaseEntity) {
			continue
		}
		val, err := tag.Float()
		if err != nil {
			return nil, err
		}
		switch tag.Code {
		case 10:
			c.Center[0] = val
		case 20:
			c.Center[1] = val
		case 40:
			c.Radius = val
		}
	}
	return c, s.Err
}

func parseArc(s *Scanner) (*Arc, error) {
	a := &Arc{BaseEntity: BaseEntity{EntityType: ArcType}}
	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 0 {
			s.PushBack()
			return a, nil
		}
		if parseCommon(s, &a.BaseEntity) {
			continue
		}
		val, err := tag.Float()
		if err != nil {
			return nil, err
		}
		switch tag.Code {
		case 10:
			a.Center[0] = val
		case 20:
			a.Center[1] = val
		case 40:
			a.Radius = val
		case 50:
			a.StartAngle = val
		case 51:
			a.EndAngle = val
		}
	}
	return a, s.Err
}

func parseLwPolyline(s *Scanner) (*LwPolyline, error) {
	l := &LwPolyline{BaseEntity: BaseEntity{EntityType: LwPolylineType}}
	var currentVertex *LwPolylineVertex

	// Helper to commit current vertex
	commitVertex := func() {
		if currentVertex != nil {
			l.Vertices = append(l.Vertices, *currentVertex)
			currentVertex = nil
		}
	}

	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 0 {
			commitVertex()
			s.PushBack()
			return l, nil
		}
		if parseCommon(s, &l.BaseEntity) {
			continue
		}

		if tag.Code == 10 { // Start of a new vertex (X)
			commitVertex()
			val, err := tag.Float()
			if err != nil {
				return nil, err
			}
			currentVertex = &LwPolylineVertex{X: val}
			continue
		}

		if tag.Code == 70 {
			val, _ := tag.Int()
			if val&1 == 1 {
				l.Closed = true
			}
			continue
		}

		if currentVertex != nil {
			val, err := tag.Float()
			if err != nil {
				return nil, err
			}
			switch tag.Code {
			case 20:
				currentVertex.Y = val
			case 42:
				currentVertex.Bulge = val
			}
		}
	}
	commitVertex()
	return l, s.Err
}

func parsePolyline(s *Scanner) (*Polyline, error) {
	p := &Polyline{BaseEntity: BaseEntity{EntityType: PolylineType}}
	// FLAGS: 70
	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 0 {
			s.PushBack()
			break
		}
		if parseCommon(s, &p.BaseEntity) {
			continue
		}
		if tag.Code == 70 {
			val, _ := tag.Int()
			if val&1 == 1 {
				p.Closed = true
			}
			// check for 3D or Mesh flags if needed
		}
	}

	// Now consume VERTEX entities until SEQEND
	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 0 {
			if tag.Value == "SEQEND" {
				// Consume SEQEND until its end (usually just 0 SEQEND then 0 NEXT)
				// Actually SEQEND is an entity, needs to be skipped or parsed.
				// We just need to consume everything until next 0.
				_ = skipEntity(s)
				return p, nil
			} else if tag.Value == "VERTEX" {
				v, err := parseVertex(s)
				if err != nil {
					return nil, err
				}
				p.Vertices = append(p.Vertices, *v)
			} else {
				// Unexpected entity inside POLYLINE sequence, possibly we misread?
				// Just push back and return what we have?
				// The spec says SEQEND is required.
				s.PushBack()
				return p, nil
			}
		}
	}
	return p, s.Err
}

func parseVertex(s *Scanner) (*Vertex, error) {
	v := &Vertex{}
	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 0 {
			s.PushBack()
			return v, nil
		}
		val, err := tag.Float()
		if err != nil {
			continue // ignore non-float tags in vertex like layer
		}
		switch tag.Code {
		case 10:
			v.X = val
		case 20:
			v.Y = val
		case 30:
			v.Z = val
		}
	}
	return v, s.Err
}

func parseSpline(s *Scanner) (*Spline, error) {
	sp := &Spline{BaseEntity: BaseEntity{EntityType: SplineType}}
	var currentControl *[3]float64

	commitControl := func() {
		if currentControl != nil {
			sp.ControlPoints = append(sp.ControlPoints, *currentControl)
			currentControl = nil
		}
	}

	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 0 {
			commitControl()
			s.PushBack()
			return sp, nil
		}
		if parseCommon(s, &sp.BaseEntity) {
			continue
		}

		// 10 is start of control point
		if tag.Code == 10 {
			commitControl()
			val, err := tag.Float()
			if err != nil {
				return nil, err
			}
			currentControl = &[3]float64{val, 0, 0}
			continue
		}

		val, err := tag.Float()
		if err == nil {
			if currentControl != nil {
				switch tag.Code {
				case 20:
					currentControl[1] = val
				case 30:
					currentControl[2] = val
				}
			}
			if tag.Code == 40 {
				sp.Knots = append(sp.Knots, val)
			}
		}

		if tag.Code == 71 {
			degree, _ := tag.Int()
			sp.Degree = degree
		}
		if tag.Code == 70 {
			flag, _ := tag.Int()
			if flag&1 == 1 {
				sp.Closed = true
			}
		}
	}
	commitControl()
	return sp, s.Err
}

func parsePoint(s *Scanner) (*Point, error) {
	p := &Point{BaseEntity: BaseEntity{EntityType: PointType}}
	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 0 {
			s.PushBack()
			return p, nil
		}
		if parseCommon(s, &p.BaseEntity) {
			continue
		}
		val, err := tag.Float()
		if err != nil {
			return nil, err
		}
		switch tag.Code {
		case 10:
			p.Coord[0] = val
		case 20:
			p.Coord[1] = val
		case 30:
			p.Coord[2] = val
		}
	}
	return p, s.Err
}

func parseText(s *Scanner) (*Text, error) {
	t := &Text{BaseEntity: BaseEntity{EntityType: TextType}}
	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 0 {
			s.PushBack()
			return t, nil
		}
		if parseCommon(s, &t.BaseEntity) {
			continue
		}

		if tag.Code == 1 {
			t.Value = tag.Value
			continue
		}

		val, err := tag.Float()
		if err != nil {
			continue
		}
		switch tag.Code {
		case 10:
			t.Point[0] = val
		case 20:
			t.Point[1] = val
		case 30:
			t.Point[2] = val
		case 40:
			t.Height = val
		}
	}
	return t, s.Err
}

func parseMText(s *Scanner) (*MText, error) {
	t := &MText{BaseEntity: BaseEntity{EntityType: MTextType}}
	var textBuf string // MText can be split across multiple code 1/3 tags

	for s.Scan() {
		tag := s.NextTag
		if tag.Code == 0 {
			s.PushBack()
			t.Value = textBuf
			return t, nil
		}
		if parseCommon(s, &t.BaseEntity) {
			continue
		}

		if tag.Code == 1 || tag.Code == 3 {
			textBuf += tag.Value
			continue
		}

		val, err := tag.Float()
		if err != nil {
			continue
		}
		switch tag.Code {
		case 10:
			t.Point[0] = val
		case 20:
			t.Point[1] = val
		case 30:
			t.Point[2] = val
		case 40:
			t.Height = val
		}
	}
	t.Value = textBuf
	return t, s.Err
}
