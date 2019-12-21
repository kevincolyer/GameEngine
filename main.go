package main

import (
	"fmt"
	"os"

	. "github.com/kevincolyer/GameEngine/GameEngine"
	"github.com/veandco/go-sdl2/sdl"
)

// func half(x, y float64) (float64, float64) {
// 	return x / 2, y / 2
// }

func wrapScreen(x, y float64) (float64, float64) {
	return Wrap(x, 0, blocksw), Wrap(y, 0, blocksh)
}

func main() {
	var ctx = New(8, 160, 100, "Asteroids", nil)

	onCreate(ctx)
	var running = true

	for running {
		running = onUpdate(ctx, ctx.Elapsed())
	}

	ctx.Destroy()
	os.Exit(0)
}

var dx, dy float64
var x, y float64
var w, h float64
var blocks, blocksw, blocksh float64

func onCreate(c *Context) {
	fmt.Println("created")
	w = 100
	h = 100

	dx = 2
	dy = 2

	blocks = c.Blocks
	blocksw = c.ScrnWidth
	blocksh = c.ScrnHeight
	c.Renderer.Clear()
	c.Renderer.Present()

}

func onUpdate(c *Context, elapsed float64) (running bool) {
	var oldLx, oldLy, nx, ny float64
	running = true

	x += dx * elapsed
	y += dy * elapsed
	x = Wrap(x, 0, c.WinWidth)
	y = Wrap(y, 0, c.WinHeight)

	c.Renderer.SetDrawColor(R256(), R256(), R256(), 255)
	nx = RandIntN(blocksw)
	ny = RandIntN(blocksh)

	c.Line(oldLx, oldLy, nx, ny)
	oldLx = nx
	oldLy = ny

	c.Renderer.SetDrawColor(R256(), R256(), R256(), 255)
	c.Triangle(RandIntN(blocksw), RandIntN(blocksh), RandIntN(blocksw), RandIntN(blocksh), RandIntN(blocksw), RandIntN(blocksh))

	R := NewRect(x, y, w, h)
	c.Renderer.SetDrawColor(255, 127, 127, 255)
	c.Renderer.FillRect(R)
	c.Renderer.SetDrawColor(0, 0, 0, 255)
	c.Renderer.DrawRect(R)

	t := c.NewText("Hello Mum!", sdl.Color{R: 255, G: 0, B: 0, A: 255})
	t.Draw(c.Renderer, 400-(t.W/2), 300-(t.H/2), 0, 0)
	t.Draw(c.Renderer, 300-(t.W/2), 200-(t.H/2), 0, 0)
	t.Draw(c.Renderer, 400-(t.W/2), 100-(t.H/2), 0, 0)

	c.Renderer.Present()
	Delay(1)
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch key := event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			running = false
		case *sdl.KeyboardEvent:
			switch string(key.Keysym.Sym) {
			case "q":
				println("Quit")
				running = false
			case "w":
				println("w pressed")
			}
			/* 			if key.State == sdl.RELEASED {
			   				println(" key released")
			   			}
			   			if key.State == sdl.PRESSED {
			   				println(" key pressed")
			   			}
			   			if key.Repeat > 0 {
			   				println(" key repeating")
						   }
			*/
		}
	}

	return running
}
