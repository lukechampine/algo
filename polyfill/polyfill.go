package main

import (
	"image/color"
	"image/color/palette"
	"runtime"
	"sync"

	. "github.com/lukechampine/algo/algo"
)

type Triangle [3]Vector

func (p Triangle) transform(m Matrix) Triangle {
	return Triangle{m.Multiply(&p[0]), m.Multiply(&p[1]), m.Multiply(&p[2])}
}

func (p Triangle) isVisible() bool {
	// a triangle is visible if the cross products of its vertices point towards the viewer
	// since we are only concerned with the Z coordinate, we can optimize a bit
	ux, uy := p[1].X-p[0].X, p[1].Y-p[0].Y
	vx, vy := p[2].X-p[0].X, p[2].Y-p[0].Y
	return (uy*vx - ux*vy) > 0
}

// algorithm taken from www-users.mat.uni.torun.pl/~wrona/3d_tutor/tri_fillers.html
func drawTri(c *Canvas, tri Triangle, col color.RGBA) {
	// compute colorIndex
	colorIndex := uint8(color.Palette(palette.Plan9).Index(col))

	// sort points by Y value
	if tri[1].Y < tri[0].Y {
		tri[1], tri[0] = tri[0], tri[1]
	}
	if tri[2].Y < tri[1].Y {
		tri[2], tri[1] = tri[1], tri[2]
	}
	if tri[1].Y < tri[0].Y {
		tri[1], tri[0] = tri[0], tri[1]
	}

	// determine linear equations for bottom half
	botFn1 := func(y float64) float64 {
		if tri[1].X == tri[0].X || tri[1].Y == tri[0].Y {
			return tri[1].X
		}
		m := (tri[1].Y - tri[0].Y) / (tri[1].X - tri[0].X)
		return (y - tri[0].Y + m*tri[0].X) / m
	}
	botFn2 := func(y float64) float64 {
		if tri[2].X == tri[0].X || tri[2].Y == tri[0].Y {
			return tri[2].X
		}
		m := (tri[2].Y - tri[0].Y) / (tri[2].X - tri[0].X)
		return (y - tri[0].Y + m*tri[0].X) / m
	}
	// sort by X value
	if botFn2(tri[0].Y+1) < botFn1(tri[0].Y+1) {
		botFn1, botFn2 = botFn2, botFn1
	}

	// draw bottom half
	for i := tri[0].Y; i <= tri[1].Y; i++ {
		l, r := botFn1(i), botFn2(i)
		// draw scanline
		for j := l; j <= r; j++ {
			c.CartSet(int(j), int(i), colorIndex)
		}
	}
	// determine linear equations for top half
	topFn1 := func(y float64) float64 {
		if tri[2].X == tri[0].X || tri[2].Y == tri[0].Y {
			return tri[2].X
		}
		m := (tri[2].Y - tri[0].Y) / (tri[2].X - tri[0].X)
		return (y - tri[0].Y + m*tri[0].X) / m
	}
	topFn2 := func(y float64) float64 {
		if tri[2].X == tri[1].X || tri[2].Y == tri[1].Y {
			return tri[2].X
		}
		m := (tri[2].Y - tri[1].Y) / (tri[2].X - tri[1].X)
		return (y - tri[2].Y + m*tri[2].X) / m
	}
	// sort by X value
	if topFn2(tri[1].Y+1) < topFn1(tri[1].Y+1) {
		topFn1, topFn2 = topFn2, topFn1
	}

	// draw top half
	for i := tri[1].Y; i <= tri[2].Y; i++ {
		l, r := topFn1(i), topFn2(i)
		// draw scanline
		for j := l; j <= r; j++ {
			c.CartSet(int(j), int(i), colorIndex)
		}
	}

	// draw outline
	c.DrawLine(Line{tri[0], tri[1]})
	c.DrawLine(Line{tri[1], tri[2]})
	c.DrawLine(Line{tri[2], tri[0]})
}

func main() {
	// use all available logical processors
	runtime.GOMAXPROCS(runtime.NumCPU())
	// canvas properties
	width, height := 700, 700
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
				tTri := figure[j].transform(tMatrix)
				if tTri.isVisible() {
					drawTri(canvas, tTri, color.RGBA{uint8(j * 20), 0, uint8(255 - j*20), 255})
				}
			}
			fw.AddFrame(canvas, num)
			wg.Done()
		}(i)
	}
	wg.Wait()

	// encode frames as a GIF
	err := fw.WriteToFile("output.gif")
	if err != nil {
		println(err.Error())
	}
}
