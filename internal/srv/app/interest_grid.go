package app

func (g *Game) ensureSnakeOnInterestGrid(id int, snake Snake) {
	for _, gridCord := range snake.GetInterestGridCords(g.cfg.InterestSize) {
		cell, ok := g.im[gridCord]
		if !ok {
			g.im[gridCord] = newInterestGridCell()
			cell = g.im[gridCord]
		}
		cell.EnsureSnake(id, snake)
		g.im[gridCord] = cell
	}
}

func (g *Game) removeSnakeFromInterestGrid(id int, snake Snake) {
	for _, gridCord := range snake.GetInterestGridCords(g.cfg.InterestSize) {
		cell, ok := g.im[gridCord]
		if !ok {
			continue
		}
		cell.RemoveSnake(id, snake)
		g.im[gridCord] = cell
	}
}

func (g *Game) updateInterestGridMap() {
	clear(g.im)

	for id, snake := range g.snakes {
		g.ensureSnakeOnInterestGrid(id, snake)
	}

	for cord := range g.apples {
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

type interestGridCord = Cord

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

func (c *interestGridCell) RemoveSnake(id int, s Snake) {
	delete(c.Snakes, id)
}

func (c *interestGridCell) EnsureApple(cord Cord) {
	c.Apples[cord] = struct{}{}
}

type InterestZone = interestGridCell

func NewInterestZone() InterestZone {
	return InterestZone{
		Snakes: make(map[int]Snake),
		Apples: make(map[Cord]struct{}),
	}
}
