package main

import (
	"container/list"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"

	. "github.com/kevincolyer/GameEngine/GameEngine"
)

// helper function - can be passed in with GameEngine.New to modify the way blocks are drawn to the screen
func wrapScreen(x, y float64) (float64, float64) {
	return Wrap(x, 0, blocksw), Wrap(y, 0, blocksh)
}

var blocksw, blocksh, blocks float64
var WHITE, BLACK, RED, GREEN, BLUE, GREY50, DARKRED, SADDLEBROWN, STEELBLUE Colour
var explodeShip bool

var fps = flag.Bool("fps", false, "Display Frames per second")
var blocksi = flag.Int("blocks", 8, "Blocks of X pixels")

func main() {
	flag.Parse()
	blocksw = 160
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
	var ctx = New(blocks, blocksw, blocksh, "Asteroids", wrapScreen)

	onCreate(ctx)
	var running = true

	for running {
		running = onUpdate(ctx, ctx.Elapsed())
	}

	ctx.Destroy()
	os.Exit(0)
}

type Object struct {
	Pos    P2D
	Vel    V2D
	Model  []P2D
	W      []P2D
	size   float64
	Angle  float64
	Da     float64
	Health int
}

func (o Object) Draw(c *Context) {
	l := len(o.W)
	if l == 1 {
		c.Point(o.W[0].X, o.W[0].Y)
		return
	}
	for i := 1; i < l; i++ {
		c.Line(o.W[i-1].X, o.W[i-1].Y, o.W[i].X, o.W[i].Y)
	}
	c.Line(o.W[0].X, o.W[0].Y, o.W[l-1].X, o.W[l-1].Y)
}

func (o Object) ScaleRotateTranslate() {
	// rotate
	for i := range o.W {
		o.W[i].X = math.Cos(o.Angle)*o.Model[i].X - math.Sin(o.Angle)*o.Model[i].Y
		o.W[i].Y = math.Sin(o.Angle)*o.Model[i].X + math.Cos(o.Angle)*o.Model[i].Y
	}
	// scale
	for i := range o.W {
		o.W[i].X = o.W[i].X * o.size
		o.W[i].Y = o.W[i].Y * o.size
	}
	// translate
	for i := range o.W {
		o.W[i].X += o.Pos.X
		o.W[i].Y += o.Pos.Y
	}
}

// GAME GLOBAL VARIABLES
var ship Object
var score int32
var hiscore int32 = 0
var worldSpeed float64
var bulletSpeed float64
var maxSpeed float64
var bullets *list.List
var rocks *list.List
var explosion [24]*Object

func onCreate(c *Context) {
	worldSpeed = 1
	bulletSpeed = worldSpeed * 0.1
	maxSpeed = math.Pow(2, 2)
	// must be before bullets.Init
	bullets = list.New()
	rocks = list.New()

	c.Clear()
	c.Present()
	resetGame()
}

func resetGame() {
	ship = Object{
		Pos:   P2D{blocksw / 2, blocksh / 2},
		Vel:   V2D{0, 0},
		Model: []P2D{P2D{0, -1}, P2D{-0.5, 0.5}, P2D{0.5, 0.5}},
		W:     []P2D{P2D{}, P2D{}, P2D{}},
		size:  6,
		Angle: 0,
	}
	if score > hiscore {
		hiscore = score
	}
	score = 0
	explodeShip = false
	bullets.Init()
	rocks.Init()
	rocks.PushBack(makeRock(blocksw/4, blocksh/2, 16))
	rocks.PushBack(makeRock(blocksw*3/4, blocksh/2, 16))
}

func makeRock(x, y float64, size float64) (rock *Object) {
	rock = &Object{size: size, Pos: P2D{X: x, Y: y}, Health: 1}
	for a := 0.0; a < 2*PI; a += 2 * PI / 20 {
		r := 0.6 + rand.Float64()*0.4
		rock.Model = append(rock.Model, P2D{math.Sin(a) * r, -math.Cos(a) * r})
	}
	rock.Da = rand.Float64()*PI/150 - PI/300
	rock.Angle = rand.Float64() * PI * 2
	rock.W = append(rock.W, rock.Model...)
	rock.Vel = V2D{(rand.Float64() - 0.5) * math.Sin(rock.Angle), -(rand.Float64() - 0.5) * math.Cos(rock.Angle)}
	return
}

func drawExplosion(c *Context, elapsed float64) {
	if explosion[0].size < 0 {
		resetGame()
		return
	}
	for _, j := range explosion {
		col := WHITE
		if rand.ExpFloat64() > 0.5 {
			col = RED
		}
		col = col.Fade(j.size / 25.0)
		// update
		j.Pos.X += j.Vel.Dx * elapsed * worldSpeed * 0.05
		j.Pos.Y += j.Vel.Dy * elapsed * worldSpeed * 0.05
		j.size -= elapsed * worldSpeed * .1
		c.SetDrawColor(col)
		c.Point(j.Pos.X, j.Pos.Y)
	}
}

