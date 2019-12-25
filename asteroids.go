package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"

	. "github.com/kevincolyer/GameEngine/GameEngine"
)

// helper function - can be passed in with GameEngine.New to modify the way blocks are drawn to the screen
func wrapScreen(x, y float64) (float64, float64) {
	return Wrap(x, 0, blocksw), Wrap(y, 0, blocksh)
}

var blocksw, blocksh, blocks float64
var WHITE, BLACK, RED, GREEN, BLUE, GREY50 Colour

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		fmt.Println("writing cpu profile to ", *cpuprofile)
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}

	}
	defer pprof.StopCPUProfile()
	// end profiling setup code

	blocksw = 160
	blocksh = 80
	blocks = 8
	WHITE = NewColour(255, 255, 255, 255)
	BLACK = NewColour(0, 0, 0, 255)
	RED = NewColour(255, 0, 0, 255)
	GREEN = NewColour(0, 255, 0, 255)
	BLUE = NewColour(0, 0, 255, 255)
	GREY50 = NewColour(127, 127, 127, 255)

	var ctx = New(blocks, blocksw, blocksh, "Asteroids", wrapScreen)

	onCreate(ctx)
	var running = true

	for running {
		running = onUpdate(ctx, ctx.Elapsed())
	}

	// profiling code
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		fmt.Println("writing memory profile to ", *memprofile)
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
	ctx.Destroy()
	os.Exit(0)
}

func drawScore(c *Context, score int32) {
	t := c.NewText(fmt.Sprintf("Score:%v", score), Colour{R: 255, G: 255, B: 255, A: 255})
	t.Draw(c.Renderer, 0, 0, 0, 0)
}

func drawFPS(c *Context, elapsed float64) {
	t := c.NewText(fmt.Sprintf("FPS:%d", int(1000/elapsed)), Colour{R: 255, G: 255, B: 255, A: 255})
	t.Draw(c.Renderer, c.WinWidth-t.W, 0, 0, 0)
}

type object struct {
	pos    P2D
	vel    V2D
	model  []P2D
	w      []P2D
	size   float64
	angle  float64
	Da     float64
	health int
}

func (o object) draw(c *Context) {
	l := len(o.w)
	if l == 1 {
		c.Point(o.w[0].X, o.w[0].Y)
		return
	}
	for i := 1; i < l; i++ {
		c.Line(o.w[i-1].X, o.w[i-1].Y, o.w[i].X, o.w[i].Y)
	}
	c.Line(o.w[0].X, o.w[0].Y, o.w[l-1].X, o.w[l-1].Y)
}

func (o object) ScaleRotateTranslate() {
	// rotate
	for i := range o.w {
		o.w[i].X = math.Cos(o.angle)*o.model[i].X - math.Sin(o.angle)*o.model[i].Y
		o.w[i].Y = math.Sin(o.angle)*o.model[i].X + math.Cos(o.angle)*o.model[i].Y
	}
	// scale
	for i := range o.w {
		o.w[i].X = o.w[i].X * o.size
		o.w[i].Y = o.w[i].Y * o.size
	}
	// translate
	for i := range o.w {
		o.w[i].X += o.pos.X
		o.w[i].Y += o.pos.Y
	}
}

// GAME GLOBAL VARIABLES
var ship object
var score int32
var worldSpeed float64
var bulletSpeed float64
var maxSpeed float64
var bullets []object
var rocks []object

func onCreate(c *Context) {
	fmt.Println("Created")
	c.Clear()
	c.Present()
	resetGame()
	worldSpeed = 1
	bulletSpeed = worldSpeed * 0.1
	maxSpeed = math.Pow(2, 2)

}

func resetGame() {
	ship = object{
		pos:   P2D{blocksw / 2, blocksh / 2},
		vel:   V2D{0, 0},
		model: []P2D{P2D{0, -1}, P2D{-0.5, 0.5}, P2D{0.5, 0.5}},
		w:     []P2D{P2D{}, P2D{}, P2D{}},
		size:  6,
		angle: 0,
	}
	score = 0
	bullets = nil
	rocks = nil
	rocks = append(rocks, makeRock(blocksw/4, blocksh/2, 16), makeRock(blocksw*3/4, blocksh/2, 16))

}

