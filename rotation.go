package main

import (
	"log"
	"os/exec"
	"runtime"
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
		Line{Vector{-100, -100, 100}, Vector{-100, -100, -100}}}
	// draw loop
	numSteps := 100
	for i := 0; i < numSteps; i++ {
		// draw each frame in a separate goroutine
		go func(num int) {
			canvas := NewCanvas(width, height)
			tMatrix := RotationMatrix(float64(num*360/numSteps), Vector{1, 1, 1})
			for j := range figure {
				canvas.DrawLine(figure[j].Transform(tMatrix).Projection())
			}
			canvas.SaveToFile(num)
		}(i)
	}
	// animation and clean-up
	log.Print("finished computation")
	exec.Command("sh", "-c", "convert -delay 10 -loop 0 input-*.png output.gif").Run()
	log.Print("output rendered as gif")
	exec.Command("sh", "-c", "rm input-*.png").Run()
}
