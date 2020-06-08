package main

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/sinomoe/gosnake"
	"time"
)

type color termbox.Attribute

const (
	cyan      = color(termbox.ColorCyan)
	black     = color(termbox.ColorBlack)
	yellow    = color(termbox.ColorYellow)
	white     = color(termbox.ColorWhite)
	bold      = color(termbox.AttrBold)
	whiteBold = color(termbox.ColorWhite) | bold
)

func drawCell(x, y int, fg, bg color, ch rune) {
	termbox.SetCell(x, y, ch, termbox.Attribute(fg), termbox.Attribute(bg))
}

func tbprint(x, y int, fg, bg color, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, termbox.Attribute(fg), termbox.Attribute(bg))
		x += runewidth.RuneWidth(c)
	}
}

func clearScreen() {
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}
}

func render() {
	if err := termbox.Flush(); err != nil {
		panic(err)
	}
}

type snake struct {
	head, body rune
}

type wall struct {
	corners                  [4]rune
	top, bottom, left, right rune
}

type style struct {
	snake snake
	food  rune
	wall  wall
}

type gameBox struct {
	game       *gosnake.Game
	xOff, yOff int
	style      style
}

func (gb *gameBox) drawSnake() {
	for _, v := range gb.game.World.Snake.Bodies {
		drawCell(v.X, v.Y, cyan, black, gb.style.snake.body)
	}
}

func (gb *gameBox) drawFood() {
	f := gb.game.World.Food
	drawCell(f.X, f.Y, yellow, black, gb.style.food)
}

func (gb *gameBox) drawWall() {
	xOff := gb.xOff
	yOff := gb.yOff
	xLen := gb.game.World.XLen
	yLen := gb.game.World.YLen
	drawCell(xOff, yOff, whiteBold, black, gb.style.wall.corners[0])
	drawCell(xOff+xLen, yOff, whiteBold, black, gb.style.wall.corners[1])
	drawCell(xOff, yOff+yLen, whiteBold, black, gb.style.wall.corners[2])
	drawCell(xOff+xLen, yOff+yLen, whiteBold, black, gb.style.wall.corners[3])
	for i := 1; i < xLen; i++ {
		drawCell(xOff+i, yOff, whiteBold, black, gb.style.wall.top)
		drawCell(xOff+i, yOff+yLen, whiteBold, black, gb.style.wall.bottom)
	}
	for i := 1; i < yLen; i++ {
		drawCell(xOff, yOff+i, whiteBold, black, gb.style.wall.left)
		drawCell(xOff+xLen, yOff+i, whiteBold, black, gb.style.wall.right)
	}
}

func (gb *gameBox) updateAndDrawAll() {
	clearScreen()
	gb.drawWall()
	gb.drawFood()
	gb.drawSnake()
	gb.drawStatus()
	render()
}

func (gb *gameBox) drawStatus() {
	G := gb.game
	scoreMsg := fmt.Sprintf("Score: %d", G.Score())
	tbprint(41, 0, black, white, scoreMsg)
	if G.IsOver() {
		tbprint(41, 1, black, white, "game over")
	}
	head := G.World.Snake.Head()
	headMsg := fmt.Sprintf("Head: (%d, %d)", head.X, head.Y)
	tbprint(41, 2, black, white, headMsg)
	food := G.World.Food
	foodMsg := fmt.Sprintf("Food: (%d, %d)", food.X, food.Y)
	tbprint(41, 3, black, white, foodMsg)
}

func main() {
	G := gosnake.GameInit(gosnake.GameConfig{
		XLen: 40,
		YLen: 20,
		BabySnake: gosnake.Snake{
			Bodies: []gosnake.Body{{18, 12}, {19, 12}, {20, 12}},
			Len:    3,
		},
		InitFood:      gosnake.Food{X: 22, Y: 16},
		WallGenerator: gosnake.DefaultWallGenerator,
	})
	gb := &gameBox{
		game: G,
		xOff: 0,
		yOff: 0,
		style: style{
			snake: snake{head: '●', body: '●'},
			food:  '◉',
			wall:  wall{corners: [4]rune{'+', '+', '+', '+'}, top: '-', bottom: '-', left: '|', right: '|'},
		},
	}
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	gb.updateAndDrawAll()

	currentKey := termbox.KeyArrowRight

	go func() {
		for {
			select {
			case <-time.After(200 * time.Millisecond):
				if G.IsOver() {
					return
				}
				switch currentKey {
				case termbox.KeyEsc:
					return
				case termbox.KeyArrowUp:
					G.WalkUp()
				case termbox.KeyArrowDown:
					G.WalkDown()
				case termbox.KeyArrowLeft:
					G.WalkLeft()
				case termbox.KeyArrowRight:
					G.WalkRight()
				}
			}
			gb.updateAndDrawAll()
		}
	}()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return
			default:
				currentKey = ev.Key
			}
		}
	}
}
