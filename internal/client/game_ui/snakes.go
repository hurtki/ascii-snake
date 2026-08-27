package game_ui

import "github.com/hurtki/ascii-snake/internal/client/domain"

func (ui *GameUI) drawSnakes(draw drawFunc, snakes []domain.Snake) {
	for _, s := range snakes {
		headChar := 'E'
		bodyChar := 'e'
		if s.PlayerID == ui.cfg.PlayerSnakeID {
			headChar = '@'
			bodyChar = '#'
		}

		draw(s.CordX, s.CordY, headChar)

		x, y := s.CordX, s.CordY
		for _, dir := range s.Body {
			x, y = apply(x, y, dir)
			draw(x, y, bodyChar)
		}
	}

}

func apply(cordX, cordY int, dir domain.Direction) (int, int) {

	switch dir {
	case domain.Up:
		return cordX - 1, cordY
	case domain.Down:
		return cordX + 1, cordY
	case domain.Left:
		return cordX, cordY - 1
	case domain.Right:
		return cordX, cordY + 1
	}
	return cordX, cordY
}
