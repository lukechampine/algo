package main

import (
	. "github.com/lukechampine/algo/algo"
)

func main() {
	// canvas properties
	width, height := 700, 700

	// drawing variables
	figure := [12]Line{
		Line{Vector{-100, 100, 100}, Vector{100, 100, 100}},
		Line{Vector{100, 100, 100}, Vector{100, -100, 100}},
		Line{Vector{100, -100, 100}, Vector{-100, -100, 100}},
		Line{Vector{-100, -100, 100}, Vector{-100, 100, 100}},
		Line{Vector{-100, 100, -100}, Vector{100, 100, -100}},
		Line{Vector{100, 100, -100}, Vector{100, -100, -100}},
		Line{Vector{100, -100, -100}, Vector{-100, -100, -100}},
		Line{Vector{-100, -100, -100}, Vector{-100, 100, -100}},
		Line{Vector{-100, 100, 100}, Vector{-100, 100, -100}},
		Line{Vector{100, 100, 100}, Vector{100, 100, -100}},
		Line{Vector{100, -100, 100}, Vector{100, -100, -100}},
		Line{Vector{-100, -100, 100}, Vector{-100, -100, -100}},
	}
	numSteps := 100
	fw := NewFrameWriter(numSteps)

	// drawing function
	fw.GenerateFrames(func(num int) *Canvas {
		canvas := NewCanvas(width, height)
		tMatrix := RotationMatrix(float64(num*360/numSteps), Vector{1, 1, 1})
		for j := range figure {
			canvas.DrawLine(figure[j].Transform(tMatrix).Projection())
		}
		return canvas
	})

	// encode frames as a gif
	err := fw.WriteToFile("output.gif")
	if err != nil {
		println(err.Error())
	}
	return
}
