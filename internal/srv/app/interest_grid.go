package app

func (g *Game) updateInterestGridMap() {
	clear(g.im)

	for id, snake := range g.snakes {
		xGridCord := snake.Cord.X / g.cfg.InterestSize
		yGridCord := snake.Cord.Y / g.cfg.InterestSize
		gridCord := Cord{X: xGridCord, Y: yGridCord}
		cell, ok := g.im[gridCord]
		if !ok {
			g.im[gridCord] = newInterestGridCell()
			cell = g.im[gridCord]
		}
		cell.EnsureSnake(id, snake)
		g.im[gridCord] = cell
	}

	for cord, _ := range g.apples {
		xGridCord := cord.X / g.cfg.InterestSize
		yGridCord := cord.Y / g.cfg.InterestSize

		gridCord := Cord{X: xGridCord, Y: yGridCord}
		cell, ok := g.im[gridCord]
		if !ok {
			g.im[gridCord] = newInterestGridCell()
			cell = g.im[gridCord]
		}
		cell.EnsureApple(cord)
		g.im[gridCord] = cell
	}
}

type interestGridCell struct {
	Snakes map[int]Snake
	Apples map[Cord]struct{}
}

func newInterestGridCell() interestGridCell {
	return interestGridCell{
		Snakes: make(map[int]Snake),
		Apples: make(map[Cord]struct{}),
	}
}

func (c *interestGridCell) EnsureSnake(id int, s Snake) {
	c.Snakes[id] = s
}

func (c *interestGridCell) EnsureApple(cord Cord) {
	c.Apples[cord] = struct{}{}
}
