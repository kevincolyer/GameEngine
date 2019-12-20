package main

import (
	. "github.com/kevincolyer/GameEngine/GameEngine"
    "github.com/veandco/go-sdl2/sdl"
    "os"
    "fmt"
)

func main() {
    var ctx=New(8,100,80,"Asteroids")
    defer ctx.Destroy()
    
    onCreate(ctx);
    var running = true
    
    for running {
        running= onUpdate(ctx,ctx.Elapsed())
    }
    
	os.Exit(0)
}

var dx, dy float32
var x, y int32
var w, h int32
var blocks,blocksw,blocksh int

func onCreate(c *Context) {
    fmt.Println("created")
	w = 100
	h = 100

	dx = 2
	dy = 2

	blocks = int(c.Blocks)
	blocksw = int(c.ScrnWidth)
	blocksh = int(c.ScrnHeight)  
    
}

func onUpdate(c *Context, elapsed float32) (running bool) {
	var oldLx, oldLy, nx, ny int32
	running = true
fmt.Println(elapsed)

		x += int32(dx*elapsed)
		y += int32(dy*elapsed)
		x = Wrap(x, 0, c.WinWidth)
		y = Wrap(y, 0, c.WinHeight)
        R:=NewRect(x,y,w,h)
        
		c.Renderer.SetDrawColor(255, 127, 127, 255)
		c.Renderer.FillRect( R )
		c.Renderer.SetDrawColor(0, 0, 0, 255)
		c.Renderer.DrawRect( R )

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
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}

    return running
}


