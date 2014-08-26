package algo

import (
	"image"
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
