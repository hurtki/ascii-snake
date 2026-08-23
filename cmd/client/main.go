package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
