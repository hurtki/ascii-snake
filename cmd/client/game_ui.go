package main

import (
	"bufio"
	"fmt"
	"strings"
)

type Snake struct {
	PlayerID     int
	CordX, CordY int
	Body         []Direction
}

type Apple struct {
	CordX, CordY int
}

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

func (ui *GameUI) drawScreen(snakes []Snake, apples []Apple) bool {
	// playerSnake, ok := cell.Snakes[playerID]
	// if !ok {
	// 	out.WriteString("\033[H\033[2J")
	// 	out.WriteString("you died, no snake with your id found\r\n")
	// 	out.Flush()
	// 	time.Sleep(2 * time.Second)
	// 	return false
	// }

	viewRadiusX := ui.cfg.InterestSize
	viewRadiusY := ui.cfg.InterestSize

	var playerSnake *Snake
	for _, s := range snakes {
		if s.PlayerID == ui.cfg.PlayerSnakeID {
			playerSnake = &s
			break
		}
	}
	if playerSnake == nil {
		// snake not on the screen
	}

	startX := playerSnake.CordX - viewRadiusX
	endX := playerSnake.CordX + viewRadiusX
	startY := playerSnake.CordY - viewRadiusY
	endY := playerSnake.CordY + viewRadiusY

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

	for x := range ui.cfg.XSize {
		draw(x, -1, '║')
		draw(x, ui.cfg.YSize, '║')
	}
	for y := range ui.cfg.YSize {
		draw(-1, y, '═')
		draw(ui.cfg.XSize, y, '═')
	}

	draw(-1, -1, '╔')
	draw(ui.cfg.XSize, -1, '╚')
	draw(-1, ui.cfg.YSize, '╗')
	draw(ui.cfg.XSize, ui.cfg.YSize, '╝')

	for _, apple := range apples {
		draw(apple.CordX, apple.CordY, 'O')
	}

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

	var sb strings.Builder
	sb.WriteString("\033[H")
	fmt.Fprintf(&sb, "cord: X:%d, Y:%d | exit: Q\r\n\r\n", playerSnake.CordX, playerSnake.CordY)

	for i := range rows {
		for j := range cols {
			sb.WriteRune(screen[i][j])
			sb.WriteByte(' ')
		}
		sb.WriteString("\r\n")
	}

	ui.out.WriteString(sb.String())
	ui.out.Flush()

	return true
}

func apply(cordX, cordY int, dir Direction) (int, int) {

	switch dir {
	case Up:
		return cordX - 1, cordY
	case Down:
		return cordX + 1, cordY
	case Left:
		return cordX, cordY - 1
	case Right:
		return cordX, cordY + 1
	}
	return cordX, cordY
}
