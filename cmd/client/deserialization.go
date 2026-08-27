package main

import (
	"encoding/binary"
	"fmt"

	"github.com/hurtki/ascii-snake/internal/client/domain"
)

const (
	SnakeTypeOfObject = byte(1)
	AppleTypeOfObject = byte(2)
)

func deserializeBinaryMapResponse(data []byte) ([]domain.Snake, []domain.Apple, error) {
	if len(data) < 1 {
		return nil, nil, fmt.Errorf("blank slice of data")
	}
	curr := 0
	snakes := make([]domain.Snake, 0)
	apples := make([]domain.Apple, 0)

	for curr < len(data) {
		objType := data[curr]
		curr++
		bytesRemains := len(data) - curr
		switch objType {
		case AppleTypeOfObject:
			if bytesRemains < 4 {
				return nil, nil, fmt.Errorf("not enough space for apple object")
			}
			xCord := binary.LittleEndian.Uint16(data[curr:])
			curr += 2
			yCord := binary.LittleEndian.Uint16(data[curr:])
			curr += 2
			apples = append(apples, domain.Apple{CordX: int(xCord), CordY: int(yCord)})
		case SnakeTypeOfObject:
			if bytesRemains < 8 {
				return nil, nil, fmt.Errorf("not enough space for snake header")
			}
			snakeID := int(binary.LittleEndian.Uint16(data[curr:]))
			curr += 2
			xCord := int(binary.LittleEndian.Uint16(data[curr:]))
			curr += 2
			yCord := int(binary.LittleEndian.Uint16(data[curr:]))
			curr += 2
			snakeBodyLength := int(binary.LittleEndian.Uint16(data[curr:]))
			curr += 2

			var bodyBytesLength int = (snakeBodyLength + 3) / 4

			bytesRemains = len(data) - curr

			if bytesRemains < bodyBytesLength {
				return nil, nil, fmt.Errorf("not enough space for snake body, remains: %d, needed: %d", bytesRemains, bodyBytesLength)
			}
			snakesBody := make([]domain.Direction, 0, snakeBodyLength)
			for i := range snakeBodyLength {
				var b int = i / 4
				shift := 6 - 2*(i%4)
				dir := (data[curr+b] >> shift) & 0x03
				switch dir {
				case byte(0):
					snakesBody = append(snakesBody, domain.Right)
				case byte(1):
					snakesBody = append(snakesBody, domain.Down)
				case byte(2):
					snakesBody = append(snakesBody, domain.Left)
				case byte(3):
					snakesBody = append(snakesBody, domain.Up)
				}
			}

			curr += bodyBytesLength

			snakes = append(snakes, domain.Snake{PlayerID: snakeID, CordX: xCord, CordY: yCord, Body: snakesBody})
		default:
			return nil, nil, fmt.Errorf("wrong type of object")
		}
	}
	return snakes, apples, nil
}
