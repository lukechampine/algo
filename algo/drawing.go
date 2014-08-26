package algo

import (
	"image"
	"image/color/palette"
	"math"
)

type Vector struct {
	X, Y, Z float64
}

func (v Vector) Sub(v2 Vector) Vector {
	return Vector{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector) Unit() Vector {
	l := v.Length()
	return Vector{v.X / l, v.Y / l, v.Z / l}
}

func (u Vector) Cross(v Vector) Vector {
	return Vector{
		u.Y*v.Z - u.Z*v.Y,
		u.Z*v.X - u.X*v.Z,
		u.X*v.Y - u.Y*v.X,
	}
}

type Matrix [3][3]float64

func (m *Matrix) Multiply(v *Vector) Vector {
	return Vector{
		v.X*m[0][0] + v.Y*m[0][1] + v.Z*m[0][2],
		v.X*m[1][0] + v.Y*m[1][1] + v.Z*m[1][2],
		v.X*m[2][0] + v.Y*m[2][1] + v.Z*m[2][2]}
}

func RotationMatrix(angle float64, v Vector) Matrix {
	s := math.Sin(angle * math.Pi / 180)
	c := math.Cos(angle * math.Pi / 180)
	v = v.Unit()
	return [3][3]float64{
		[3]float64{v.X*v.X*(1-c) + c, v.X*v.Y*(1-c) - v.Z*s, v.X*v.Y*(1-c) + v.Y*s},
		[3]float64{v.X*v.Y*(1-c) + v.Z*s, v.Y*v.Y*(1-c) + c, v.Y*v.Z*(1-c) - v.X*s},
		[3]float64{v.X*v.Z*(1-c) - v.Y*s, v.Y*v.Z*(1-c) + v.X*s, v.Z*v.Z*(1-c) + c}}
}

type Line struct {
	V1, V2 Vector
}

func (l *Line) Transform(m Matrix) Line {
	return Line{m.Multiply(&l.V1), m.Multiply(&l.V2)}
}

func (l Line) Projection() Line {
	return Line{Vector{l.V1.X, l.V1.Y, 0}, Vector{l.V2.X, l.V2.Y, 0}}
}

type Canvas struct {
	image.Paletted
}

func NewCanvas(width, height int) *Canvas {
	canvas := new(Canvas)
	canvas.Paletted = *image.NewPaletted(image.Rect(0, 0, width, height), palette.Plan9)
	canvas.Clear()
	return canvas
}

func (c *Canvas) Clear() {
	for i := range c.Pix {
		c.Pix[i] = 0xFF // white
	}
}

// CartSet adjusts its input x and y coordinates to fit a Cartesian plane.
// In other words, it places (0,0) at the center of the Canvas.
func (c *Canvas) CartSet(x, y int, colorIndex uint8) {
	c.SetColorIndex(c.Bounds().Size().X/2+x, c.Bounds().Size().Y/2-y, colorIndex)
}

func (c *Canvas) DrawLine(l Line) {
	length := l.V2.Sub(l.V1).Length()
	unit := l.V2.Sub(l.V1).Unit()
	for i := float64(0); i < length; i++ {
		x := l.V1.X + i*unit.X
		y := l.V1.Y + i*unit.Y
		c.CartSet(int(x), int(y), 0x00) // black
	}
}
