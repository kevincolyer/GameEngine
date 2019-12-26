package GameEngine

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	fontPath = "GameEngine/OpenSans-Regular.ttf"
	fontSize = 32
)

type TransformFunc func(x, y float64) (float64, float64)

type Context struct {
	WindowTitle       string
	WinWidth          float64
	WinHeight         float64
	Window            *sdl.Window
	Renderer          *sdl.Renderer
	Blocks            float64
	ScrnWidth         float64
	ScrnHeight        float64
	lastTick          time.Time
	font              *ttf.Font
	screenXYtransform TransformFunc
}

var ctx Context

func New(b, sw, sh float64, title string, tf TransformFunc) *Context {
	fmt.Println("starting Game engine")
	ctx = Context{
		Blocks:            b,
		ScrnWidth:         sw,
		ScrnHeight:        sh,
		WinWidth:          sw * b,
		WinHeight:         sh * b,
		WindowTitle:       title,
		lastTick:          time.Now(),
		screenXYtransform: tf,
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
	// Font init
	if err = ttf.Init(); err != nil {
		println("Font initialisation failed", err)
		os.Exit(3)
	}
	if ctx.font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		println("Failed to load font", fontPath, err)
		os.Exit(4)
	}

	return &ctx
}

// cleans up window and renderer
func (c *Context) Destroy() {
	c.Window.Destroy()
	c.Renderer.Destroy()
	c.font.Close()
	ttf.Quit()
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

// Holds rendering information for displaying ttf to a texture for renderer
type TextTexture struct {
	surface *sdl.Surface
	texture *sdl.Texture
	W       float64
	H       float64
}
func NewSdlColor(r, g, b, a float64) sdl.Color {
    return sdl.Color{uint8(r), uint8(g), uint8(b), uint8(a)}
}


type Colour struct {
    R float64
    G float64
    B float64
    A float64
}
func NewColour(r,g,b,a float64) Colour {
	return Colour{r,g,b,a}
}

const PI = 3.141592
// fades a colour by a percentage
func (c Colour) Fade(percent float64) Colour {
    percent=Clamp01(percent) // just to keep things sane!
    return Colour{R: c.R*percent, G: c.G*percent, B: c.B*percent, A: c.A}
}

func (c Colour) ToSDLColor() sdl.Color {
    return NewSdlColor(c.R,c.G,c.B,c.A)
}

func (c Colour) Unpack() (uint8,uint8,uint8,uint8) {
	return uint8(c.R),uint8(c.G),uint8(c.B),uint8(c.A) 
}


type V2D struct {
    Dx float64
    Dy float64
}

type P2D struct {
    X float64
    Y float64
//    Z float64 // for depth buffering
}


type zbuffer struct {
    buf [][]float64
    w int
    h int
}

const  NEGINF float64  = -1000000

// Clears zbuffer
func (z *zbuffer) Clear() {
    for x:=0;x<z.w;x++ {
        for y:=0;y<z.h;y++ {
            z.buf[x][y]=NEGINF
        }
    }
}

// sets z buffer to depth z if is nearer than prev val. Returns true or false
func (zb *zbuffer) SetIfNearer(x,y,z float64) bool {
    // nearer is > then NEGINF
    if zb.buf[int(x)][int(y)]>z {
        zb.buf[int(x)][int(y)]=z
        return true
    }
    return false
}

func NewZbuffer(w,h float64) *zbuffer {
    z := &zbuffer{
		buf: make( [][]float64, int(w), int(h) ) ,
        w: int(w),
        h: int(h),
    }
    z.Clear()
    return z 
}

// New texture for text display. Use .Draw to paint
func (c *Context) NewText(stringToDisplay string, clr Colour ) (t *TextTexture) {
	t = &TextTexture{}
	var err error
	if t.surface, err = c.font.RenderUTF8Blended(stringToDisplay, clr.ToSDLColor()); err != nil {
		panic("can't render font to buffer")
	}
	t.texture, err = c.Renderer.CreateTextureFromSurface(t.surface)
	t.H = float64(t.surface.H)
	t.W = float64(t.surface.W)
	t.surface.Free()
	return t
}

// Draws TextTexture to renderer at x,y with optional w,h (default to full w,h)
func (t *TextTexture) Draw(r *sdl.Renderer, x, y, w, h float64) {
	if w == 0 {
		w = t.W
	}
	if h == 0 {
		h = t.H
	}

	d := NewRect(x, y, w, h)
	s := NewRect(0, 0, w, h)
	// Draw the text around the center of the window
	r.Copy(t.texture, s, d)
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
func R256() float64 {
	return float64(rand.Intn(256))
}

// line drawing (blocks)
func (c *Context) Line(x0, y0, x1, y1 float64) {
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
		c.Point(x0, y0)
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

// Clear renderer
func (c *Context) Clear() {
	c.Renderer.Clear()
}

// Render all to screen
func (c *Context) Present() {
	c.Renderer.Present()
}

// Set Drawing color
func (c *Context) SetDrawColor(rgba Colour) {
	// color := &sdl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	c.Renderer.SetDrawColor(rgba.Unpack())
}

// Struct holding key status information
type KeyStatus struct {
	Key       string // key pressed as string
	Pressed   bool
	Released  bool
	Repeating bool
	Modifier  uint16
	Event     bool
}

// Poll quit and key events. Returns running=True and a Key struct
func (c *Context) PollQuitandKeys() (running bool, keys KeyStatus) {
	running = true
	keys.Event = false // unless something happens!
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch key := event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			running = false
		case *sdl.KeyboardEvent:
			keys.Key = string(key.Keysym.Sym)
			keys.Pressed = (key.State == sdl.PRESSED)
			keys.Released = (key.State == sdl.RELEASED)
			keys.Pressed = (key.Repeat > 0)
			keys.Modifier = key.Keysym.Mod
			keys.Event = true
		}
	}
	return running, keys
}

// Draws a point (blocks)
func (c *Context) Point(x0, y0 float64) {
	if c.screenXYtransform != nil {
		x0, y0 = c.screenXYtransform(x0, y0)
	}
	c.Renderer.FillRect(NewRect(x0*c.Blocks, y0*c.Blocks, c.Blocks, c.Blocks))
}

// Draws a point (blocks)
func (c *Context) PointScale(x0, y0, scale float64) {
	// if c.screenXYtransform != nil {
	// 	x0, y0 = c.screenXYtransform(x0, y0)
	// }
	c.Renderer.FillRect(NewRect(x0*c.Blocks/scale, y0*c.Blocks/scale, c.Blocks/scale, c.Blocks/scale))
}

// calculates the elapsed time between updates
func (c *Context) Elapsed() float64 {
	t := time.Now()
	elapsed := float64(t.Sub(c.lastTick))
	c.lastTick = t
	if elapsed == 0 {
		elapsed++
	}
	return elapsed / (1000*1000 *10)
}

func Sign(s float64) float64 {
	if s==0 { return 0 }
	if s<0 {return -1 }
	return 1
}
