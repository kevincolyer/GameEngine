package GameEngine

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"os"
    "time"
)

type Context struct {
    WindowTitle string
    WinWidth int32
    WinHeight int32
	Window *sdl.Window
	Renderer *sdl.Renderer
    Blocks int64
    ScrnWidth int64
    ScrnHeight int64
       }


var ctx Context 

func New(b,sw,sh int, title string) *Context {
    fmt.Println("starting Game engine")
    ctx = Context{
        Blocks: int64(b),
        ScrnWidth: int64(sw),
        ScrnHeight: int64(sh),
        WinWidth: int32(sw*b),
        WinHeight: int32(sh*b),
        WindowTitle: title,
    };
    var err error
    ctx.Window, err = sdl.CreateWindow(ctx.WindowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		ctx.WinWidth, ctx.WinHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit( 1)
	}

	ctx.Renderer, err = sdl.CreateRenderer(ctx.Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit( 2)
	}

    return &ctx
}

func (c *Context) Destroy() {
	c.Window.Destroy()
	c.Renderer.Destroy()
    
}

func NewRect(x, y, w, h int32) *sdl.Rect {
    return &sdl.Rect{x, y, w, h}
}

func Delay(s uint32) {
    sdl.Delay(s)
}

// rand intn for sdl
func RandIntN(i int) int32 {
    return int32(rand.Intn(i))
}


func Wrap(num, low, hi int32) int32 {
	if num < low {
		return hi
	}
	if num > hi {
		return low
	}
	return num
}

func Clamp(num, low, hi int32) int32 {
	if num < low {
		return low
	}
	if num > hi {
		return hi
	}
	return num
}

func Clamp01(num float64) float64 {
	if num < 0 {
		return 0
	}
	if num > 1 {
		return 1
	}
	return num
}

func R256() uint8 {
	return uint8(rand.Intn(256))
}

func (c *Context) Line(x0, y0, x1, y1 int32) {
    blocks:=int32(c.Blocks)
    
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
		c.Renderer.FillRect(&sdl.Rect{x0 * blocks, y0 * blocks, blocks, blocks})
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

func (c *Context) Triangle(x0, y0, x1, y1, x2, y2 int32) {
// 	blocks:=int32(c.Blocks)
    
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
		x3 := int32(float64(x0) + (float64(y1-y0) / float64(y2-y0) * float64(x2-x0)))
		y3 := y1
		//bf
		c.bottomFlatTriangle(x0, y0, x1, y1, x3, y3)
		//then
		//tf
		c.topFlatTriangle(x1, y1, x3, y3, x2, y2)
	}
}

func (c *Context) bottomFlatTriangle(x1, y1, x2, y2, x3, y3 int32) {
// 	blocks:=int32(c.Blocks)
	invslope1 := float64(x2-x1) / float64(y2-y1)
	invslope2 := float64(x3-x1) / float64(y3-y1)
	curx1 := float64(x1)
	curx2 := float64(x1)

	for scanlineY := y1; scanlineY <= y2; scanlineY++ {
		c.Line(int32(curx1), scanlineY, int32(curx2), scanlineY)
		curx1 += invslope1
		curx2 += invslope2
	}
}

func (c *Context) topFlatTriangle(x1, y1, x2, y2, x3, y3 int32) {
// 	blocks:=int32(c.Blocks)
	invslope1 := float64(x3-x1) / float64(y3-y1)
	invslope2 := float64(x3-x2) / float64(y3-y2)
	curx1 := float64(x3)
	curx2 := float64(x3)

	for scanlineY := y3; scanlineY > y1; scanlineY-- {
		c.Line(int32(curx1), scanlineY, int32(curx2), scanlineY)
		curx1 -= invslope1
		curx2 -= invslope2
	}
}

func (c *Context) Point(x0, y0 int32) {
    blocks:=int32(c.Blocks)
	c.Renderer.FillRect(&sdl.Rect{x0 * blocks, y0 * blocks, blocks, blocks})
}

var lastTick=time.Now()

func (c *Context) Elapsed() float32 {
    t:=time.Now()
    elapsed := t.Sub(lastTick)
    lastTick=t
    if elapsed==0 { elapsed++ }
    return float32(1/float32(elapsed)*1000*1000)
}
