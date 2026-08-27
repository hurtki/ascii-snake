package main

import (
	"fmt"
	"time"

	"github.com/hurtki/ascii-snake/internal/srv/app"
)

func main() {
	cfg := app.GameConfig{
		BaseSnakeLength: 5,
		InterestSize:    10,
		XSize:           200,
		YSize:           200,
		TickTime:        time.Millisecond * 33,
	}

	game := app.InitGame(cfg)

	go game.Start()

	snakeID, err := game.AddPlayer()
	if err != nil {
		fmt.Printf("Error adding a player: %v\n", err)
		return
	}

	go func() {
		time.Sleep(time.Second / 4)
		game.AddMove(snakeID, app.Move{Direction: app.Down})
	}()
	for {
		fmt.Println(game.GetInterestCellForSnake(snakeID))
	}
}
