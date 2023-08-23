package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 800
	screenHeight = 600
	blockSize    = 40
)

type block struct {
	x, y int32
}

var blocks []block

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Tetris", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()

		for _, block := range blocks {
			drawBlock(renderer, block.x, block.y)
		}

		renderer.Present()

		sdl.Delay(16)
	}
}

func drawBlock(renderer *sdl.Renderer, x, y int32) {
	renderer.SetDrawColor(0, 128, 255, 255)
	renderer.FillRect(&sdl.Rect{X: x, Y: y, W: blockSize, H: blockSize})
}
