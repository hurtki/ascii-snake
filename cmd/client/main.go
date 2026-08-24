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

	for {
		addr := Input(ctx, "Enter address:")
		if addr == "" {
			break
		}
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		gameConn, err := NewGameConnection(ctx, addr)
		if err != nil {
			fmt.Printf("can't initialize connection with remote game: %s\n", err.Error())
			continue
		}
		fmt.Println("established connection with game", gameConn.GetGameUICfg())
		break
	}

	return

	oldState, err := term.MakeRaw(os.Stdin.Fd())
	if err == nil {
		defer func() {
			_ = term.Restore(os.Stdin.Fd(), oldState)
		}()
	}

	// get server address from user

	// initialize connection with server ( net.go )

	// initialize UI instance

	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h\033[2J\033[H")

	// wire connection to UI instance in order to start updating the TUI

	// wire input to network instance in order to send user's moves

	<-ctx.Done()
}

func Input(ctx context.Context, prompt string) string {
	if prompt != "" {
		fmt.Print(prompt)
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
