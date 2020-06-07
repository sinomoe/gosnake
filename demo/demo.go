package main

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/sinomoe/gosnake"
	"time"
)

func drawSnake(s gosnake.Snake) {
	for _, v := range s.Bodies {
		termbox.SetCell(v.X, v.Y, '●', termbox.ColorCyan, termbox.ColorBlack)
	}
}

func drawFood(f gosnake.Food) {
	termbox.SetCell(f.X, f.Y, '◉', termbox.ColorYellow, termbox.ColorBlack)
}

func drawWall(w gosnake.World) {
	xOff := 0
	yOff := 0
	xLen := w.XLen
	yLen := w.YLen
	termbox.SetCell(xOff, yOff, '+', termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	termbox.SetCell(xOff+xLen, yOff, '+', termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	termbox.SetCell(xOff, yOff+yLen, '+', termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	termbox.SetCell(xOff+xLen, yOff+yLen, '+', termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	for i := 1; i < xLen; i++ {
		termbox.SetCell(xOff+i, yOff, '-', termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
		termbox.SetCell(xOff+i, yOff+yLen, '-', termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	}
	for i := 1; i < yLen; i++ {
		termbox.SetCell(xOff, yOff+i, '|', termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
		termbox.SetCell(xOff+xLen, yOff+i, '|', termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	}
}

func updateAndDrawAll(G *gosnake.Game) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawWall(G.World)
	drawFood(G.World.Food)
	drawSnake(G.World.Snake)
	drawStatus(G)
	termbox.Flush()
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func drawStatus(G *gosnake.Game) {
	scoreMsg := fmt.Sprintf("Score: %d", G.World.Snake.Len-3)
	tbprint(41, 0, termbox.ColorBlack, termbox.ColorWhite, scoreMsg)
	if G.IsOver() {
		tbprint(41, 1, termbox.ColorBlack, termbox.ColorWhite, "game over")
	}
	head := G.World.Snake.Head()
	headMsg := fmt.Sprintf("Head: (%d, %d)", head.X, head.Y)
	tbprint(41, 2, termbox.ColorBlack, termbox.ColorWhite, headMsg)
	food := G.World.Food
	foodMsg := fmt.Sprintf("Food: (%d, %d)", food.X, food.Y)
	tbprint(41, 3, termbox.ColorBlack, termbox.ColorWhite, foodMsg)
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
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	updateAndDrawAll(G)

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
			updateAndDrawAll(G)
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
