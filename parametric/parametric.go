package main

import (
	"math"
	"runtime"
	"sync"

	. "github.com/lukechampine/algo/algo"
)

func drawParaFn(c *Canvas, step float64) {
	t := float64(0)
	for i := 0; i < 100; i++ {
		c.DrawLine(Line{paraFn(t), paraFn(t + step)})
		t += step
	}
}

func paraFn(t float64) Vector {
	x := 340 * math.Exp(-0.1*t) * math.Cos(t)
	y := 340 * math.Exp(-0.1*t) * math.Sin(t)
	return Vector{x, y, 0}
}

func main() {
	// use all available logical processors
	runtime.GOMAXPROCS(runtime.NumCPU())
	// canvas properties
	width, height := 700, 700
	// draw loop
	step := 0.0
	numSteps := 100
	inc := math.Pi / float64(numSteps)
	fw := NewFrameWriter(numSteps)
	var wg sync.WaitGroup
	wg.Add(numSteps)
	for i := 0; i < numSteps; i++ {
		// draw each frame in a separate goroutine
		go func(num int) {
			canvas := NewCanvas(width, height)
			drawParaFn(canvas, step+float64(num)*inc)
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
