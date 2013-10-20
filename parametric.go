package main

import (
	"log"
	"math"
	"os/exec"
	"runtime"
)

func (c Canvas) DrawParaFn(step float64) {
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
	frameChan := make(chan int, numSteps)
	for i := 0; i < numSteps; i++ {
		frameChan <- i
		// draw each frame in a separate goroutine
		go func() {
			num := <-frameChan
			canvas := NewCanvas(width, height)
			canvas.DrawParaFn(step + float64(num)*inc)
			canvas.SaveToFile(num)
		}()
	}
	log.Print("finished computation")
	exec.Command("sh", "-c", "convert -delay 10 -loop 0 input-*.png output.gif").Run()
	log.Print("output rendered as gif")
	exec.Command("sh", "-c", "rm input-*.png").Run()
}
