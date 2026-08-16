package ws

import (
	"encoding/binary"

	"github.com/hurtki/ascii-snake/internal/srv/app"
)

const (
	SnakeTypeOfObject = byte(1)
	AppleTypeOfObject = byte(2)
)

func serializeInterestZone(interestZone app.InterestZone) []byte {
	approxLen := len(interestZone.Apples)*10 + len(interestZone.Snakes)*5
	res := make([]byte, 0, approxLen)
	curr := 0

	for snakeID, snake := range interestZone.Snakes {

		neededForSnakeHeader := 9
		res = append(res, make([]byte, neededForSnakeHeader)...)

		res[curr] = SnakeTypeOfObject
		curr++
		binary.LittleEndian.PutUint16(res[curr:curr+2], uint16(snakeID))
		curr += 2
		binary.LittleEndian.PutUint16(res[curr:curr+2], uint16(snake.Cord.X))
		curr += 2
		binary.LittleEndian.PutUint16(res[curr:curr+2], uint16(snake.Cord.Y))
		curr += 2
		binary.LittleEndian.PutUint16(res[curr:curr+2], uint16(len(snake.Body)))
		curr += 2

		neededForBody := (len(snake.Body) + 3) / 4
		res = append(res, make([]byte, neededForBody)...)

		for i, dir := range snake.Body {
			var b int = i / 4
			shift := 6 - 2*(i%4)
			res[curr+b] |= byte(dir) << shift
		}

		curr += neededForBody
	}
	for appleCord := range interestZone.Apples {
		neededForApple := 5
		res = append(res, make([]byte, neededForApple)...)

		res[curr] = AppleTypeOfObject
		curr++
		binary.LittleEndian.PutUint16(res[curr:curr+2], uint16(appleCord.X))
		curr += 2
		binary.LittleEndian.PutUint16(res[curr:curr+2], uint16(appleCord.Y))
		curr += 2
	}

	return res
}
