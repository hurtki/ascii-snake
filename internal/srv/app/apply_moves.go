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
}
