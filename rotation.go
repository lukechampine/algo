package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

type Vector struct {
	X, Y, Z float64
}

func (v *Vector) Sub(v2 Vector) Vector {
	return Vector{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v *Vector) Unit() Vector {
	l := v.Length()
	return Vector{v.X / l, v.Y / l, v.Z / l}
}

type Matrix [3][3]float64

func (m *Matrix) Multiply(v Vector) Vector {
	return Vector{v.X*m[0][0] + v.Y*m[0][1] + v.Z*m[0][2],
		/*	   */ v.X*m[1][0] + v.Y*m[1][1] + v.Z*m[1][2],
		/*	   */ v.X*m[2][0] + v.Y*m[2][1] + v.Z*m[2][2]}
}

// generate a transformation matrix that rotates along vector v by angle degrees
func RotationMatrix(angle float64, v Vector) Matrix {
	s := math.Sin(angle * math.Pi / 180)
	c := math.Cos(angle * math.Pi / 180)
	v = v.Unit()
	return [3][3]float64{[3]float64{v.X*v.X*(1-c) + c, v.X*v.Y*(1-c) - v.Z*s, v.X*v.Y*(1-c) + v.Y*s},
		/*            */ [3]float64{v.X*v.Y*(1-c) + v.Z*s, v.Y*v.Y*(1-c) + c, v.Y*v.Z*(1-c) - v.X*s},
		/*            */ [3]float64{v.X*v.Z*(1-c) - v.Y*s, v.Y*v.Z*(1-c) + v.X*s, v.Z*v.Z*(1-c) + c}}
}

type Line struct {
	V1, V2 Vector
}

func (l *Line) Transform(m Matrix) Line {
	return Line{m.Multiply(l.V1), m.Multiply(l.V2)}
}

// project 3D vector onto 2D space so we can draw it
func (l Line) Projection() Line {
	return Line{Vector{l.V1.X, l.V1.Y, 0}, Vector{l.V2.X, l.V2.Y, 0}}
}

type Canvas struct {
	image.RGBA
}

func NewCanvas(r image.Rectangle) *Canvas {
	canvas := new(Canvas)
	canvas.RGBA = *image.NewRGBA(r)
	return canvas
}

func (c Canvas) Clear() {
	size := c.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			color := color.RGBA{255, 255, 255, 255}
			c.Set(x, y, color)
		}
	}
}

func (c Canvas) DrawLine(color color.RGBA, l Line) {
	delta := l.V2.Sub(l.V1)
	length := delta.Length()
	x_step, y_step := delta.X/length, delta.Y/length
	limit := int(length + 0.5)
	for i := 0; i < limit; i++ {
		// adjust for canvas coordinate system
		x := l.V1.X + float64(i)*x_step + float64(c.Bounds().Size().X/2)
		y := -(l.V1.Y + float64(i)*y_step) + float64(c.Bounds().Size().Y/2)
		c.Set(int(x), int(y), color)
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
	outFile, err := os.Create(outFilename)
	defer outFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(outFile, &c)
}

func main() {
	// use all available logical processors
	runtime.GOMAXPROCS(runtime.NumCPU())
	// canvas properties
	width, height := 700, 700
	// figure
	lines := make([]Line, 12)
	lines[0] = Line{Vector{-100, 100, 100}, Vector{100, 100, 100}}
	lines[1] = Line{Vector{100, 100, 100}, Vector{100, -100, 100}}
	lines[2] = Line{Vector{100, -100, 100}, Vector{-100, -100, 100}}
	lines[3] = Line{Vector{-100, -100, 100}, Vector{-100, 100, 100}}
	lines[4] = Line{Vector{-100, 100, -100}, Vector{100, 100, -100}}
	lines[5] = Line{Vector{100, 100, -100}, Vector{100, -100, -100}}
	lines[6] = Line{Vector{100, -100, -100}, Vector{-100, -100, -100}}
	lines[7] = Line{Vector{-100, -100, -100}, Vector{-100, 100, -100}}
	lines[8] = Line{Vector{-100, 100, 100}, Vector{-100, 100, -100}}
	lines[9] = Line{Vector{100, 100, 100}, Vector{100, 100, -100}}
	lines[10] = Line{Vector{100, -100, 100}, Vector{100, -100, -100}}
	lines[11] = Line{Vector{-100, -100, 100}, Vector{-100, -100, -100}}
	// draw loop
	numSteps := 100
	frameChan := make(chan int, numSteps)
	for i := 0; i < numSteps; i++ {
		frameChan <- i
		// draw each frame in a separate goroutine
		go func() {
			num := <-frameChan
			canvas := NewCanvas(image.Rect(0, 0, width, height))
			canvas.Clear()
			tMatrix := RotationMatrix(float64(num*360/numSteps), Vector{1, 1, 1})
			for j := range lines {
				canvas.DrawLine(color.RGBA{0, 0, 0, 255}, lines[j].Transform(tMatrix).Projection())
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
