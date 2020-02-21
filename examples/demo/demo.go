package main

import (
	"fmt"
	"os"

	. "github.com/kevincolyer/GameEngine/GameEngine"
)

func wrapScreen(x, y float64) (float64, float64) {
	return Wrap(x, 0, blocksw), Wrap(y, 0, blocksh)
}


var coin *Sprite
var dungeon *SpriteSheet
var err error

func main() {
	var ctx = New(8, 160, 80, "Demo", nil)

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
var blocksw, blocksh, blocks float64

var coinx, coiny float64

func onCreate(c *Context) {
	fmt.Println("created")
	w = 100
	h = 100

	dx = 2
	dy = 2

	blocks = c.Blocks
	blocksw = c.ScrnWidth
	blocksh = c.ScrnHeight
	c.Clear()
	c.Present()
	coin, err = NewSprite("../../assets/coin2.png")
	if err != nil {
		panic(err)
	}
	dungeon, err = NewSpriteSheet("../../assets/set-cave_bright.png", 10, 2)
	if err != nil {
		panic(err)
	}
}

var tick float64

func onUpdate(c *Context, elapsed float64) (running bool) {
	var oldLx, oldLy, nx, ny float64
	tick += elapsed
	x += dx * elapsed
	y += dy * elapsed
	x = Wrap(x, 0, c.WinWidth)
	y = Wrap(y, 0, c.WinHeight)

	c.SetDrawColor(Colour{R256(), R256(), R256(), 255})
	nx = RandIntN(blocksw)
	ny = RandIntN(blocksh)

	c.Line(oldLx, oldLy, nx, ny)
	oldLx = nx
	oldLy = ny

	c.SetDrawColor(Colour{R256(), R256(), R256(), 255})
	c.Triangle(RandIntN(blocksw), RandIntN(blocksh), RandIntN(blocksw), RandIntN(blocksh), RandIntN(blocksw), RandIntN(blocksh))

	R := NewRect(x, y, w, h)
	c.SetDrawColor(Colour{255, 127, 127, 255})
	c.Renderer.FillRect(R)
	c.SetDrawColor(Colour{0, 0, 0, 255})
	c.Renderer.DrawRect(R)

	c.SetDrawColor(Colour{R: 255, G: 0, B: 0, A: 255})
	c.DrawText(1, 1, 2, "Hello Mum!")
	c.DrawText(1, 20, 1, "Hello Mum!")

	coinx = x / blocks
	coiny = y / blocks
	coin.DrawSprite(c, coinx, coiny)
	coin.DrawPartialSprite(c, coinx+20, coiny, 8, 8, 4, 4)
	dungeon.DrawSpriteFromSheet(c, coinx, coiny+20, 0, 0)
	dungeon.DrawSpriteFromSheetI(c, coinx+20, coiny+40, tick)

	c.Present()
	Delay(1)
	running, keys := c.PollQuitandKeys()
	if keys.Event {
		if keys.Key == "q" {
			running = false
		}
		println(keys.Key, " pressed")
	}

	return running
}
