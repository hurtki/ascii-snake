package main

import (
	"context"
	"os"
	"time"

	"github.com/hurtki/ascii-snake/internal/client/domain"
)

type gameInput interface {
	SendMove(domain.Direction)
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
					gameInput.SendMove(domain.Up)
				case 's', 'S':
					gameInput.SendMove(domain.Down)
				case 'a', 'A':
					gameInput.SendMove(domain.Left)
				case 'd', 'D':
					gameInput.SendMove(domain.Right)
				case 'q', 'Q', 3:
					cancel()
					return
				}
			} else if n == 3 && buf[0] == 27 && buf[1] == 91 {
				switch buf[2] {
				case 'A':
					gameInput.SendMove(domain.Up)
				case 'B':
					gameInput.SendMove(domain.Down)
				case 'D':
					gameInput.SendMove(domain.Left)
				case 'C':
					gameInput.SendMove(domain.Right)
				}
			}
		}
	}
}
