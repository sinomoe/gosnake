package gosnake

import "testing"

func TestBodyCollision(t *testing.T) {
	G := GameInit(DefaultConfig)
	if G.IsOver() {
		t.Error("step 0 error collision")
		return
	}
	G.WalkUp()
	if G.IsOver() {
		t.Error("step 1 error collision")
		return
	}
	G.WalkLeft()
	if G.IsOver() {
		t.Error("step 2 error collision")
		return
	}
	G.WalkDown()
	if !G.IsOver() {
		t.Error("step 3 should collision")
		return
	}
}

func TestWallCollision(t *testing.T) {
	G := GameInit(DefaultConfig)
	for i := 0; i < 49; i++ {
		G.WalkRight()
		if G.IsOver() {
			t.Errorf("step %d shodld not collision", i+1)
			return
		}
	}
	G.WalkRight()
	if !G.IsOver() {
		t.Error("step 50 should collision")
		return
	}
}

func TestEatFood(t *testing.T) {
	G := GameInit(DefaultConfig)
	initLen := DefaultConfig.BabySnake.Len
	initFoodX := DefaultConfig.InitFood.X
	initFoodY := DefaultConfig.InitFood.Y
	for i := 0; i < 2; i++ {
		G.WalkRight()
		if G.IsOver() {
			t.Errorf("step %d should not game over", i+1)
			return
		}
		if G.World.Snake.Len != initLen {
			t.Errorf("step %d: error snake Len(%d)", i+1, G.World.Snake.Len)
			return
		}
	}
	G.WalkDown()
	if G.IsOver() {
		t.Errorf("step %d shodld not over", 3)
		return
	}
	if G.World.Snake.Len != initLen {
		t.Errorf("step %d: error snake Len(%d)", 3, G.World.Snake.Len)
		return
	}
	G.WalkDown()
	if G.IsOver() {
		t.Errorf("step %d shodld not over", 4)
		return
	}
	if G.World.Snake.Len == initLen {
		t.Errorf("step %d: snake Len should be %d", 4, G.World.Snake.Len+1)
		return
	}
	if G.Score() == 0 {
		t.Errorf("step %d: score should be %d", 4, 1)
		return
	}

	// check new food target
	if G.World.Food.X == initFoodX && G.World.Food.Y == initFoodY {
		t.Error("new target not yet generated")
		return
	}
}
