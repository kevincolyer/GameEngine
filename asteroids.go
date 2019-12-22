package main

import (
	"fmt"
	"os"
	"math"
	. "github.com/kevincolyer/GameEngine/GameEngine"
)

// helper function - can be passed in with GameEngine.New to modify the way blocks are drawn to the screen
func wrapScreen(x, y float64) (float64, float64) {
	return Wrap(x, 0, blocksw), Wrap(y, 0, blocksh)
}

var blocksw,blocksh,blocks float64
var WHITE, BLACK, RED, GREEN, BLUE, GREY50 Colour

func main() {
	blocksw=160
	blocksh=80
	blocks=8
 	WHITE  = NewColour(255,255,255,255)
 	BLACK  = NewColour(0,0,0,255      )
 	RED    = NewColour(255,0,0,255    )
 	GREEN  = NewColour(0,255,0,255    )
 	BLUE   = NewColour(0,0,255,255    )
 	GREY50 = NewColour(127,127,127,255)

	var ctx = New(blocks, blocksw, blocksh, "Asteroids", wrapScreen)

	onCreate(ctx)
	var running = true

	for running {
		running = onUpdate(ctx, ctx.Elapsed())
	}

	ctx.Destroy()
	os.Exit(0)
}

func  drawScore(c *Context, score int32) {
	t:=c.NewText(fmt.Sprintf("Score:%v",score), Colour{R: 255, G: 255, B: 255, A: 255})
	t.Draw(c.Renderer,0,0,0,0)
}

func drawFPS(c *Context, elapsed float64) {
	t:=c.NewText(fmt.Sprintf("FPS:%d",int(1000/elapsed)), Colour{R: 255, G: 255, B: 255, A: 255})
	t.Draw(c.Renderer,c.WinWidth-t.W,0,0,0)
}


type object struct {
	pos P2D
	vel V2D
	model []P2D
	w []P2D
	size float64
	angle float64
}

func (o object) draw(c *Context) {
	l:=len(o.w)
	if l==1 { c.Point(o.w[0].X,o.w[0].Y); return }
	for i:=1;i<l;i++ {
		c.Line(o.w[i-1].X, o.w[i-1].Y, o.w[i].X, o.w[i].Y)
	}
	c.Line(o.w[0].X,o.w[0].Y,o.w[l-1].X,o.w[l-1].Y)
}

// GAME GLOBAL VARIABLES
var ship object
var score int32
var worldSpeed float64
var bullets []object

func onCreate(c *Context) {
	fmt.Println("Created")
	c.Clear()
	c.Present()
	ship=object{
		pos: P2D{blocksw/2,blocksh/2}, 
		vel: V2D{0,0},
		model: []P2D{P2D{0,-1},P2D{-0.5,0.5},P2D{0.5,0.5}},
		w: []P2D{P2D{},P2D{},P2D{}},
		size: 6,
		angle: 0,
	}
	worldSpeed=0.1
}


func onUpdate(c *Context, elapsed float64) (running bool) {
	// boilerplate to start
    running, keys := c.PollQuitandKeys()
	if keys.Event {
		if keys.Key == "q" {
			running = false
		}
	}
	// Update code here...
	c.SetDrawColor(BLACK)
	c.Clear()
	// keys
	if keys.Key=="a" { ship.angle=ship.angle-1*elapsed}
	if keys.Key=="d" { ship.angle=ship.angle+1*elapsed}
	if keys.Key=="w" { 
		ship.vel.Dx= math.Sin(ship.angle)*elapsed*worldSpeed + ship.vel.Dx 
		ship.vel.Dy= -math.Cos(ship.angle)*elapsed*worldSpeed + ship.vel.Dy
	}
	if keys.Key=="x" { ship.vel.Dx=0; ship.vel.Dy=0}
	if keys.Key==" " {
		bull:=object{
			pos: P2D{ship.pos.X,ship.pos.Y},
			vel: V2D{math.Sin(ship.angle)*elapsed,-math.Cos(ship.angle)*elapsed},
		}
		bullets=append(bullets,bull)
		println(len(bullets))
	}
	// manipulations
		// ship
	ship.pos.X=Wrap(ship.pos.X+ship.vel.Dx*elapsed,0,blocksw)
	ship.pos.Y=Wrap(ship.pos.Y+ship.vel.Dy*elapsed,0,blocksh)
		// bullets
		for i:=range bullets {
			bullets[i].pos.X=Wrap(bullets[i].pos.X+bullets[i].vel.Dx,0,blocksw)
			bullets[i].pos.Y=Wrap(bullets[i].pos.Y+bullets[i].vel.Dy,0,blocksh)
		}

	// rotate
	for i := range ship.w {
		ship.w[i].X=math.Cos(ship.angle)*ship.model[i].X-math.Sin(ship.angle)*ship.model[i].Y
		ship.w[i].Y=math.Sin(ship.angle)*ship.model[i].X+math.Cos(ship.angle)*ship.model[i].Y
	}
	// scale
	for i := range ship.w {
		ship.w[i].X=ship.w[i].X*ship.size
		ship.w[i].Y=ship.w[i].Y*ship.size
	}
	// translate
	for i := range ship.w {
		ship.w[i].X+=ship.pos.X
		ship.w[i].Y+=ship.pos.Y
	}
	// draw
	c.SetDrawColor(WHITE)
	ship.draw(c)
	for i:=range bullets {
		c.Point(bullets[i].pos.X,bullets[i].pos.Y)
	}	
	
	
	drawScore(c,score)
	drawFPS(c,elapsed)
    // boilerplate to finish
    c.Present()
	Delay(1)
	return running
}
