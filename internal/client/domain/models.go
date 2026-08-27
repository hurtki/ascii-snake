package domain

type Snake struct {
	PlayerID     int
	CordX, CordY int
	Body         []Direction
}

type Apple struct {
	CordX, CordY int
}

type Direction uint8

const (
	Right Direction = iota
	Down
	Left
	Up
)
