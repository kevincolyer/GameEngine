package main

import (
	"fmt"
	"os"

	. "github.com/kevincolyer/GameEngine/GameEngine"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	var ctx = New(8, 160, 100, "Asteroids")
	defer ctx.Destroy()

	onCreate(ctx)
	var running = true

	for running {
		running = onUpdate(ctx, ctx.Elapsed())
	}

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

}

func onUpdate(c *Context, elapsed float64) (running bool) {
	var oldLx, oldLy, nx, ny float64
	running = true

	x += dx * elapsed
	y += dy * elapsed
	x = Wrap(x, 0, c.WinWidth)
	y = Wrap(y, 0, c.WinHeight)

	R := NewRect(x, y, w, h)
	c.Renderer.SetDrawColor(255, 127, 127, 255)
	c.Renderer.FillRect(R)
	c.Renderer.SetDrawColor(0, 0, 0, 255)
	c.Renderer.DrawRect(R)

	c.Renderer.SetDrawColor(R256(), R256(), R256(), 255)
	nx = RandIntN(blocksw)
	ny = RandIntN(blocksh)

	c.Line(oldLx, oldLy, nx, ny)
	oldLx = nx
	oldLy = ny

	c.Renderer.SetDrawColor(R256(), R256(), R256(), 255)
	c.Triangle(RandIntN(blocksw), RandIntN(blocksh), RandIntN(blocksw), RandIntN(blocksh), RandIntN(blocksw), RandIntN(blocksh))

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
