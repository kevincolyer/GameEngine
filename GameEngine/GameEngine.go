package GameEngine

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Context struct {
	WindowTitle string
	WinWidth    float64
	WinHeight   float64
	Window      *sdl.Window
	Renderer    *sdl.Renderer
	Blocks      float64
	ScrnWidth   float64
	ScrnHeight  float64
	lastTick    time.Time
}

var ctx Context

func New(b, sw, sh float64, title string) *Context {
	fmt.Println("starting Game engine")
	ctx = Context{
		Blocks:      b,
		ScrnWidth:   sw,
		ScrnHeight:  sh,
		WinWidth:    sw * b,
		WinHeight:   sh * b,
		WindowTitle: title,
		lastTick:    time.Now(),
	}
	var err error
	ctx.Window, err = sdl.CreateWindow(ctx.WindowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(ctx.WinWidth), int32(ctx.WinHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit(1)
	}

	ctx.Renderer, err = sdl.CreateRenderer(ctx.Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(2)
	}

	return &ctx
}

// cleans up window and renderer
func (c *Context) Destroy() {
	c.Window.Destroy()
	c.Renderer.Destroy()
}

// helper function to build and SDL rectangle from float64's
func NewRect(x, y, w, h float64) *sdl.Rect {
	return &sdl.Rect{int32(x), int32(y), int32(w), int32(h)}
}

// Delay in milliseconds
func Delay(s uint32) {
	sdl.Delay(s)
}

// rand intn for sdl
func RandIntN(i float64) float64 {
	return float64(rand.Intn(int(i)))
}

// returns number wraped around a low and hi boundary
func Wrap(num, low, hi float64) float64 {
	if num < low {
		return hi - (low - num)
	}
	if num > hi {
		return low + (num - hi)
	}
	return num
}

// returns number clamped between low and hi
func Clamp(num, low, hi float64) float64 {
	if num < low {
		return low
	}
	if num > hi {
		return hi
	}
	return num
}

// returns number clamped between 0 and 1
func Clamp01(num float64) float64 {
	if num < 0 {
		return 0
	}
	if num > 1 {
		return 1
	}
	return num
}

// random number 0-255
func R256() uint8 {
	return uint8(rand.Intn(256))
}

// line drawing (blocks)
func (c *Context) Line(x0, y0, x1, y1 float64) {
	blocks := c.Blocks
	x0 = math.Round(x0)
	x1 = math.Round(x1)
	y0 = math.Round(y0)
	y1 = math.Round(y1)

	// from rosetta code
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}
	var sx, sy float64
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
		c.Renderer.FillRect(NewRect(x0*blocks, y0*blocks, blocks, blocks))
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

// Draws outline Triangle (blocks)
func (c *Context) Triangle(x0, y0, x1, y1, x2, y2 float64) {

	if x0 == x1 && x1 == x2 && y0 == y1 && y1 == y2 {
		c.Point(x0, y0)
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
		c.bottomFlatTriangle(x0, y0, x1, y1, x2, y2)
	} else if y0 == y1 {
		c.topFlatTriangle(x0, y0, x1, y1, x2, y2)
		// if top flat triangle

	} else {
		//get new vertex in middle of x0,y0 x2,y2 face at y1
		x3 := x0 + (y1-y0)/(y2-y0)*(x2-x0)
		y3 := y1
		//bf
		c.bottomFlatTriangle(x0, y0, x1, y1, x3, y3)
		//then
		//tf
		c.topFlatTriangle(x1, y1, x3, y3, x2, y2)
	}
}

// helper for Trinagle
func (c *Context) bottomFlatTriangle(x1, y1, x2, y2, x3, y3 float64) {
	// 	blocks:=int32(c.Blocks)

	invslope1 := (x2 - x1) / (y2 - y1)
	invslope2 := (x3 - x1) / (y3 - y1)
	curx1 := x1
	curx2 := x1

	for scanlineY := y1; scanlineY <= y2; scanlineY++ {
		c.Line(curx1, scanlineY, curx2, scanlineY)
		curx1 += invslope1
		curx2 += invslope2

	}
}

// helper for Triangle
func (c *Context) topFlatTriangle(x1, y1, x2, y2, x3, y3 float64) {

	invslope1 := (x3 - x1) / (y3 - y1)
	invslope2 := (x3 - x2) / (y3 - y2)
	curx1 := x3
	curx2 := x3

	for scanlineY := y3; scanlineY > y1; scanlineY-- {
		c.Line(curx1, scanlineY, curx2, scanlineY)
		curx1 -= invslope1
		curx2 -= invslope2
	}
}

// Draws a point (blocks)
func (c *Context) Point(x0, y0 float64) {
	c.Renderer.FillRect(NewRect(x0*c.Blocks, y0*c.Blocks, c.Blocks, c.Blocks))
}

// calculates the elapsed time between updates
func (c *Context) Elapsed() float64 {
	t := time.Now()
	elapsed := float64(t.Sub(c.lastTick))
	c.lastTick = t
	if elapsed == 0 {
		elapsed++
	}
	return 1 / elapsed * 1000 * 1000
}
