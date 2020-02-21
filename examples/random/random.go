package main

import (
	"fmt"
	"os"

	. "github.com/kevincolyer/GameEngine/GameEngine"
)

func wrapScreen(x, y float64) (float64, float64) {
	return Wrap(x, 0, blocksw), Wrap(y, 0, blocksh)
}

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
var rand *gnuRand
var lastRandSeed uint32 = 0
var BLACK = Colour{0, 0, 0, 255}
var WHITE = Colour{255, 255, 255, 225}

func onCreate(c *Context) {
	blocks = c.Blocks
	blocksw = c.ScrnWidth
	blocksh = c.ScrnHeight
	c.Clear()
	c.Present()
	rand = NewgnuRand(0)
	lastRandSeed = rand.Seed(0)
}

func onUpdate(c *Context, elapsed float64) (running bool) {
	running, keys := c.PollQuitandKeys()
	if keys.Event {
		if keys.Key == "q" {
			running = false
		}
		if keys.Key == " " && keys.Released {
			lastRandSeed = rand.Seed(rand.Rnd())
			println(lastRandSeed)
		}
		fmt.Println("pressed ", keys.Key)
	}
	c.Clear()
	lastRandSeed = rand.Seed(lastRandSeed)

	for x := 0.0; x < blocksw; x++ {
		for y := 0.0; y < blocksh; y++ {
			star := rand.Rand() < 0.15
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



// gnuRand object pseudo random number generator - good for 32bits
type gnuRand struct {
	seed uint32
	Size float32
}

// gnuRand constructor - pseudo random number generator - good for 32bits
func NewgnuRand(s uint32) *gnuRand {
	return &gnuRand{seed: uint32(s), Size: 0x7FFFFFFF}
}

// gnuRand ... Seeds
func (g *gnuRand) Seed(s uint32) uint32 {
	g.seed = s
	return g.seed
}

// Rnd integer fullrange 0 -> upper bound
func (g *gnuRand) Rnd() uint32 {
	g.seed = g.seed*1103515245 + 12345
	return g.seed
	//return float32(gnuRandSeed) / float32(0x7FFFFFFF)
}

// Rand float between 0 and 1
func (g *gnuRand) Rand() float32 {
	return float32(g.Rnd()) / g.Size
}

// Rand float between min and max (signed)
func (g *gnuRand) RandF(min, max float32) float32 {
	return g.Rand()/(max-min) + min
}

// Rand integer between min and max (signed)
func (g *gnuRand) RandI(min, max int32) int32 {
	return int32(g.Rnd())%(max-min) + min
}
