package gosnake

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	DefaultWallGenerator WallGenerator
	DefaultConfig        GameConfig
)

func init() {
	rand.Seed(time.Now().Unix())
	DefaultWallGenerator = func(w World) Wall {
		return func(c Coordinates) bool {
			if c.x <= 0 || c.y <= 0 || c.x >= w.XLen || c.y >= w.YLen {
				return true
			}
			return false
		}
	}
	DefaultConfig = GameConfig{
		XLen: 100,
		YLen: 100,
		BabySnake: Snake{
			Bodies: []Body{{46, 50}, {47, 50}, {48, 50}, {49, 50}, {50, 50}},
			Len:    5,
		},
		InitFood:      Food{X: 52, Y: 52},
		WallGenerator: DefaultWallGenerator,
	}
}

type Coordinates struct {
	x, y int
}

func (c Coordinates) String() string {
	return fmt.Sprintf("(%d, %d)", c.x, c.y)
}

func (c Coordinates) Equal(co Coordinates) bool {
	return c.x == co.x && c.y == co.y
}

type direction int

const (
	up direction = iota + 1
	down
	left
	right
)

type Body struct {
	X, Y int
}

func (b Body) Coordinates() Coordinates {
	return Coordinates{
		x: b.X,
		y: b.Y,
	}
}

type Bodies []Body

func (bs Bodies) head() Body {
	return bs[len(bs)-1]
}
func (bs Bodies) any(f func(b Body) bool) (existed bool) {
	for _, b := range bs {
		if exist := f(b); exist {
			return true
		}
	}
	return false
}
func (bs Bodies) match(f func(b Body) bool) (matched int) {
	matched = 0
	for _, b := range bs {
		if exist := f(b); exist {
			matched++
		}
	}
	return matched
}
func (bs *Bodies) appendHead(b Body) {
	*bs = append(*bs, b)
}
func (bs *Bodies) shiftTail() {
	*bs = (*bs)[1:]
}

type Snake struct {
	Bodies    Bodies
	Len       int
	Direction direction
	score     int
}

func (s Snake) Head() Body {
	return s.Bodies.head()
}

// assume snake is not too long
func (s Snake) onSnake(c Coordinates) bool {
	existed := s.Bodies.any(func(b Body) bool {
		return b.Coordinates().Equal(c)
	})
	return existed
}
func (s Snake) detectBodyCollision() (detected bool) {
	head := s.Head()
	matched := s.Bodies.match(func(b Body) bool {
		return b.Coordinates().Equal(head.Coordinates())
	})
	if matched >= 2 {
		return true
	}
	return false
}

func (s *Snake) appendHead(b Body) {
	s.Bodies.appendHead(b)
}
func (s *Snake) shiftTail() {
	s.Bodies.shiftTail()
}
func (s *Snake) walk(w *World, direction direction) {
	s.Direction = direction
	head := s.Head()
	newHead := head
	switch direction {
	case up:
		newHead.Y--
	case down:
		newHead.Y++
	case left:
		newHead.X--
	case right:
		newHead.X++
	}
	s.appendHead(newHead)
	if w.Food.Coordinates().Equal(newHead.Coordinates()) {
		s.Len++ // eat food
		s.score++
		w.RefreshFood()
		return
	}
	s.shiftTail()
}

type Food struct {
	X, Y int
}

func (f Food) Coordinates() Coordinates {
	return Coordinates{
		x: f.X,
		y: f.Y,
	}
}

type WallGenerator func(w World) Wall
type Wall func(c Coordinates) bool

func (wa Wall) onWall(c Coordinates) bool {
	return wa(c)
}

type World struct {
	XLen, YLen int
	Snake      Snake
	Food       Food
	wall       Wall
}

func (w World) detectWallCollision() (detected bool) {
	head := w.Snake.Head()
	return w.wall.onWall(head.Coordinates())
}
func (w World) detectCollision() (detected bool) {
	return w.detectWallCollision() || w.Snake.detectBodyCollision()
}
func (w *World) RefreshFood() {
	X := rand.Intn(w.XLen - 1)
	Y := rand.Intn(w.YLen - 1)
	for w.Snake.onSnake(Coordinates{X, Y}) || w.wall.onWall(Coordinates{X, Y}) {
		X = rand.Intn(w.XLen - 1)
		Y = rand.Intn(w.YLen - 1)
	}
	w.Food = Food{
		X: X,
		Y: Y,
	}
}

type Game struct {
	World  World
	isOver bool
}

func (G Game) Score() int   { return G.World.Snake.score }
func (G Game) IsOver() bool { return G.isOver }
func (G *Game) WalkUp() {
	G.World.Snake.walk(&G.World, up)
	if detected := G.World.detectCollision(); detected {
		G.isOver = true
	}
}
func (G *Game) WalkDown() {
	G.World.Snake.walk(&G.World, down)
	if detected := G.World.detectCollision(); detected {
		G.isOver = true
	}
}
func (G *Game) WalkLeft() {
	G.World.Snake.walk(&G.World, left)
	if detected := G.World.detectCollision(); detected {
		G.isOver = true
	}
}
func (G *Game) WalkRight() {
	G.World.Snake.walk(&G.World, right)
	if detected := G.World.detectCollision(); detected {
		G.isOver = true
	}
}

type GameConfig struct {
	XLen, YLen    int
	BabySnake     Snake
	InitFood      Food
	WallGenerator WallGenerator
}

func GameInit(c GameConfig) *Game {
	world := World{
		XLen:  c.XLen,
		YLen:  c.YLen,
		Snake: c.BabySnake,
		Food:  c.InitFood,
	}
	world.wall = c.WallGenerator(world)

	return &Game{
		World:  world,
		isOver: false,
	}
}
