package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	. "github.com/kevincolyer/GameEngine/GameEngine"
)

var blocksw, blocksh, blocks float64
var WHITE, BLACK, RED, GREEN, BLUE, GREY50, DARKRED, SADDLEBROWN, STEELBLUE Colour

var fps = flag.Bool("fps", false, "Display Frames per second")
var blocksi = flag.Int("blocks", 8, "Blocks of X pixels")

var x, y, z, angle float64

const ww = 10
const wh = 20

var world [wh][ww]int

const horizon float64 = 20
const screenz float64 = 0.5 // for now...

func main() {
	flag.Parse()
	blocksw = 180
	blocksh = 80
	blocks = float64(*blocksi)
	WHITE = NewColour(255, 255, 255, 255)
	BLACK = NewColour(0, 0, 0, 255)
	RED = NewColour(255, 0, 0, 255)
	DARKRED = NewColour(139, 0, 0, 255)
	GREEN = NewColour(0, 255, 0, 255)
	BLUE = NewColour(0, 0, 255, 255)
	GREY50 = NewColour(127, 127, 127, 255)
	SADDLEBROWN = NewColour(139, 69, 19, 255)
	STEELBLUE = NewColour(70, 130, 180, 255)
	var ctx = New(blocks, blocksw, blocksh, "Dogenstein", nil)

	onCreate(ctx)
	var running = true

	for running {
		running = onUpdate(ctx, ctx.Elapsed())
	}

	ctx.Destroy()
	os.Exit(0)
}

// GAME GLOBAL VARIABLES
const FOV float64 = PI / 2

var worldSpeed float64
var distScreen = 0.5

func onCreate(c *Context) {
	worldSpeed = 0.05

	c.Clear()
	c.Present()
	resetGame()
}

func resetGame() {
	x = 3
	y = 3
	angle = PI / 4
	world = [wh][ww]int{ // y then x
		[ww]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 1, 0, 1, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 1, 0, 1, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 1, 0, 1, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 1, 0, 1, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		[ww]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

}

func onUpdate(c *Context, elapsed float64) (running bool) {
	// boilerplate to start
	// println(elapsed)
	ews := elapsed * worldSpeed
	running, keys := c.PollQuitandKeys()
	if keys.Event {
		if keys.Key == "q" {
			running = false
		}
	}
	// Update code here...
	c.SetDrawColor(BLACK)
	c.Clear()
	// keys //////////////////////////////////////////
	if keys.Key == "a" {
		angle -= ews
	}
	if keys.Key == "d" {
		angle += ews
	}
	if keys.Key == "w" {
		x += math.Sin(angle) * ews
		y -= math.Cos(angle) * ews
	}
	if keys.Key == "s" {
		x -= math.Sin(angle) * ews
		y += math.Cos(angle) * ews
	}
	if keys.Key == " " {
	}

	// manipulations /////////////////////////////////////

	// Draw
	///////////////////////////////////////////////////
	screenmid := blocksh / 2
	// screenmidw := blocksw / 2
	for bx := 0.0; bx < blocksw; bx++ {
		a := angle + bx/blocksw*FOV - FOV/2
		// march a ray from screen distance to horizon
		// z is distance to hitting a block. give up at horizon
		// z = math.Sqrt(screenz*screenz + (bx-screenmidw)*(bx-screenmidw))
		z = screenz
		for z < horizon {
			z += 0.01
			tx := x + math.Sin(a)*z
			ty := y - math.Cos(a)*z
			if tx >= ww || tx < 0 || ty >= wh || ty < 0 {
				z = horizon
				break
			}
			if world[int(ty)][int(tx)] == 1 {
				// hit a wall

				break
			}
		}
		// draw from top to bottom
		wallt := math.Trunc(screenmid - (screenmid / z))
		wallb := math.Trunc(blocksh - wallt)
		c.SetDrawColor(BLACK)
		for by := 0.0; by < blocksh; by++ {
			if by >= wallt && by < wallb {
				c.SetDrawColor(GREY50.Fade(1.5 / z))
			}
			if by >= wallb {
				c.SetDrawColor(RED.Fade(1 - (blocksh-by)/screenmid))
			}
			c.Point(bx, by)
		}
	}

	// Draw text and 'top' layers
	c.SetDrawColor(WHITE)
	if *fps {
		c.DrawText(1, 17, 4, fmt.Sprintf("fps:%d", int(100/elapsed)))
	}
	c.DrawText(0, 0, 4, fmt.Sprintf("x: %v y: %v a: %v", x, y, angle))
	// boilerplate to finish
	c.Present()
	Delay(1)
	return running
}
