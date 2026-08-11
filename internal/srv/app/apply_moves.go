package app

func (g *Game) applyMoves(moves []Move) {
	moveMap := make(map[int]Move, len(g.snakes))

	for _, m := range moves {
		if _, ok := moveMap[m.SnakeID]; !ok {
			moveMap[m.SnakeID] = m
		}
	}

	// delete all snakes that went to the border
	for snakeID, snake := range g.snakes {
		// if there was no move registered for a snake
		// then add it a move that goes "forward"
		if _, ok := moveMap[snakeID]; !ok {
			moveMap[snakeID] = Move{SnakeID: snakeID, Direction: snake.Body[0].Opposite()}
		}

		inbound := snake.Cord.
			// apply a move
			Apply(moveMap[snakeID].Direction).
			// check if it's inbouds
			InBound(g.cfg.xSize, g.cfg.ySize)

		if !inbound {
			// if snake went to border, than delete it from map
			delete(g.snakes, snakeID)
		}
	}

	xLen := (g.cfg.xSize + g.cfg.InterestSize - 1) / g.cfg.InterestSize
	yLen := (g.cfg.ySize + g.cfg.InterestSize - 1) / g.cfg.InterestSize

	interestGrid := make([][]interestGridCell, xLen)

	for i := range interestGrid {
		interestGrid[i] = make([]interestGridCell, yLen)
	}

	for id, snake := range g.snakes {
		interestGrid[snake.Cord.X/g.cfg.InterestSize][snake.Cord.Y/g.cfg.InterestSize].EnsureSnake(id, snake)
	}

	for cord, _ := range g.apples {
		interestGrid[cord.X/g.cfg.InterestSize][cord.Y/g.cfg.InterestSize].EnsureApple(cord)
	}

	for _, move := range moveMap {
		resultCord := g.snakes[move.SnakeID].Cord.Apply(move.Direction)
		interestGridCell := interestGrid[resultCord.X/g.cfg.InterestSize][resultCord.Y/g.cfg.InterestSize]

		hitSnake := false
		for id, potentialSnake := range interestGridCell.Snakes {
			if potentialSnake.Cord == resultCord {
				// if we hit someones head, delete both snakes
				delete(g.snakes, id)
				delete(g.snakes, move.SnakeID)
				hitSnake = true
				break
			}
			if potentialSnake.Contains(resultCord) {
				delete(g.snakes, move.SnakeID)
				hitSnake = true
				break
			}
		}
		if hitSnake {
			continue
		}

		hitApple := false
		for cord, _ := range interestGridCell.Apples {
			if resultCord == cord {
				hitApple = true
				// hit the apple
				s := g.snakes[move.SnakeID]
				s.Cord = s.Cord.Apply(move.Direction)
				s.Body = append([]Direction{move.Direction.Opposite()}, s.Body...)
				g.snakes[move.SnakeID] = s
				delete(g.apples, cord)
				break
			}
		}

		if hitApple {
			continue
		}

		s := g.snakes[move.SnakeID]
		s.Cord = s.Cord.Apply(move.Direction)
		s.Body = append([]Direction{move.Direction.Opposite()}, s.Body...)[:len(s.Body)]
		g.snakes[move.SnakeID] = s
	}

}

type interestGridCell struct {
	Snakes map[int]Snake
	Apples map[Cord]struct{}
}

func (c *interestGridCell) EnsureSnake(id int, s Snake) {
	c.Snakes[id] = s
}

func (c *interestGridCell) EnsureApple(cord Cord) {
	c.Apples[cord] = struct{}{}
}
