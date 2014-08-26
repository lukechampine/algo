package main

import (
	"image/color"

	. "github.com/lukechampine/algo/algo"
)

func main() {
	// canvas properties
	width, height := 700, 700

	// drawing variables
	// figure (vertices must be defined in clockwise order)
	figure := [12]Triangle{
		// front (Z = 100)
		Triangle{Vector{-100, 100, 100}, Vector{100, 100, 100}, Vector{-100, -100, 100}},
		Triangle{Vector{-100, -100, 100}, Vector{100, 100, 100}, Vector{100, -100, 100}},
		// top (Y = 100)
		Triangle{Vector{-100, 100, -100}, Vector{100, 100, -100}, Vector{-100, 100, 100}},
		Triangle{Vector{-100, 100, 100}, Vector{100, 100, -100}, Vector{100, 100, 100}},
		// right (X = 100)
		Triangle{Vector{100, 100, 100}, Vector{100, 100, -100}, Vector{100, -100, 100}},
		Triangle{Vector{100, -100, 100}, Vector{100, 100, -100}, Vector{100, -100, -100}},
		// back (Z = -100)
		Triangle{Vector{100, 100, -100}, Vector{-100, 100, -100}, Vector{100, -100, -100}},
		Triangle{Vector{100, -100, -100}, Vector{-100, 100, -100}, Vector{-100, -100, -100}},
		// bottom (Y = -100)
		Triangle{Vector{-100, -100, 100}, Vector{100, -100, 100}, Vector{-100, -100, -100}},
		Triangle{Vector{-100, -100, -100}, Vector{100, -100, 100}, Vector{100, -100, -100}},
		// left (X = -100)
		Triangle{Vector{-100, 100, -100}, Vector{-100, 100, 100}, Vector{-100, -100, -100}},
		Triangle{Vector{-100, -100, -100}, Vector{-100, 100, 100}, Vector{-100, -100, 100}},
	}
	numSteps := 100
	fw := NewFrameWriter(numSteps)

	// drawing function
	fw.GenerateFrames(func(num int) *Canvas {
		canvas := NewCanvas(width, height)
		tMatrix := RotationMatrix(float64(num*360/numSteps), Vector{1, 1, 1})
		for j := range figure {
			tTri := figure[j].Transform(tMatrix)
			if tTri.IsVisible() {
				canvas.FillTriangle(tTri, color.RGBA{uint8(j * 20), 0, uint8(255 - j*20), 255})
				// draw the outline last, since we want it on top
				canvas.DrawTriangle(tTri)
			}
		}
		return canvas
	})

	// encode frames as a GIF
	err := fw.WriteToFile("output.gif")
	if err != nil {
		println(err.Error())
	}
}
