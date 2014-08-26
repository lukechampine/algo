package algo

import (
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
