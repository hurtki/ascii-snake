package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hurtki/ascii-snake/internal/srv/app"
	"golang.org/x/term"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	cfg := app.GameConfig{
		BaseSnakeLength:              10,
		InterestSize:                 15,
		XSize:                        100,
		YSize:                        100,
		TickTime:                     time.Millisecond * 100,
		PlayerSpawnPaddingFromBorder: 1,
	}
	game := app.InitGame(cfg)

	go game.Start()

	playerID, err := game.AddPlayer()
	if err != nil {
		fmt.Printf("error adding a player: %v\r\n", err)
		return
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err == nil {
		defer func() {
			_ = term.Restore(fd, oldState)
		}()
	}

	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h\033[2J\033[H")

	go handleInput(ctx, cancel, game, playerID)

	ticker := time.NewTicker(cfg.TickTime)
	defer ticker.Stop()

	out := bufio.NewWriter(os.Stdout)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !drawScreen(out, game, playerID, cfg) {
				return
			}
		}
	}
}

func handleInput(ctx context.Context, cancel context.CancelFunc, game *app.Game, playerID int) {
	buf := make([]byte, 3)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			if n == 1 {
				switch buf[0] {
				case 'w', 'W':
					game.AddMove(playerID, app.Move{Direction: app.Up})
				case 's', 'S':
					game.AddMove(playerID, app.Move{Direction: app.Down})
				case 'a', 'A':
					game.AddMove(playerID, app.Move{Direction: app.Left})
				case 'd', 'D':
					game.AddMove(playerID, app.Move{Direction: app.Right})
				case 'q', 'Q', 3:
					cancel()
					return
				}
			} else if n == 3 && buf[0] == 27 && buf[1] == 91 {
				switch buf[2] {
				case 'A':
					game.AddMove(playerID, app.Move{Direction: app.Up})
				case 'B':
					game.AddMove(playerID, app.Move{Direction: app.Down})
				case 'D':
					game.AddMove(playerID, app.Move{Direction: app.Left})
				case 'C':
					game.AddMove(playerID, app.Move{Direction: app.Right})
				}
			}
		}
	}
}

func drawScreen(out *bufio.Writer, g *app.Game, playerID int, cfg app.GameConfig) bool {
	cell := g.GetInterestCellForSnake(playerID)

	playerSnake, ok := cell.Snakes[playerID]
	if !ok {
		out.WriteString("\033[H\033[2J")
		out.WriteString("you died, no snake with your id found\r\n")
		out.Flush()
		time.Sleep(2 * time.Second)
		return false
	}

	viewRadiusX := cfg.InterestSize
	viewRadiusY := cfg.InterestSize

	startX := playerSnake.Cord.X - viewRadiusX
	endX := playerSnake.Cord.X + viewRadiusX
	startY := playerSnake.Cord.Y - viewRadiusY
	endY := playerSnake.Cord.Y + viewRadiusY

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

	for x := range cfg.XSize {
		draw(x, -1, '║')
		draw(x, cfg.YSize, '║')
	}
	for y := range cfg.YSize {
		draw(-1, y, '═')
		draw(cfg.XSize, y, '═')
	}

	draw(-1, -1, '╔')
	draw(cfg.XSize, -1, '╚')
	draw(-1, cfg.YSize, '╗')
	draw(cfg.XSize, cfg.YSize, '╝')

	for appCord := range cell.Apples {
		draw(appCord.X, appCord.Y, 'O')
	}

	for id, s := range cell.Snakes {
		headChar := 'E'
		bodyChar := 'e'
		if id == playerID {
			headChar = '@'
			bodyChar = '#'
		}

		draw(s.Cord.X, s.Cord.Y, headChar)

		currPos := s.Cord
		for _, dir := range s.Body {
			currPos = currPos.Apply(dir)
			draw(currPos.X, currPos.Y, bodyChar)
		}
	}

	var sb strings.Builder
	sb.WriteString("\033[H")
	fmt.Fprintf(&sb, "cord: X:%d, Y:%d | exit: Q\r\n\r\n", playerSnake.Cord.X, playerSnake.Cord.Y)

	for i := range rows {
		for j := range cols {
			sb.WriteRune(screen[i][j])
			sb.WriteByte(' ')
		}
		sb.WriteString("\r\n")
	}

	out.WriteString(sb.String())
	out.Flush()

	return true
}
