package game_ui

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/hurtki/ascii-snake/internal/client/domain"
)

type GameUICfg struct {
	XSize, YSize  int
	InterestSize  int
	PlayerSnakeID int
}

type GameUI struct {
	out *bufio.Writer
	cfg GameUICfg
}

func NewGameUI(out *bufio.Writer, cfg GameUICfg) *GameUI {
	return &GameUI{
		out: out,
		cfg: cfg,
	}
}

type drawFunc func(int, int, rune)

func (ui *GameUI) DrawScreen(snakes []domain.Snake, apples []domain.Apple) {
	viewRadiusX := ui.cfg.InterestSize
	viewRadiusY := ui.cfg.InterestSize

	var playerSnake *domain.Snake
	for _, s := range snakes {
		if s.PlayerID == ui.cfg.PlayerSnakeID {
			playerSnake = &s
			break
		}
	}
	var startX, endX, startY, endY int
	if playerSnake == nil && len(snakes) > 0 {
		playerSnake = &snakes[0]
	}

	if playerSnake != nil {
		startX = playerSnake.CordX - viewRadiusX
		endX = playerSnake.CordX + viewRadiusX
		startY = playerSnake.CordY - viewRadiusY
		endY = playerSnake.CordY + viewRadiusY
	}

	rows := endX - startX + 1
	cols := endY - startY + 1

	screen := make([][]rune, rows)
	for i := range screen {
		screen[i] = make([]rune, cols)
		for j := range screen[i] {
			screen[i][j] = '.'
		}
	}
	draw := func(x, y int, char rune) {
		if x >= startX && x <= endX && y >= startY && y <= endY {
			screen[x-startX][y-startY] = char
		}
	}

	ui.drawBorders(draw)

	ui.drawApples(draw, apples)

	ui.drawSnakes(draw, snakes)

	var sb strings.Builder

	sb.WriteString("\033[H")
	fmt.Fprintf(&sb, "cord: X:%d, Y:%d | apples from srv: %d exit: Q\r\n\r\n",
		playerSnake.CordX,
		playerSnake.CordY,
		len(apples),
	)

	for i := range rows {
		for j := range cols {
			sb.WriteRune(screen[i][j])
			sb.WriteByte(' ')
		}
		sb.WriteString("\r\n")
	}

	ui.out.WriteString(sb.String())
	ui.out.Flush()
}
