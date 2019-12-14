package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"os"
)

var winTitle string = "Go-SDL2 Render"
var winWidth, winHeight int32 = 800, 600

func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	// 	var points []sdl.Point
	// 	var rect sdl.Rect
	// 	var rects []sdl.Rect

	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer renderer.Destroy()

	var dx, dy, x, y int32
	var w, h int32
	w = 100
	h = 100

	dx = 1
	dy = 1

	blocks := 2
	blocksw := int(winWidth) / blocks
	blocksh := int(winHeight) / blocks

	running := true
	var oldLx, oldLy, nx, ny int32
	for running {

		//         for i:=0;i<1000;i++ {
		//
		//         renderer.SetDrawColor(r256(),r256(),r256(),255)
		//         renderer.FillRect(&sdl.Rect{int32(rand.Intn(blocksw)*blocks),int32(rand.Intn(blocksh)*blocks),int32(blocks),int32(blocks)})
		//         }

		x += dx
		y += dy
		x = wrap(x, 0, winWidth)
		y = wrap(y, 0, winHeight)

		renderer.SetDrawColor(255, 127, 127, 255)
		renderer.FillRect(&sdl.Rect{x, y, w, h})
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.DrawRect(&sdl.Rect{x, y, w, h})

		renderer.SetDrawColor(r256(), r256(), r256(), 255)
		nx = int32(rand.Intn(blocksw))
		ny = int32(rand.Intn(blocksh))
		Line(oldLx, oldLy, nx, ny, int32(blocks), renderer)
		oldLx = nx
		oldLy = ny

		renderer.SetDrawColor(r256(), r256(), r256(), 255)
		Triangle(int32(rand.Intn(blocksw)), int32(rand.Intn(blocksh)), int32(rand.Intn(blocksw)), int32(rand.Intn(blocksh)), int32(rand.Intn(blocksw)), int32(rand.Intn(blocksh)), int32(blocks), renderer)

		renderer.Present()

		sdl.Delay(1)
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}

	return 0
}

func wrap(num, low, hi int32) int32 {
	if num < low {
		return hi
	}
	if num > hi {
		return low
	}
	return num
}

func clamp(num, low, hi int32) int32 {
	if num < low {
		return low
	}
	if num > hi {
		return hi
	}
	return num
}

func clamp01(num float64) float64 {
	if num < 0 {
		return 0
	}
	if num > 1 {
		return 1
	}
	return num
}

func r256() uint8 {
	return uint8(rand.Intn(256))
}

func Line(x0, y0, x1, y1, blocks int32, r *sdl.Renderer) {
	// from rosetta code
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}
	var sx, sy int32
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	for {
		r.FillRect(&sdl.Rect{x0 * blocks, y0 * blocks, blocks, blocks})
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func Triangle(x0, y0, x1, y1, x2, y2, blocks int32, r *sdl.Renderer) {
	if x0 == x1 && x1 == x2 && y0 == y1 && y1 == y2 {
		Point(x0, y0, blocks, r)
		return
	}
	// sort verticies in ascending order
	if y0 > y1 {
		x1, x0 = x0, x1
		y1, y0 = y0, y1
	}
	if y0 > y2 {
		x2, x0 = x0, x2
		y2, y0 = y0, y2
	}
	if y1 > y2 {
		x2, x1 = x1, x2
		y2, y1 = y1, y2
	}
	// if bottom flat triangle
	if y1 == y2 {
		bottomFlatTriangle(x0, y0, x1, y1, x2, y2, blocks, r)
	} else if y0 == y1 {
		topFlatTriangle(x0, y0, x1, y1, x2, y2, blocks, r)
		// if top flat triangle

	} else {
		//get new vertex in middle of x0,y0 x2,y2 face at y1
		x3 := int32(float64(x0) + (float64(y1-y0) / float64(y2-y0) * float64(x2-x0)))
		y3 := y1
		//bf
		bottomFlatTriangle(x0, y0, x1, y1, x3, y3, blocks, r)
		//then
		//tf
		topFlatTriangle(x1, y1, x3, y3, x2, y2, blocks, r)
	}
}

func bottomFlatTriangle(x1, y1, x2, y2, x3, y3, blocks int32, r *sdl.Renderer) {
	invslope1 := float64(x2-x1) / float64(y2-y1)
	invslope2 := float64(x3-x1) / float64(y3-y1)
	curx1 := float64(x1)
	curx2 := float64(x1)

	for scanlineY := y1; scanlineY <= y2; scanlineY++ {
		Line(int32(curx1), scanlineY, int32(curx2), scanlineY, blocks, r)
		curx1 += invslope1
		curx2 += invslope2
	}
}

func topFlatTriangle(x1, y1, x2, y2, x3, y3, blocks int32, r *sdl.Renderer) {
	invslope1 := float64(x3-x1) / float64(y3-y1)
	invslope2 := float64(x3-x2) / float64(y3-y2)
	fmt.Println("   tf invslope1 ", invslope1)
	fmt.Println("   tf invslope2 ", invslope2)
	curx1 := float64(x3)
	curx2 := float64(x3)

	for scanlineY := y3; scanlineY > y1; scanlineY-- {
		Line(int32(curx1), scanlineY, int32(curx2), scanlineY, blocks, r)
		curx1 -= invslope1
		curx2 -= invslope2
	}
}

func Point(x0, y0, blocks int32, r *sdl.Renderer) {
	r.FillRect(&sdl.Rect{x0 * blocks, y0 * blocks, blocks, blocks})
}

func main() {
	os.Exit(run())
}
