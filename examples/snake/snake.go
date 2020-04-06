package main

import (
	// 	"container/list"
	"flag"
	"fmt"
	"math"
	// 	"math/rand"
	"os"

	. "github.com/kevincolyer/GameEngine/GameEngine"
)

// helper function - can be passed in with GameEngine.New to modify the way blocks are drawn to the screen
func wrapScreen(x, y float64) (float64, float64) {
	return Wrap(x, 0, blocksw), Wrap(y, 0, blocksh)
}

var blocksw, blocksh, blocks float64
var WHITE, BLACK, RED, GREEN, BLUE, GREY50, DARKRED, SADDLEBROWN, BROWN, STEELBLUE Colour
var explodeShip bool

var fps = flag.Bool("fps", false, "Display Frames per second")
var blocksi = flag.Int("blocks", 4, "Blocks of X pixels")

const PI = 3.1415927
const TWOPI = PI * 2
const HALFPI = PI / 2

func main() {
	flag.Parse()
	blocksw = 320
	blocksh = 160
	blocks = float64(*blocksi)
	WHITE = NewColour(255, 255, 255, 255)
	BLACK = NewColour(0, 0, 0, 255)
	RED = NewColour(255, 0, 0, 255)
	DARKRED = NewColour(139, 0, 0, 255)
	GREEN = NewColour(0, 255, 0, 255)
	BLUE = NewColour(0, 0, 255, 255)
	GREY50 = NewColour(127, 127, 127, 255)
	SADDLEBROWN = NewColour(139, 69, 19, 255)
	BROWN = NewColour(101, 60, 15, 255)
	STEELBLUE = NewColour(70, 130, 180, 255)
	var ctx = New(blocks, blocksw, blocksh, "SSSSNake!", wrapScreen)

	onCreate(ctx)
	var running = true

	for running {
		running = onUpdate(ctx, ctx.Elapsed())
	}

	ctx.Destroy()
	os.Exit(0)
}

var player snake

type snake struct {
	segments  []seg
	size      float64
	length    float64
	direction float64
	speed     float64
}

type seg struct {
	x float64
	y float64
}

// GAME GLOBAL VARIABLES
var score int32
var hiscore int32 = 0
var worldSpeed float64

func onCreate(c *Context) {
	worldSpeed = 1

	c.Clear()
	c.Present()
	resetGame()
	player = snake{size: 5.0, length: 5.0, direction: PI + HALFPI, speed: 0.25}
	player.segments = make([]seg, 200, 200)
	x := blocksw / 2
	y := blocksh / 2
	for i := 0.0; i < player.length; i++ {
		player.segments[int(i)].x = x + i*player.size*1.5
		player.segments[int(i)].y = y
	}
}

func resetGame() {
}

func onUpdate(c *Context, elapsed float64) (running bool) {
	// boilerplate to start
	// println(elapsed)
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
		player.direction = player.direction + 0.5*elapsed*worldSpeed
	}
	if keys.Key == "d" {
		player.direction = player.direction - 0.5*elapsed*worldSpeed
	}
	if keys.Key == "n" {
		// add new section
		player.length++
		player.segments[int(player.length-1)].x = player.segments[int(player.length-2)].x
		player.segments[int(player.length-1)].y = player.segments[int(player.length-2)].y
	}
	// update head
	x := player.segments[0].x
	y := player.segments[0].y
	d := player.direction
	player.segments[0].x = x + math.Sin(d)*player.speed*worldSpeed*elapsed
	player.segments[0].y = y + math.Cos(d)*player.speed*worldSpeed*elapsed

	player.segments[0].x, player.segments[0].y = wrapScreen(player.segments[0].x, player.segments[0].y)
	// update body - keep min of player.size from previous segment
	for j := 1.0; j < player.length; j++ {
		i := int(j)
		dx := player.segments[i-1].x - player.segments[i].x
		dy := player.segments[i-1].y - player.segments[i].y
		// only wrap if moderate distance from player
		// 		wrap := dx*dx+dy*dy < player.size*1.5*player.size*1.5

		a := math.Atan2(dx, dy)
		ndx := math.Sin(a) * 1.5 * player.size
		ndy := math.Cos(a) * 1.5 * player.size
		player.segments[i].x = player.segments[i-1].x - ndx
		player.segments[i].y = player.segments[i-1].y - ndy
		/*
			if wrap {
				player.segments[i].x, player.segments[i].y = wrapScreen(player.segments[i].x, player.segments[i].y)
			}*/

	}
	// Draw
	///////////////////////////////////////////////////
	// draw bottom layers

	for i := player.length - 1; i >= 1.0; i-- {
		if math.Mod(i, 2) == 0 {
			c.SetDrawColor(SADDLEBROWN)
		} else {
			c.SetDrawColor(BROWN)
		}
		c.DrawFillCircle(player.segments[int(i)].x, player.segments[int(i)].y, player.size)
	}
	c.SetDrawColor(WHITE)
	c.DrawFillCircle(player.segments[0].x, player.segments[0].y, player.size)
	c.SetDrawColor(SADDLEBROWN)
	// eyes
	c.Point(player.segments[0].x+math.Cos(player.direction)*2, player.segments[0].y-math.Sin(player.direction)*2)
	c.Point(player.segments[0].x-math.Cos(player.direction)*2, player.segments[0].y+math.Sin(player.direction)*2)

	// Draw text and 'top' layers
	c.SetDrawColor(DARKRED)
	c.DrawText(1, 1, 2, fmt.Sprintf("hi:%v score:%v", hiscore, score))
	if *fps {
		c.DrawText(1, 17, 4, fmt.Sprintf("fps:%d", int(100/elapsed)))
	}
	// boilerplate to finish
	c.Present()
	Delay(1)
	return running
}
