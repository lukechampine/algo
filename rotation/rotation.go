package main

import (
	"runtime"
	"sync"

	. "github.com/lukechampine/algo/algo"
)

func main() {
	// use all available logical processors
	runtime.GOMAXPROCS(runtime.NumCPU())
	// canvas properties
	width, height := 700, 700
	// figure
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
	// draw loop
	numSteps := 100
	fw := NewFrameWriter(numSteps)
	var wg sync.WaitGroup
	wg.Add(numSteps)
	for i := 0; i < numSteps; i++ {
		// draw each frame in a separate goroutine
		go func(num int) {
			canvas := NewCanvas(width, height)
			tMatrix := RotationMatrix(float64(num*360/numSteps), Vector{1, 1, 1})
			for j := range figure {
				canvas.DrawLine(figure[j].Transform(tMatrix).Projection())
			}
			fw.AddFrame(canvas, num)
			wg.Done()
		}(i)
	}
	wg.Wait()

	// encode frames as a gif
	err := fw.WriteToFile("output.gif")
	if err != nil {
		println(err.Error())
	}
	return
}
