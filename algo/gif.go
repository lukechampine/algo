package algo

import (
	"image"
	"image/gif"
	"os"
	"sync"
)

type FrameWriter struct {
	gif.GIF
	sync.Mutex
}

func NewFrameWriter(numFrames int) *FrameWriter {
	fw := new(FrameWriter)
	fw.Image = make([]*image.Paletted, numFrames)
	fw.Delay = make([]int, numFrames)
	for i := range fw.Delay {
		fw.Delay[i] = 10
	}
	return fw
}

func (fw *FrameWriter) AddFrame(c *Canvas, index int) {
	fw.Lock()
	fw.Image[index] = &c.Paletted
	fw.Unlock()
}

func (fw *FrameWriter) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	return gif.EncodeAll(file, &fw.GIF)
}
