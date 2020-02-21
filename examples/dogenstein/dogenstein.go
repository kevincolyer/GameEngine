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
var blocksi = flag.Int("blocks", 4, "Blocks of X pixels")

var x, y, z, angle float64
var comment string

var wall *Sprite
var lamp *Sprite
var ball *Sprite
var err error

const ww = 10
const wh = 20

var world [wh][ww]int

const horizon float64 = 15
const screenz float64 = 0.5 // for now...

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

type spriteObject struct {
	sprite *Sprite
	x, y   float64
	dx, dy float64
	alive  float64
}

var objects []spriteObject
var bullets []spriteObject
var depthbuffer *ZBuffer
var zbuff []float64

// GAME GLOBAL VARIABLES
const FOV float64 = PI / 2
const FOV2 float64 = FOV / 2

var worldSpeed float64
var distScreen = 0.5
var commentTicker = 0.0

// to help undistort the fishbowl effect
var rayAngles []float64
var blocksPerUnit = 10.0 // tunable parameter for changing the rayAngles

func onCreate(c *Context) {
	worldSpeed = 0.1
	wall, err = NewSprite("../../assets/wall.png")
	if err != nil {
		panic("couldn't load sprite")
	}
	lamp, err = NewSprite("../../assets/lamppost.png")
	if err != nil {
		panic("couldn't load sprite")
	}
	ball, err = NewSprite("../../assets/tennisball.png")
	if err != nil {
		panic("couldn't load sprite")
	}
	rayAngles = make([]float64, int(blocksw))
	// for i := range rayAngles {
	// 	yEye := distScreen
	// 	xEye := (float64(i) - (blocksw / 2.0)) / blocksPerUnit
	// 	rayAngles[i] = math.Atan2(yEye, xEye)
	// }

	objects = []spriteObject{
		spriteObject{sprite: lamp, x: 3, y: 4},
		spriteObject{sprite: lamp, x: 4, y: 8},
		spriteObject{sprite: lamp, x: 3, y: 12},
	}
	bullets = []spriteObject{}
	c.Clear()
	c.Present()
	resetGame()
	depthbuffer = NewZBuffer(blocksw, 1)
	zbuff = make([]float64, int(blocksw))
}

