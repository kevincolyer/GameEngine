package GameEngine

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/pbnjay/pixfont"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	//"github.com/veandco/go-sdl2/ttf"
)

// TransformFunc ... type of function to modify Point drawing
type TransformFunc func(x, y float64) (float64, float64)

// Context ... contains context object for game engine
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
	screenXYtransform TransformFunc
}

// New ... create the GameEngine and initialises
func New(b, sw, sh float64, title string, tf TransformFunc) *Context {
	fmt.Println("starting Game engine")
	ctx := Context{
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

	return &ctx
}

// Destroy ... cleans up window and renderer
func (c *Context) Destroy() {
	c.Window.Destroy()
	c.Renderer.Destroy()
	// c.font.Close()
	ttf.Quit()
}

// NewRect ... helper function to build and SDL rectangle from float64's
func NewRect(x, y, w, h float64) *sdl.Rect {
	return &sdl.Rect{int32(x), int32(y), int32(w), int32(h)}
}

// Delay in milliseconds
func Delay(s uint32) {
	sdl.Delay(s)
}

// RandIntN ... rand intn for sdl
func RandIntN(i float64) float64 {
	return float64(rand.Intn(int(i)))
}

// NewSdlColor ... takes floats and returs an sdl suitalbe color object
func NewSdlColor(r, g, b, a float64) sdl.Color {
	return sdl.Color{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// Colour struct for float64 colour values
type Colour struct {
	R float64
	G float64
	B float64
	A float64
}

// NewColour gives a new colour object
func NewColour(r, g, b, a float64) Colour {
	return Colour{r, g, b, a}
}

// PI constant
const PI = 3.141592

// Fade a colour by a percentage
func (c Colour) Fade(normalised float64) Colour {
	normalised = Clamp01(normalised) // just to keep things sane!
	return Colour{R: c.R * normalised, G: c.G * normalised, B: c.B * normalised, A: c.A}
}

// ToSDLColor ... returns a sdl.Color struct from a Colour struct
func (c Colour) ToSDLColor() sdl.Color {
	return NewSdlColor(c.R, c.G, c.B, c.A)
}

// Unpack ... transforms Colour struct to 4 uint8 values
func (c Colour) Unpack() (uint8, uint8, uint8, uint8) {
	return uint8(c.R), uint8(c.G), uint8(c.B), uint8(c.A)
}

// V2D ... struct for holding a 2D vector
type V2D struct {
	Dx float64
	Dy float64
}

// V3D ... struct for holding a 2D vector
type V3D struct {
	DX float64
	DY float64
	DZ float64
}

// P2D ... Struct for holding a 2D point
type P2D struct {
	X float64
	Y float64
	// for depth buffering?
}

// P3D ... Struct for holding a 3D point
type P3D struct {
	X float64
	Y float64
	Z float64
}

// ZBuffer struct
type ZBuffer struct {
	buf []float64
	w   int
	h   int
}

// NEGINF ... helper const = -1_000_000
const NEGINF float64 = -1000000

// NewZBuffer returns a pointer to a initilised and cleared z buffer
func NewZBuffer(w, h float64) *ZBuffer {
	z := &ZBuffer{
		buf: make([]float64, int(h)*int(w)),
		w:   int(w),
		h:   int(h),
	}
	z.Clear()
	return z
}

// Clear ... ZBuffer
func (z *ZBuffer) Clear() {
	for x := range z.buf {
		z.buf[x] = NEGINF
	}
}

// SetIfNearer ... sets z buffer to depth d if is nearer than prev val. Returns true or false
func (z *ZBuffer) SetIfNearer(x, y, d float64) bool {
	// nearer is > then NEGINF
	if z.buf[int(x)+int(y)*z.w] < d {
		z.buf[int(x)+int(y)*z.w] = d
		return true
	}
	return false
}

// DrawText to screen with a scaling factor to reduce
func (c *Context) DrawText(x, y, scale float64, text string) {
	pixfonttest := &pixfont.StringDrawable{}
	pixfont.DrawString(pixfonttest, 00, 00, text, color.White)
	x1 := x
	y1 := y
	for _, i := range []byte(pixfonttest.String()) {
		if i == 10 {
			y1++
			x1 = x
			continue
		}
		if i == 'X' {
			c.PointScale(x1, y1, scale)
		}
		x1++
	}
}

// Wrap ... returns number wraped around a low and hi boundary
func Wrap(num, low, hi float64) float64 {
	if num < low {
		return hi - (low - num)
	}
	if num > hi {
		return low + (num - hi)
	}
	return num
}

// Clamp ... returns number clamped between low and hi
func Clamp(num, low, hi float64) float64 {
	if num < low {
		return low
	}
	if num > hi {
		return hi
	}
	return num
}

// Clamp01 returns number clamped between 0 and 1
func Clamp01(num float64) float64 {
	if num < 0 {
		return 0
	}
	if num > 1 {
		return 1
	}
	return num
}

// R256 random number 0-255
func R256() float64 {
	return float64(rand.Intn(256))
}

// Line drawing (blocks)
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

// Triangle Draws outline Triangle (blocks)
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

// helper for Triangle
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

// Present ... Renders all to screen
func (c *Context) Present() {
	c.Renderer.Present()
}

// SetDrawColor for next use to Colour struct
func (c *Context) SetDrawColor(rgba Colour) {
	// color := &sdl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	c.Renderer.SetDrawColor(rgba.Unpack())
}

// KeyStatus  ... Struct holding key status information
type KeyStatus struct {
	Key       string // key pressed as string
	Pressed   bool
	Released  bool
	Repeating bool
	Modifier  uint16
	Event     bool
}

// PollQuitandKeys ... checks for events. Returns running=True and a Key struct
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

// DrawCircle ... draws circle using point function
func (c *Context) DrawCircle(x, y, radius float64) {
	// from Rosetta Code
	// Circle plots a circle with center x, y and radius r.
	// Limiting behavior:
	// r < 0 plots no pixels.
	// r = 0 plots a single pixel at x, y.
	// r = 1 plots four pixels in a diamond shape around the center pixel at x, y.
	if radius < 0 {
		return
	}
	// Bresenham algorithm
	x1, y1, err := -radius, 0.0, 2-2*radius
	for {
		c.Point(x-x1, y+y1)
		c.Point(x+y1, y+x1)
		c.Point(x+x1, y-y1)
		c.Point(x-y1, y-x1)
		radius = err
		if radius > x1 {
			x1++
			err += x1*2 + 1
		}
		if radius <= y1 {
			y1++
			err += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}
}

// DrawFillCircle ... draws filled circle using point function
func (c *Context) DrawFillCircle(x, y, radius float64) {
	// from Rosetta Code - adapted to draw fill circle
	// Circle plots a circle with center x, y and radius r.
	// Limiting behavior:
	// r < 0 plots no pixels.
	// r = 0 plots a single pixel at x, y.
	// r = 1 plots four pixels in a diamond shape around the center pixel at x, y.
	if radius < 0 {
		return
	}
	// Bresenham algorithm
	x1, y1, err := -radius, 0.0, 2-2*radius
	for {
		c.Line(x-x1, y+y1, x, y+y1)
		c.Line(x+y1, y+x1, x, y+x1)
		c.Line(x+x1, y-y1, x, y-y1)
		c.Line(x-y1, y-x1, x, y-x1)
		radius = err
		if radius > x1 {
			x1++
			err += x1*2 + 1
		}
		if radius <= y1 {
			y1++
			err += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}
}

// Point ... Draws a blocky point transformed to screen with optional transform applied from func stored in Context
func (c *Context) Point(x0, y0 float64) {
	if c.screenXYtransform != nil {
		x0, y0 = c.screenXYtransform(x0, y0)
	}
	c.Renderer.FillRect(NewRect(x0*c.Blocks, y0*c.Blocks, c.Blocks, c.Blocks))
}

// PointScale ... Draws a blocky point but scaled down by a factore (used mainly in text drawing) (blocks)
func (c *Context) PointScale(x0, y0, scale float64) {
	c.Renderer.FillRect(NewRect(x0*c.Blocks/scale, y0*c.Blocks/scale, c.Blocks/scale, c.Blocks/scale))
}

// Elapsed ... calculates the elapsed time between updates
func (c *Context) Elapsed() float64 {
	t := time.Now()
	elapsed := float64(t.Sub(c.lastTick))
	c.lastTick = t
	if elapsed == 0 {
		elapsed++
	}
	return elapsed / (1000 * 1000 * 10)
}

// Sign ... helper func - returns sign of input as -1,0,1
func Sign(s float64) float64 {
	if s == 0 {
		return 0
	}
	if s < 0 {
		return -1
	}
	return 1
}

// Sprite struct
type Sprite struct {
	image.Image
	W float64
	H float64
}

// NewSprite ... loads and builds a new sprite given a filename. Returns pointer to sprite structure
func NewSprite(filename string) (s *Sprite, err error) {
	infile, err := os.Open(filename)
	if err != nil {
		return
	}
	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	i, err := png.Decode(infile)
	if err != nil {
		return
	}
	bounds := i.Bounds()
	s = &Sprite{i, float64(bounds.Max.X - bounds.Min.X), float64(bounds.Max.Y - bounds.Min.Y)}
	return
}

// DrawSprite ... at x,y location
func (s *Sprite) DrawSprite(c *Context, x, y float64) {
	for i := 0.0; i < s.W; i++ {
		for j := 0.0; j < s.H; j++ {
			r, g, b, a := s.At(int(i), int(j)).RGBA()
			// no blending for now!
			if a > 0 {
				c.Renderer.SetDrawColor(uint8(r), uint8(g), uint8(b), uint8(a))
				c.Point(x+i, y+j)
			}
		}
	}
}

// DrawPartialSprite ... draws a rectangle from a sprite at x,y given offset ox and oy into sprite size w, h. No bounds checking
func (s *Sprite) DrawPartialSprite(c *Context, x, y, ox, oy, w, h float64) {
	for i := ox; i < ox+w; i++ {
		for j := oy; j < oy+h; j++ {
			r, g, b, a := s.At(int(i), int(j)).RGBA()
			// no blending for now!
			if a > 0 {
				c.Renderer.SetDrawColor(uint8(r), uint8(g), uint8(b), uint8(a))
				c.Point(x+i-ox, y+j-oy)
			}
		}
	}
}

// SampleSprite ... Samples from normal x, y of sprite
func (s *Sprite) SampleSprite(nx, ny float64) (rgba Colour) {
	bounds := s.Bounds()
	x := int(math.Trunc(nx * float64(bounds.Max.X-bounds.Min.X)))
	y := int(math.Trunc(ny * float64(bounds.Max.Y-bounds.Min.Y)))
	r, g, b, a := s.At(x, y).RGBA()
	rgba = NewColour(float64(r), float64(g), float64(b), float64(a))
	return
}

// SpriteSheet ...
type SpriteSheet struct {
	Sheet         *Sprite
	SpritesPerRow float64
	SpritesPerCol float64
	SpriteW       float64
	SpriteH       float64
}

// NewSpriteSheet ...
func NewSpriteSheet(filename string, NumPerCol, NumPerRow float64) (sh *SpriteSheet, err error) {
	sh = &SpriteSheet{SpritesPerRow: NumPerRow, SpritesPerCol: NumPerCol}
	sh.Sheet, err = NewSprite(filename)
	if err != nil {
		return nil, err
	}
	sh.SpriteW = sh.Sheet.W / sh.SpritesPerCol
	sh.SpriteH = sh.Sheet.H / sh.SpritesPerRow
	return
}

// DrawSpriteFromSheet ... given x and y coord draws the sprite from row x col of spritesheet
func (s *SpriteSheet) DrawSpriteFromSheet(c *Context, x, y, row, col float64) {
	col = math.Mod(col, s.SpritesPerCol)
	row = math.Mod(row, s.SpritesPerRow)
	ox := col * s.SpriteW
	oy := row * s.SpriteH
	s.Sheet.DrawPartialSprite(c, x, y, ox, oy, s.SpriteW, s.SpriteH)
}

// DrawSpriteFromSheetI ... indexes spritesheet as linear array.
func (s *SpriteSheet) DrawSpriteFromSheetI(c *Context, x, y, i float64) {
	i = math.Trunc(math.Mod(i, s.SpritesPerCol*s.SpritesPerRow))
	col := math.Mod(i, s.SpritesPerCol)
	row := math.Trunc(i / s.SpritesPerCol)
	ox := col * s.SpriteW
	oy := row * s.SpriteH
	s.Sheet.DrawPartialSprite(c, x, y, ox, oy, s.SpriteW, s.SpriteH)
}
