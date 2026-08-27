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

	"github.com/charmbracelet/x/term"
	"github.com/hurtki/ascii-snake/internal/client/game_ui"
)

const v string = "1.0.3"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	gameConn := &GameConnection{}
	var err error

	for {
		addr := Input(
			ctx,
			"ASCII Snake | OS Project | github:hurtki/ascii-snake | %s\nEnter address:",
			v,
		)
		if addr == "" {
			return
		}

		ctx, cancel := context.WithTimeout(ctx, time.Second*3)
		defer cancel()
		gameConn, err = NewGameConnection(ctx, addr)
		if err != nil {
			fmt.Printf("can't initialize connection with remote game: %s\n", err.Error())
			continue
		}
		fmt.Println("established connection with game", gameConn.GetGameUICfg())
		break
	}

	oldState, err := term.MakeRaw(os.Stdin.Fd())
	if err == nil {
		defer func() {
			_ = term.Restore(os.Stdin.Fd(), oldState)
		}()
	}

	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h\033[2J\033[H")

	gameUI := game_ui.NewGameUI(bufio.NewWriter(os.Stdout), gameConn.GetGameUICfg())

	go handleInput(ctx, cancel, gameConn)

	fmt.Print("\033[2J")

	gameConn.StartDrawing(ctx, gameUI)
	<-ctx.Done()

	gameConn.Close()
}

func Input(ctx context.Context, prompt string, args ...any) string {
	if prompt != "" {
		fmt.Printf(prompt, args...)
	}

	res := make(chan string, 1)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			res <- ""
		}
		res <- strings.TrimRight(text, "\r\n")
	}()

	select {
	case <-ctx.Done():
		return ""
	case text := <-res:
		return text
	}
}
