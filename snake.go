package gosnake

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
	DefaultConfig = GameConfig{
		XLen: 100,
		YLen: 100,
		XOff: 0,
		YOff: 0,
		BabySnake: Snake{
			Bodies: []Body{{48, 50}, {49, 50}, {50, 50}},
			Len:    3,
		},
		InitFood: Food{X: 52, Y: 52},
	}
}

type World struct {
	XLen, YLen       int
	XOffset, YOffset int
}

func (w World) CheckWallCollision(s Snake) (ok bool) {
	LX := w.XOffset
	RX := w.XLen + w.XOffset
	TY := w.YOffset
	BY := w.YLen + w.YOffset

	head := s.Bodies[s.Len-1]
	if head.X <= LX || head.X >= RX || head.Y >= BY || head.Y <= TY {
		return false
	}
	return true
}

type Body struct {
	X, Y int
}

type Snake struct {
	Bodies []Body
	Len    int
}

// assume snake is not too long
func (s Snake) isOnSnake(X, Y int) bool {
	for _, b := range s.Bodies {
		if b.X == X && b.Y == Y {
			return true
		}
	}
	return false
}

func (s Snake) CheckBodyCollision() (ok bool) {
	head := s.Bodies[s.Len-1]
	for i := 0; i < s.Len-2; i++ {
		body := s.Bodies[i]
		if body.X == head.X && body.Y == head.Y {
			return false
		}
	}
	return true
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (s *Snake) appendHead(b Body) {
	s.Bodies = append(s.Bodies, b)
}
func (s *Snake) shiftTail() {
	s.Bodies = s.Bodies[1:]
}
func (s *Snake) Walk(G *Game, direction Direction) {
	head := s.Bodies[s.Len-1]
	newHead := head
	switch direction {
	case Up:
		newHead.Y--
	case Down:
		newHead.Y++
	case Left:
		newHead.X--
	case Right:
		newHead.X++
	}
	s.appendHead(newHead)
	if G.Food.X == newHead.X && G.Food.Y == newHead.Y {
		s.Len++ // eat food
		G.RefreshFood()
		return
	}
	s.shiftTail()
}

type Food struct {
	X, Y int
}

type Game struct {
	World  World
	Food   Food
	Snake  Snake
	isOver bool
}

type GameConfig struct {
	XLen, YLen, XOff, YOff int
	BabySnake              Snake
	InitFood               Food
}

var DefaultConfig GameConfig

func GameInit(c GameConfig) *Game {
	world := World{
		XLen:    c.XLen,
		YLen:    c.YLen,
		XOffset: c.XOff,
		YOffset: c.YOff,
	}
	snake := c.BabySnake
	food := c.InitFood

	return &Game{
		World:  world,
		Snake:  snake,
		Food:   food,
		isOver: false,
	}
}
func (G *Game) RefreshFood() {
	X := rand.Intn(G.World.XLen-2) + G.World.XOffset + 1
	Y := rand.Intn(G.World.YLen-2) + G.World.YOffset + 1
	for G.Snake.isOnSnake(X, Y) {
		X = rand.Intn(G.World.XLen-2) + G.World.XOffset + 1
		Y = rand.Intn(G.World.YLen-2) + G.World.YOffset + 1
	}
	G.Food = Food{
		X: X,
		Y: Y,
	}
}
func (G *Game) CheckCollision() (ok bool) {
	return G.World.CheckWallCollision(G.Snake) && G.Snake.CheckBodyCollision()
}
func (G *Game) IsOver() bool { return G.isOver }
func (G *Game) WalkUp() {
	G.Snake.Walk(G, Up)
	if ok := G.CheckCollision(); !ok {
		G.isOver = true
	}
}
func (G *Game) WalkDown() {
	G.Snake.Walk(G, Down)
	if ok := G.CheckCollision(); !ok {
		G.isOver = true
	}
}
func (G *Game) WalkLeft() {
	G.Snake.Walk(G, Left)
	if ok := G.CheckCollision(); !ok {
		G.isOver = true
	}
}
func (G *Game) WalkRight() {
	G.Snake.Walk(G, Right)
	if ok := G.CheckCollision(); !ok {
		G.isOver = true
	}
}
