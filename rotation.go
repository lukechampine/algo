package main

import (
	"log"
	"math"
	"os/exec"
	"runtime"
)

type Matrix [3][3]float64

func (m Matrix) Multiply(v Vector) Vector {
	return Vector{
		v.X*m[0][0] + v.Y*m[0][1] + v.Z*m[0][2],
		v.X*m[1][0] + v.Y*m[1][1] + v.Z*m[1][2],
		v.X*m[2][0] + v.Y*m[2][1] + v.Z*m[2][2]}
}

func (l Line) Transform(m Matrix) Line {
	return Line{m.Multiply(l.V1), m.Multiply(l.V2)}
}

// project a 3D vector onto 2D space so we can draw it
func (l Line) Projection() Line {
	return Line{Vector{l.V1.X, l.V1.Y, 0}, Vector{l.V2.X, l.V2.Y, 0}}
}

// generate a transformation matrix that rotates along vector v by angle degrees
func RotationMatrix(angle float64, v Vector) Matrix {
	s := math.Sin(angle * math.Pi / 180)
	c := math.Cos(angle * math.Pi / 180)
	v = v.Unit()
	return [3][3]float64{
		[3]float64{v.X*v.X*(1-c) + c, v.X*v.Y*(1-c) - v.Z*s, v.X*v.Y*(1-c) + v.Y*s},
		[3]float64{v.X*v.Y*(1-c) + v.Z*s, v.Y*v.Y*(1-c) + c, v.Y*v.Z*(1-c) - v.X*s},
		[3]float64{v.X*v.Z*(1-c) - v.Y*s, v.Y*v.Z*(1-c) + v.X*s, v.Z*v.Z*(1-c) + c}}
}

func main() {
	// use all available logical processors
	runtime.GOMAXPROCS(runtime.NumCPU())
	// canvas properties
	width, height := 700, 700
	// figure
	lines := [12]Line{
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
	frameChan := make(chan int, numSteps)
	for i := 0; i < numSteps; i++ {
		frameChan <- i
		// draw each frame in a separate goroutine
		go func() {
			num := <-frameChan
			canvas := NewCanvas(width, height)
			tMatrix := RotationMatrix(float64(num*360/numSteps), Vector{1, 1, 1})
			for j := range lines {
				canvas.DrawLine(lines[j].Transform(tMatrix).Projection())
			}
			canvas.SaveToFile(num)
		}()
	}
	// animation and clean-up
	log.Print("finished computation")
	exec.Command("sh", "-c", "convert -delay 10 -loop 0 input-*.png output.gif").Run()
	log.Print("output rendered as gif")
	exec.Command("sh", "-c", "rm input-*.png").Run()
}
