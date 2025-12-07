package dxf

// EntityType represents the type of a DXF entity.
type EntityType string

const (
	LineType       EntityType = "LINE"
	CircleType     EntityType = "CIRCLE"
	ArcType        EntityType = "ARC"
	LwPolylineType EntityType = "LWPOLYLINE"
	PolylineType   EntityType = "POLYLINE"
	SplineType     EntityType = "SPLINE"
	PointType      EntityType = "POINT"
	TextType       EntityType = "TEXT"
	MTextType      EntityType = "MTEXT"
)

// Entity is the interface that all DXF entities implement.
type Entity interface {
	Type() EntityType
	Layer() string
}

// BaseEntity contains common properties for all entities.
type BaseEntity struct {
	EntityType EntityType
	LayerName  string
}

func (e *BaseEntity) Type() EntityType {
	return e.EntityType
}

func (e *BaseEntity) Layer() string {
	return e.LayerName
}

// Line represents a LINE entity.
type Line struct {
	BaseEntity
	Start [3]float64
	End   [3]float64
}

// Circle represents a CIRCLE entity.
type Circle struct {
	BaseEntity
	Center [3]float64
	Radius float64
}

// Arc represents an ARC entity.
type Arc struct {
	BaseEntity
	Center     [3]float64
	Radius     float64
	StartAngle float64
	EndAngle   float64
}

// LwPolyline represents a LWPOLYLINE entity.
type LwPolyline struct {
	BaseEntity
	Vertices []LwPolylineVertex
	Closed   bool
}

type LwPolylineVertex struct {
	X, Y, Z float64
	Bulge   float64
}

// Polyline represents a POLYLINE entity.
// Note: This often comes with VERTEX entities following it, terminated by SEQEND.
type Polyline struct {
	BaseEntity
	Vertices []Vertex
	Closed   bool
	Is3D     bool
	IsMesh   bool
}

type Vertex struct {
	X, Y, Z float64
}

// Spline represents a SPLINE entity.
type Spline struct {
	BaseEntity
	ControlPoints [][3]float64
	Knots         []float64
	Closed        bool
	Degree        int
}

// Point represents a POINT entity.
type Point struct {
	BaseEntity
	Coord [3]float64
}

// Text represents a TEXT entity.
type Text struct {
	BaseEntity
	Point  [3]float64
	Height float64
	Value  string
	// Rotation, Oblique, etc. can be added later
}

// MText represents an MTEXT entity.
type MText struct {
	BaseEntity
	Point  [3]float64
	Height float64
	Value  string
}
