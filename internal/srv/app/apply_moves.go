package app

func (g *Game) applyMoves() {
	for snakeID, move := range g.moves {
		moveSnake := g.snakes[snakeID]

		if moveSnake.Body[0] == move.Direction {
			// if move is "backwards", replace it with "forward"
			move.Direction = move.Direction.Opposite()
		}
		// for every move check grid cell, where player's head moved
		// and for optimisation we can check only snakes that are in that spicific grid cell
		resultCord := moveSnake.Cord.Apply(move.Direction)

		if !resultCord.InBound(g.cfg.XSize, g.cfg.YSize) {
			g.removeSnakeFromInterestGrid(snakeID, moveSnake)
			delete(g.snakes, snakeID)
			continue
		}

		gridCordX := resultCord.X / g.cfg.InterestSize
		gridCordY := resultCord.Y / g.cfg.InterestSize

		interestGridCell := g.im[Cord{X: gridCordX, Y: gridCordY}]

		hitSnake := false
		for id, potentialSnake := range interestGridCell.Snakes {
			if potentialSnake.Cord == resultCord {
				// if we hit someones head, delete both snakes
				delete(g.snakes, id)
				delete(g.snakes, snakeID)
				hitSnake = true

				g.removeSnakeFromInterestGrid(id, potentialSnake)

				break
			}
			if potentialSnake.Contains(resultCord) {
				delete(g.snakes, snakeID)
				hitSnake = true
				break
			}
		}
		if hitSnake {
			g.removeSnakeFromInterestGrid(snakeID, moveSnake)
			continue
		}

		hitApple := false
		for cord := range interestGridCell.Apples {
			if resultCord == cord {
				hitApple = true
				// hit the apple
				s := g.snakes[snakeID]
				s.Cord = s.Cord.Apply(move.Direction)
				s.Body = append([]Direction{move.Direction.Opposite()}, s.Body...)
				g.snakes[snakeID] = s
				delete(g.apples, cord)
				break
			}
		}

		if hitApple {
			g.ensureSnakeOnInterestGrid(snakeID, moveSnake)
			continue
		}

		s := g.snakes[snakeID]
		s.Cord = s.Cord.Apply(move.Direction)
		s.Body = append([]Direction{move.Direction.Opposite()}, s.Body...)[:len(s.Body)]
		g.snakes[snakeID] = s

		g.removeSnakeFromInterestGrid(snakeID, moveSnake)
		g.ensureSnakeOnInterestGrid(snakeID, s)
	}
}
