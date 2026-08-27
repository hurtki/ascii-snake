package game_ui

import (
	"bufio"
	"fmt"
	"strings"
	"time"

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
	var centerX, centerY int
	var playerSnake *domain.Snake
	var controlComment string

	for _, s := range snakes {
		if s.PlayerID == ui.cfg.PlayerSnakeID {
			playerSnake = &s
			break
		}
	}

	switch {
	case playerSnake != nil:
		centerX = playerSnake.CordX
		centerY = playerSnake.CordY
		controlComment = "playing"
	case len(snakes) > 0:
		centerX = snakes[0].CordX
		centerY = snakes[0].CordY
		controlComment = "spectating"
	// case len(snakes) == 0:
	default:
		centerX = ui.cfg.XSize / 2
		centerY = ui.cfg.YSize / 2
		controlComment = "map center"
	}

	startX := centerX - ui.cfg.InterestSize
	endX := centerX + ui.cfg.InterestSize
	startY := centerY - ui.cfg.InterestSize
	endY := centerY + ui.cfg.InterestSize

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
	fmt.Fprintf(&sb, "X:%d, Y:%d | apples: %d | package: %d:%d:%d | %s exit: Q\r\n\r\n",
		centerX,
		centerY,
		len(apples),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		controlComment,
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
