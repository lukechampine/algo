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
	X, Y float64
}

func (v *Vector) Sub(v2 Vector) Vector {
	return Vector{v.X - v2.X, v.Y - v2.Y}
}

func (v *Vector) Length() float64 {
	return math.Hypot(v.X, v.Y)
}

type Line struct {
	V1, V2 Vector
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

func (c Canvas) DrawParaFn(step float64) {
	t := float64(0)
	for i := 0; i < 100; i++ {
		c.DrawLine(color.RGBA{0, 0, 0, 255}, Line{paraFn(t), paraFn(t + step)})
		t += step
	}
}

func paraFn(t float64) Vector {
	x := 340 * math.Exp(-0.1*t) * math.Cos(t)
	y := 340 * math.Exp(-0.1*t) * math.Sin(t)
	return Vector{x, y}
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
			canvas := NewCanvas(image.Rect(0, 0, width, height))
			canvas.Clear()
			canvas.DrawParaFn(step + float64(num)*inc)
			canvas.SaveToFile(num)
		}()
	}
	log.Print("finished computation")
	exec.Command("sh", "-c", "convert -delay 10 -loop 0 input-*.png output.gif").Run()
	log.Print("output rendered as gif")
	exec.Command("sh", "-c", "rm input-*.png").Run()
}
