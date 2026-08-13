package app

import "maps"

func (g *Game) GetInterestCellForSnake(snakeID int) interestGridCell {
	g.mu.RLock()

	g.sf.Do("", func() (any, error) {
		<-g.AfterTickCh
		return nil, nil
	})

	s, ok := g.snakes[snakeID]
	if !ok {
		if len(g.snakes) > 0 {
			// user's snake not on the field
			// find other random and use it instead
			for id, snake := range g.snakes {
				snakeID = id
				s = snake
			}
		} else {
			// no one is on the map
			// just center of the map
			s = Snake{Cord: Cord{g.cfg.xSize / 2, g.cfg.ySize / 2}}
		}
	}

	gridCordX := s.Cord.X / g.cfg.InterestSize
	gridCordY := s.Cord.Y / g.cfg.InterestSize

	res := newInterestGridCell()

	for x := gridCordX - 1; x <= gridCordX+1; x++ {
		for y := gridCordY - 1; y <= gridCordY+1; y++ {
			if cell, ok := g.im[Cord{X: x, Y: y}]; ok {
				maps.Copy(res.Snakes, cell.Snakes)
				maps.Copy(res.Apples, cell.Apples)
			}
		}
	}

	// if in zone of interest there is only player's snake
	// find the nearest one
	if len(res.Snakes) < 2 && len(g.snakes) > 1 {
		cord := Cord{X: gridCordX, Y: gridCordY}

		stepLength := 1
		dir := Right

		found := false

		xGridSize := (g.cfg.xSize + g.cfg.InterestSize - 1) / g.cfg.InterestSize
		yGridSize := (g.cfg.ySize + g.cfg.InterestSize - 1) / g.cfg.InterestSize

		for !found {
			for range 2 {
				for s := 0; s < stepLength; s++ {
					cord = cord.Apply(dir)

					if !cord.InBound(xGridSize, yGridSize) {
						continue
					}

					cell, ok := g.im[cord]
					if !ok {
						continue
					}
					_, ok = cell.Snakes[snakeID]

					if (len(cell.Snakes) == 1 && !ok) || (len(cell.Snakes) > 1) {
						for id, outOfInterestSnake := range g.im[cord].Snakes {
							if id != snakeID {
								res.EnsureSnake(id, outOfInterestSnake)
								found = true
							}
						}
					}

				}
				dir = dir.NextClockwise()
			}
			stepLength++
		}
	}

	g.mu.RUnlock()

	return res
}
