package main

import (
	"image/color"
	"log"
	"os/exec"
	"runtime"
)

type Polygon [3]Vector

func (p Polygon) Transform(m Matrix) Polygon {
	return Polygon{m.Multiply(p[0]), m.Multiply(p[1]), m.Multiply(p[2])}
}

func (p Polygon) Projection() Polygon {
	return Polygon{Vector{p[0].X, p[0].Y, 0}, Vector{p[1].X, p[1].Y, 0}, Vector{p[2].X, p[2].Y, 0}}
}

// algorithm taken from www-users.mat.uni.torun.pl/~wrona/3d_tutor/tri_fillers.html
func (c Canvas) DrawPoly(poly Polygon, color color.RGBA) {
	// sort points by y value
	if poly[1].Y < poly[0].Y {
		poly[1], poly[0] = poly[0], poly[1]
	}
	if poly[2].Y < poly[1].Y {
		poly[2], poly[1] = poly[1], poly[2]
	}
	if poly[1].Y < poly[0].Y {
		poly[1], poly[0] = poly[0], poly[1]
	}

	// calculate delta values
	l, r := poly[0], poly[0]
	dx1, dx2, dx3 := 0.0, 0.0, 0.0
	if poly[1].Y-poly[0].Y > 0 {
		dx1 = (poly[1].X - poly[0].X) / (poly[1].Y - poly[0].Y)
	}
	if poly[2].Y-poly[0].Y > 0 {
		dx2 = (poly[2].X - poly[0].X) / (poly[2].Y - poly[0].Y)
	}
	if poly[2].Y-poly[1].Y > 0 {
		dx3 = (poly[2].X - poly[1].X) / (poly[2].Y - poly[1].Y)
	}

	// main draw routine
	if dx1 > dx2 {
		// draw bottom half
		for ; l.Y <= poly[1].Y; l.Y, l.X, r.Y, r.X = l.Y+1, l.X+dx2, r.Y+1, r.X+dx1 {
			for i := l.X; i < r.X; i++ {
				c.SmartSet(int(i), int(l.Y), color)
			}
		}
		// draw top half
		r = poly[1]
		for ; l.Y <= poly[2].Y; l.Y, l.X, r.Y, r.X = l.Y+1, l.X+dx2, r.Y+1, r.X+dx3 {
			for i := l.X; i < r.X; i++ {
				c.SmartSet(int(i), int(l.Y), color)
			}
		}
	} else {
		for ; l.Y <= poly[1].Y; l.Y, l.X, r.Y, r.X = l.Y+1, l.X+dx1, r.Y+1, r.X+dx2 {
			for i := l.X; i < r.X; i++ {
				c.SmartSet(int(i), int(l.Y), color)
			}
		}
		l = poly[1]
		for ; l.Y <= poly[2].Y; l.Y, l.X, r.Y, r.X = l.Y+1, l.X+dx3, r.Y+1, r.X+dx2 {
			for i := l.X; i < r.X; i++ {
				c.SmartSet(int(i), int(l.Y), color)
			}
		}
	}

	// draw outline
	c.DrawLine(Line{poly[0], poly[1]})
	c.DrawLine(Line{poly[1], poly[2]})
	c.DrawLine(Line{poly[2], poly[0]})
}

func main() {
	// use all available logical processors
	runtime.GOMAXPROCS(runtime.NumCPU())
	// canvas properties
	width, height := 700, 700
	// figure
	figure := [12]Polygon{
		// front
		Polygon{Vector{-100, 100, 100}, Vector{100, 100, 100}, Vector{-100, -100, 100}},
		Polygon{Vector{100, -100, 100}, Vector{100, 100, 100}, Vector{-100, -100, 100}},
		// back
		Polygon{Vector{-100, 100, -100}, Vector{100, 100, -100}, Vector{-100, -100, -100}},
		Polygon{Vector{100, -100, -100}, Vector{100, 100, -100}, Vector{-100, -100, -100}},
		// top
		Polygon{Vector{-100, 100, 100}, Vector{100, 100, 100}, Vector{-100, 100, -100}},
		Polygon{Vector{100, 100, -100}, Vector{100, 100, 100}, Vector{-100, 100, -100}},
		// bottom
		Polygon{Vector{-100, -100, 100}, Vector{100, -100, 100}, Vector{-100, -100, -100}},
		Polygon{Vector{100, -100, -100}, Vector{100, -100, 100}, Vector{-100, -100, -100}},
		// left
		Polygon{Vector{100, -100, 100}, Vector{100, 100, 100}, Vector{100, -100, -100}},
		Polygon{Vector{100, 100, -100}, Vector{100, 100, 100}, Vector{100, -100, -100}},
		// right
		Polygon{Vector{-100, -100, 100}, Vector{-100, 100, 100}, Vector{-100, -100, -100}},
		Polygon{Vector{-100, 100, -100}, Vector{-100, 100, 100}, Vector{-100, -100, -100}}}
	// draw loop
	numSteps := 50
	frameChan := make(chan int, numSteps)
	for i := 0; i < numSteps; i++ {
		frameChan <- i
		// draw each frame in a separate goroutine
		go func() {
			num := <-frameChan
			canvas := NewCanvas(width, height)
			tMatrix := RotationMatrix(float64(num*360/numSteps), Vector{1, 1, 1})
			for j := range figure {
				canvas.DrawPoly(figure[j].Transform(tMatrix).Projection(), color.RGBA{uint8(j * 20), 0, uint8(255 - j*20), 255})
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
