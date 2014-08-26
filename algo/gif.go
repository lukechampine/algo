package algo

import (
	"image"
	"image/gif"
	"os"
	"runtime"
	"sync"
)

func init() {
	// use all available logical processors
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type FrameWriter struct {
	gif.GIF
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

func (fw *FrameWriter) GenerateFrames(frameFunc func(int) *Canvas) {
	var wg sync.WaitGroup
	wg.Add(len(fw.Image))
	for i := range fw.Image {
		go func(index int) {
			fw.Image[index] = &frameFunc(index).Paletted
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func (fw *FrameWriter) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	return gif.EncodeAll(file, &fw.GIF)
}
