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
	var ctx = New(2, 640, 320, "Random", nil)

	onCreate(ctx)
	var running = true

	for running {
		running = onUpdate(ctx, ctx.Elapsed())
	}

	ctx.Destroy()
	os.Exit(0)
}

var blocksw, blocksh, blocks float64

func onCreate(c *Context) {
	blocks = c.Blocks
	blocksw = c.ScrnWidth
	blocksh = c.ScrnHeight
	c.Clear()
	c.Present()
}

var prng = "0"

var BLACK = Colour{0, 0, 0, 255}
var WHITE = Colour{255, 255, 255, 225}
var lastProcGenSeed uint32 = 0
var lastGnuRandSeed uint32 = 0

func onUpdate(c *Context, elapsed float64) (running bool) {
	running, keys := c.PollQuitandKeys()
	if keys.Event {
		if keys.Key == "q" {
			running = false
		}
		if keys.Key >= "0" && keys.Key < "2" {
			prng = keys.Key
			fmt.Println(prng)
		}
		if keys.Key == " " {
			lastGnuRandSeed = uint32(gnuRand() * 0x7FFFFFFF)
			lastProcGenSeed = uint32(procGenRnd() * 0x7FFFFFFF)
		}
		fmt.Println("pressed ", keys.Key)
	}
	c.Clear()
	nProcGen = lastProcGenSeed
	gnuRandSeed = lastGnuRandSeed

	for x := 0.0; x < blocksw; x++ {
		for y := 0.0; y < blocksh; y++ {
			star := false
			if prng == "0" {
				star = procGenRnd() < 0.15
			}
			if prng == "1" {
				star = gnuRand() < 0.15
			}
			if star {
				c.SetDrawColor(WHITE)
			} else {
				c.SetDrawColor(BLACK)
			}
			c.Point(x, y)
		}
	}
	c.Present()
	Delay(1)

	return running
}

// https://github.com/OneLoneCoder/olcPixelGameEngine/blob/master/Videos/OneLoneCoder_PGE_ProcGen_Universe.cpp
var nProcGen uint32 = 0

func procGenRnd() float32 {
	nProcGen += 0xe120fc15
	var tmp uint64 = uint64(nProcGen) * 0x4a39b70d
	var m1 uint32 = uint32((tmp >> 32) ^ tmp)
	tmp = uint64(m1) * 0x12fad5c9
	var m2 uint32 = uint32((tmp >> 32) ^ tmp)
	return float32(m2) / float32(0x7FFFFFFF)
}

var gnuRandSeed uint32 = 0

// gnuRand ... pseudo random number generator - good for 32bits
func gnuRand() float32 {
	gnuRandSeed = gnuRandSeed*1103515245 + 12345
	return float32(gnuRandSeed) / float32(0x7FFFFFFF)
}
