package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"
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

type Matrix [3][3]float64

func (m Matrix) Multiply(v Vector) Vector {
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

func (l Line) Transform(m Matrix) Line {
	return Line{m.Multiply(l.V1), m.Multiply(l.V2)}
}

func (l Line) Projection() Line {
	return Line{Vector{l.V1.X, l.V1.Y, 0}, Vector{l.V2.X, l.V2.Y, 0}}
}

type Canvas struct {
	image.RGBA
}

func NewCanvas(width, height int) *Canvas {
	canvas := new(Canvas)
	canvas.RGBA = *image.NewRGBA(image.Rect(0, 0, width, height))
	size := canvas.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			color := color.RGBA{255, 255, 255, 255}
			canvas.Set(x, y, color)
		}
	}
	return canvas
}

func (c Canvas) SmartSet(x, y int, color color.RGBA) {
	c.Set(x+c.Bounds().Size().X/2, -y+c.Bounds().Size().Y/2, color)
}

func (c Canvas) DrawLine(l Line) {
	length := l.V2.Sub(l.V1).Length()
	unit := l.V2.Sub(l.V1).Unit()
	for i := 0; i < int(length); i++ {
		x := l.V1.X + float64(i)*unit.X
		y := l.V1.Y + float64(i)*unit.Y
		// adjust for canvas coordinate system
		c.SmartSet(int(x), int(y), color.RGBA{0, 0, 0, 255})
	}
}

func (c Canvas) SaveToFile(num int) {
	outFilename := ""
	if num < 10 {
		outFilename = "input-00" + strconv.Itoa(num) + ".png"
	} else if num < 100 {
		outFilename = "input-0" + strconv.Itoa(num) + ".png"
	} else {
		outFilename = "input-" + strconv.Itoa(num) + ".png"
	}
	outFile, _ := os.Create(outFilename)
	defer outFile.Close()
	png.Encode(outFile, &c)
}
