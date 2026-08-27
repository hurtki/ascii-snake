package game_ui

import "github.com/hurtki/ascii-snake/internal/client/domain"

func (ui *GameUI) drawApples(draw drawFunc, apples []domain.Apple) {
	for _, apple := range apples {
		draw(apple.CordX, apple.CordY, 'O')
	}
}