func makeRock(x, y float64, size float64) (rock object) {
	rock = object{size: size, pos: P2D{X: x, Y: y}}
	for a := 0.0; a < 2*PI; a += 2 * PI / 20 {
		r := 0.6 + rand.Float64()*0.4
		rock.model = append(rock.model, P2D{math.Sin(a) * r, -math.Cos(a) * r})
	}
	rock.Da = rand.Float64()*PI/150 - PI/300
	rock.angle = rand.Float64() * PI * 2
	rock.w = append(rock.w, rock.model...)
	rock.vel = V2D{(rand.Float64() - 0.5) * math.Sin(rock.angle), -(rand.Float64() - 0.5) * math.Cos(rock.angle)}
	return
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
	// keys
	if keys.Key == "a" {
		ship.angle = ship.angle - 1*elapsed*worldSpeed
	}
	if keys.Key == "d" {
		ship.angle = ship.angle + 1*elapsed*worldSpeed
	}
	if keys.Key == "w" && (ship.vel.Dx*ship.vel.Dx+ship.vel.Dy*ship.vel.Dy) < maxSpeed {
		ship.vel.Dx = math.Sin(ship.angle)*elapsed*worldSpeed*0.5 + ship.vel.Dx
		ship.vel.Dy = -math.Cos(ship.angle)*elapsed*worldSpeed*0.5 + ship.vel.Dy
	}
	if keys.Key == "x" {
		ship.vel.Dx = 0
		ship.vel.Dy = 0
	}
	if keys.Key == " " {
		bull := object{
			pos: P2D{ship.pos.X, ship.pos.Y},
			vel: V2D{
				(math.Abs(ship.vel.Dx) + bulletSpeed) * math.Sin(ship.angle),
				(math.Abs(ship.vel.Dx) + bulletSpeed) * -math.Cos(ship.angle),
			},
			health: 1000,
		}
		bullets = append(bullets, bull)
	}
	// manipulations
	// ship
	ship.pos.X = Wrap(ship.pos.X+ship.vel.Dx*elapsed, 0, blocksw)
	ship.pos.Y = Wrap(ship.pos.Y+ship.vel.Dy*elapsed, 0, blocksh)
	// bullets
	for i := range bullets {
		if bullets[i].health > 0 {
			bullets[i].pos.X += bullets[i].vel.Dx
			bullets[i].pos.Y += bullets[i].vel.Dy
			bullets[i].health--
			if bullets[i].pos.X > blocksw || bullets[i].pos.X < 0 || bullets[i].pos.Y > blocksh || bullets[i].pos.Y < 0 {
				bullets[i].health = 0
			}
		}
	}
	if len(bullets) > 1 && bullets[0].health == 0 {
		bullets = bullets[1:]
	}
	for i := range rocks {
		rocks[i].angle = rocks[i].angle + rocks[i].Da*elapsed
		rocks[i].pos.X = Wrap(rocks[i].pos.X+rocks[i].vel.Dx*elapsed, 0, blocksw)
		rocks[i].pos.Y = Wrap(rocks[i].pos.Y+rocks[i].vel.Dy*elapsed, 0, blocksh)
		rocks[i].ScaleRotateTranslate()
	}

	// // rotate
	// for i := range ship.w {
	// 	ship.w[i].X=math.Cos(ship.angle)*ship.model[i].X-math.Sin(ship.angle)*ship.model[i].Y
	// 	ship.w[i].Y=math.Sin(ship.angle)*ship.model[i].X+math.Cos(ship.angle)*ship.model[i].Y
	// }
	// // scale
	// for i := range ship.w {
	// 	ship.w[i].X=ship.w[i].X*ship.size
	// 	ship.w[i].Y=ship.w[i].Y*ship.size
	// }
	// // translate
	// for i := range ship.w {
	// 	ship.w[i].X+=ship.pos.X
	// 	ship.w[i].Y+=ship.pos.Y
	// }
	ship.ScaleRotateTranslate()

	// draw
	c.SetDrawColor(GREY50)
	for _, r := range rocks {
		r.draw(c)
	}

	c.SetDrawColor(WHITE)
	ship.draw(c)

	for i := range bullets {
		if bullets[i].health > 0 {
			c.Point(bullets[i].pos.X, bullets[i].pos.Y)
		}
	}

	drawScore(c, score)
	//drawFPS(c,elapsed)
	// boilerplate to finish
	c.Present()
	Delay(1)
	return running
}
