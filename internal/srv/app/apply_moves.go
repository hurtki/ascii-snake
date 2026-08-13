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

	g.updateInterestGridMap()

	for _, move := range moveMap {
		// for every move check grid cell, where player's head moved
		// and for optimisation we can check only snakes that are in that spicific grid cell
		resultCord := g.snakes[move.SnakeID].Cord.Apply(move.Direction)
		gridCordX := resultCord.X / g.cfg.InterestSize
		gridCordY := resultCord.Y / g.cfg.InterestSize

		interestGridCell := g.im[Cord{X: gridCordX, Y: gridCordY}]

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
		for cord := range interestGridCell.Apples {
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
