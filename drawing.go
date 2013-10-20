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

type Line struct {
	V1, V2 Vector
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

func (c Canvas) DrawLine(l Line) {
	length := l.V2.Sub(l.V1).Length()
	unit := l.V2.Sub(l.V1).Unit()
	for i := 0; i < int(length); i++ {
		x := l.V1.X + float64(i)*unit.X
		y := l.V1.Y + float64(i)*unit.Y
		// adjust for canvas coordinate system
		c.Set(int(x)+c.Bounds().Size().X/2, int(-y)+c.Bounds().Size().Y/2, color.RGBA{0, 0, 0, 255})
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
