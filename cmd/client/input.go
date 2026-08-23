package main

import (
	"context"
	"os"
	"time"
)

type Direction uint8

const (
	Right Direction = iota
	Down
	Left
	Up
)

type gameInput interface {
	SendMove(Direction)
}

func handleInput(ctx context.Context, cancel context.CancelFunc, gameInput gameInput) {
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
					gameInput.SendMove(Up)
				case 's', 'S':
					gameInput.SendMove(Down)
				case 'a', 'A':
					gameInput.SendMove(Left)
				case 'd', 'D':
					gameInput.SendMove(Right)
				case 'q', 'Q', 3:
					cancel()
					return
				}
			} else if n == 3 && buf[0] == 27 && buf[1] == 91 {
				switch buf[2] {
				case 'A':
					gameInput.SendMove(Up)
				case 'B':
					gameInput.SendMove(Down)
				case 'D':
					gameInput.SendMove(Left)
				case 'C':
					gameInput.SendMove(Right)
				}
			}
		}
	}
}