func makeExplosion() {
	k := 0
	for i := 0; i < 3; i++ {
		x0 := ship.W[i].X
		y0 := ship.W[i].Y
		x1 := ship.W[(i+1)%3].X
		y1 := ship.W[(i+1)%3].Y
		dx := (x0 - x1) / 9
		dy := (y0 - y1) / 9
		// split into 8 points
		for j := 0.0; j < 8; j++ {
			x := x0 + dx*j
			y := y0 + dy*j
			theta := math.Atan2(x-ship.Pos.X, y-ship.Pos.Y) + rand.Float64()*PI/3 - PI/6

			explosion[k] = &Object{
				size: 25,
				Pos:  P2D{x, y},
				Vel: V2D{
					ship.Vel.Dx + math.Sin(theta)*4,
					ship.Vel.Dy + math.Cos(theta)*4,
				},
			}
			k++
		}
	}
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
		ship.Angle = ship.Angle - 1*elapsed*worldSpeed
	}
	if keys.Key == "d" {
		ship.Angle = ship.Angle + 1*elapsed*worldSpeed
	}
	if keys.Key == "w" && (ship.Vel.Dx*ship.Vel.Dx+ship.Vel.Dy*ship.Vel.Dy) < maxSpeed {
		ship.Vel.Dx = math.Sin(ship.Angle)*elapsed*worldSpeed*0.5 + ship.Vel.Dx
		ship.Vel.Dy = -math.Cos(ship.Angle)*elapsed*worldSpeed*0.5 + ship.Vel.Dy
	}
	if keys.Key == "x" {
		ship.Vel.Dx = 0
		ship.Vel.Dy = 0
	}
	if keys.Key == " " && explodeShip == false {
		bull := &Object{
			Pos: P2D{ship.Pos.X, ship.Pos.Y},
			Vel: V2D{
				(math.Abs(ship.Vel.Dx) + bulletSpeed) * math.Sin(ship.Angle),
				(math.Abs(ship.Vel.Dx) + bulletSpeed) * -math.Cos(ship.Angle),
			},
			Health: 1000,
		}
		bullets.PushFront(bull)
	}
	// manipulations /////////////////////////////////////
	// ship
	ship.Pos.X = Wrap(ship.Pos.X+ship.Vel.Dx*elapsed, 0, blocksw)
	ship.Pos.Y = Wrap(ship.Pos.Y+ship.Vel.Dy*elapsed, 0, blocksw)

	// bullets(
	for b := bullets.Front(); b != nil; b = b.Next() {
		v := b.Value.(*Object)
		if v.Health > 0 {
			v.Pos.X += v.Vel.Dx
			v.Pos.Y += v.Vel.Dy
			v.Health--
			if v.Pos.X > blocksw || v.Pos.X < 0 || v.Pos.Y > blocksh || v.Pos.Y < 0 {
				v.Health = 0
				bullets.Remove(b)
			}
		}
	}

	// rocks
	for r := rocks.Front(); r != nil; r = r.Next() {
		rock := r.Value.(*Object)
		if rock.Health > 0 {

			rock.Angle = rock.Angle + rock.Da*elapsed
			rock.Pos.X = Wrap(rock.Pos.X+rock.Vel.Dx*elapsed, 0, blocksw)
			rock.Pos.Y = Wrap(rock.Pos.Y+rock.Vel.Dy*elapsed, 0, blocksh)
			rock.ScaleRotateTranslate()

			// collision detection
			// ship
			rockx := rock.Pos.X
			rocky := rock.Pos.Y
			dx := rockx - ship.Pos.X
			dy := rocky - ship.Pos.Y
			if dx*dx+dy*dy < (ship.size+rock.size)*(ship.size+rock.size) && explodeShip != true {
				explodeShip = true
				makeExplosion()
			}

			// bullets
			for b := bullets.Front(); b != nil; b = b.Next() {
				v := b.Value.(*Object)
				if v.Health == 0 {
					continue
				}
				dx = rockx - v.Pos.X
				dy = rocky - v.Pos.Y
				if dx*dx+dy*dy < rock.size*rock.size {
					// hit!
					// remove bullet, rock and increment score
					rock.Health = 0
					score += (16 - int32(rock.size)) * 10
					v.Health = 0
					if rock.size > 4 {
						// make two more rocks
						rocks.PushFront(makeRock(rock.Pos.X+6, rock.Pos.Y-6, rock.size/2))
						rocks.PushFront(makeRock(rock.Pos.X-6, rock.Pos.Y+6, rock.size/2))
					}
					rocks.Remove(r)
				}
			} // end bullets loop
		}
	} // end rocks loop

	// add two rocks if all rocks destroyed
	if rocks.Len() == 0 {
		score += 1000
		rocks.PushFront(makeRock(Wrap(ship.Pos.X+blocksw/2, 0, blocksw), rand.Float64()*blocksh, 16))
		rocks.PushFront(makeRock(Wrap(ship.Pos.X-blocksw/2, 0, blocksw), rand.Float64()*blocksh, 16))
	}

	// rotate scale and translate
	ship.ScaleRotateTranslate()

	// Draw
	///////////////////////////////////////////////////
	c.SetDrawColor(SADDLEBROWN)
	for r := rocks.Front(); r != nil; r = r.Next() {
		rock := r.Value.(*Object)
		if rock.Health > 0 {
			rock.Draw(c)
		}
	}
	if explodeShip == false {
		c.SetDrawColor(WHITE)
		ship.Draw(c)
	} else {
		drawExplosion(c, elapsed)
	}

	c.SetDrawColor(STEELBLUE)
	for b := bullets.Front(); b != nil; b = b.Next() {
		v := b.Value.(*Object)
		if v.Health > 0 {
			c.Point(v.Pos.X, v.Pos.Y)
		}
	}

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
