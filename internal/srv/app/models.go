package app

import (
	"errors"
)

type Snake struct {
	// Cord is snake's head coordinates
	Cord Cord
	// Body contains slice of directions that represents snake body
	// for example if head is on (0, 0) and snake is like <--
	// then directions are going to be [right, right]
	// every directions says where is the next snake's cell
	Body []Direction
}

func (s Snake) Contains(cord Cord) bool {
	c := s.Cord
	for _, d := range s.Body {
		if c == cord {
			return true
		}
		c = c.Apply(d)
	}

	return c == cord
}

type Cord struct {
	X int
	Y int
}

func (c Cord) Apply(d Direction) Cord {
	switch d {
	case Up:
		c.X--
	case Down:
		c.X++
	case Left:
		c.Y--
	case Right:
		c.Y++
	}
	return c
}

func (c Cord) InBound(xSize, ySize int) bool {
	if c.X < 0 || c.Y < 0 || c.X >= xSize || c.Y >= ySize {
		return false
	}
	return true
}

// Direction is used to determine a move on the Plot
type Direction uint8

const (
	Up Direction = iota
	Down
	Left
	Right
)

var opposites = [...]Direction{
	Up:    Down,
	Down:  Up,
	Left:  Right,
	Right: Left,
}

func (d Direction) Opposite() Direction {
	return opposites[d]
}

func NewDirection(d uint8) (Direction, error) {
	if d > 3 {
		return Direction(0), errors.New("not existing direction")
	}
	return Direction(d), nil
}

type Move struct {
	SnakeID   int
	Direction Direction
}