func resetGame() {
	x = 3
	y = 3
	angle = PI / 2
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
	if commentTicker <= 0.0 {
		comment = ""
		commentTicker = 0
	} else {
		commentTicker -= elapsed
	}
	// Update code here...
	c.SetDrawColor(BLACK)
	c.Clear()
	depthbuffer.Clear()
	// keys //////////////////////////////////////////
	if keys.Key == "a" {
		angle -= ews
		if angle < 0 {
			angle += PI * 2
		}
	}
	if keys.Key == "d" {
		angle += ews
		if angle >= PI*2 {
			angle -= PI * 2
		}
	}
	if keys.Key == "=" {
		blocksPerUnit += ews
	}
	if keys.Key == "-" {
		blocksPerUnit -= ews
		if blocksPerUnit < 0 {
			blocksPerUnit = 0.1
		}
	}

	nx := x
	ny := y
	// forward backward
	if keys.Key == "w" {
		nx = x + math.Sin(angle)*ews
		ny = y + math.Cos(angle)*ews
	}
	if keys.Key == "s" {
		nx = x - math.Sin(angle)*ews
		ny = y - math.Cos(angle)*ews
	}
	// strafe left / right
	if keys.Key == "n" {
		nx = x - math.Cos(angle)*ews
		ny = y + math.Sin(angle)*ews
	}
	if keys.Key == "m" {
		nx = x + math.Cos(angle)*ews
		ny = y - math.Sin(angle)*ews
	}
	nx = Clamp(nx, 0, float64(ww))
	ny = Clamp(ny, 0, float64(wh))

	if world[int(math.Floor(ny))][int(math.Floor(nx))] == 0 {
		x = nx
		y = ny
	} else {
		comment = "BUMP! Ooops"
		commentTicker = 20.0
	}

	if keys.Key == " " {
	}

	// manipulations /////////////////////////////////////

	// Draw
	///////////////////////////////////////////////////
	screenmid := blocksh / 2
	// for i := range rayAngles {
	// 	yEye := distScreen
	// 	xEye := (float64(i) - (blocksw / 2.0)) / blocksPerUnit
	// 	rayAngles[i] = math.Atan2(yEye, xEye)
	// }
	// screenmidw := blocksw / 2
	for bx := 0.0; bx < blocksw; bx++ {
		a := angle + bx/blocksw*FOV - FOV2
		// a := angle + rayAngles[int(bx)]
		// println(int(bx), angle, rayAngles[int(bx)])
		// march a ray from screen distance to horizon
		// z is distance to hitting a block. give up at horizon
		// z = math.Sqrt(screenz*screenz + (bx-screenmidw)*(bx-screenmidw))
		z = screenz
		var tx, ty float64
		hitWall := false
		for z < horizon {
			z += 0.01
			tx = x + math.Sin(a)*z
			ty = y + math.Cos(a)*z
			if tx >= ww || tx < 0 || ty >= wh || ty < 0 {
				z = horizon
				break
			}
			if world[int(ty)][int(tx)] == 1 {
				// hit a wall
				hitWall = true
				break
			}
		}
		// if hit wall check if hit corner...
		if hitWall {

		}

		// draw from top to bottom
		wallt := math.Trunc(screenmid - (screenmid / z))
		wallb := math.Trunc(blocksh - wallt)
		c.SetDrawColor(BLACK) // ceiling to walltop
		for by := 0.0; by < blocksh; by++ {
			// wall top to bottom
			if by >= wallt && by < wallb {
				//c.SetDrawColor(WHITE.Fade(1.5 / z))
				// texture draw
				// find normalised x and y to sample texture with
				ny := (by - wallt) / (wallb - wallt)
				var nx float64
				// workout which side of cube hit
				cx := math.Trunc(tx) + 0.5
				cy := math.Trunc(ty) + 0.5

				// line to eye
				le := ty - cy
				// line to ray point
				lr := tx - cx
				// angle between
				wangle := math.Atan2(le, lr)
				const PI4 = PI / 4
				const PI3_4 = 0.75 * PI
				// angle (less pi/4) as rotate axis a bit to indicate which side is hit
				// +pi side (clockwise 2 quadrants)
				if wangle >= -PI4 && wangle < PI4 {
					nx = ty - math.Trunc(ty)
				}
				if wangle >= PI4 && wangle < PI3_4 {
					nx = tx - math.Trunc(tx)
				}
				// -PI side (anti-clockwise quadrants)
				if wangle < -PI4 && wangle >= -PI3_4 {
					nx = tx - math.Trunc(tx)
				}
				if wangle >= PI3_4 || wangle < -PI3_4 {
					nx = ty - math.Trunc(ty)
				}
				c.SetDrawColor(wall.SampleSprite(nx, ny))
			}
			// wall bottom
			if by >= wallb {
				c.SetDrawColor(RED.Fade(1 - (blocksh-by)/screenmid))
			}
			c.Point(bx, by)
		}
		//depthbuffer.SetIfNearer(bx, 0, z)
		zbuff[int(bx)] = z
	}

	eyex := math.Sin(angle)
	eyey := math.Cos(angle)
	// Draw static objects - sprites
	// TODO draw sdynamic objects - bullets
	for _, o := range objects {

		// is object in field of view?
		oVecx := o.x - x
		oVecy := o.y - y
		z := math.Sqrt(oVecx*oVecx + oVecy*oVecy)
		oAngle := math.Atan2(eyey, eyex) - math.Atan2(oVecy, oVecx)
		if oAngle < -PI {
			oAngle += PI * 2
		}
		if oAngle > PI {
			oAngle -= PI * 2
		}
		if math.Abs(oAngle) < FOV2 && z >= 0.5 && z < horizon {
			// draw sprite
			oCeil := screenmid - (blocksh / z)
			oFloor := blocksh - oCeil
			oHeight := oFloor - oCeil
			oAspectRatio := o.sprite.H / o.sprite.W
			oWidth := oHeight / oAspectRatio
			oMidObject := (0.5*(oAngle/FOV2) + 0.5) * blocksw
			for lx := 0.0; lx < oWidth; lx++ {
				for ly := 0.0; ly < oHeight; ly++ {
					sampleX := lx / oWidth
					sampleY := ly / oHeight
					clr := o.sprite.SampleSprite(sampleX, sampleY)
					oCol := oMidObject + lx - (oWidth / 2)
					if oCol >= 0 && oCol < blocksw {
						if clr.A > 0 && zbuff[int(oCol)] >= z {
							c.SetDrawColor(clr)
							c.Point(oCol, ly+oCeil)
							zbuff[int(oCol)] = z
						}
					}
				}
			}
		}
	}
	// Draw text and 'top' layers
	c.SetDrawColor(STEELBLUE)
	if *fps {
		c.DrawText(1, 17, 4, fmt.Sprintf("fps:%d", int(100/elapsed)))
	}
	c.DrawText(0, 0, 4, fmt.Sprintf("x: %.2f y: %.2f a: %.2f blocksPerUnit %.2f    %v", x, y, angle, blocksPerUnit, comment))
	// boilerplate to finish
	c.Present()
	Delay(1)
	return running
}
