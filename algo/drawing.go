package algo

import (
	"image"
	"image/color"
	"image/color/palette"
)

type Canvas struct {
	image.Paletted
}

func NewCanvas(width, height int) *Canvas {
	canvas := new(Canvas)
	canvas.Paletted = *image.NewPaletted(image.Rect(0, 0, width, height), palette.Plan9)
	canvas.Clear()
	return canvas
}

func (c *Canvas) Clear() {
	for i := range c.Pix {
		c.Pix[i] = 0xFF // white
	}
}

// CartSet adjusts its input x and y coordinates to fit a Cartesian plane.
// In other words, it places (0,0) at the center of the Canvas.
func (c *Canvas) CartSet(x, y int, colorIndex uint8) {
	c.SetColorIndex(c.Bounds().Size().X/2+x, c.Bounds().Size().Y/2-y, colorIndex)
}

func (c *Canvas) DrawLine(l Line) {
	length := l.V2.Sub(l.V1).Length()
	unit := l.V2.Sub(l.V1).Unit()
	for i := float64(0); i < length; i++ {
		x := l.V1.X + i*unit.X
		y := l.V1.Y + i*unit.Y
		c.CartSet(int(x), int(y), 0x00) // black
	}
}

func (c *Canvas) DrawTriangle(tri Triangle) {
	c.DrawLine(Line{tri[0], tri[1]})
	c.DrawLine(Line{tri[1], tri[2]})
	c.DrawLine(Line{tri[2], tri[0]})
}

// algorithm taken from www-users.mat.uni.torun.pl/~wrona/3d_tutor/tri_fillers.html
func (c *Canvas) FillTriangle(tri Triangle, col color.RGBA) {
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
}
