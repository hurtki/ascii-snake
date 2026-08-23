package app

import "maps"

func (g *Game) GetInterestZoneForSnakeAfterTick(snakeID int) InterestZone {

	g.sf.Do("", func() (any, error) {
		<-g.AfterTickCh
		return nil, nil
	})

	g.mu.RLock()

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
			s = Snake{Cord: Cord{g.cfg.XSize / 2, g.cfg.YSize / 2}}
		}
	}

	gridCord := s.Cord.ToInterestGrid(g.cfg.InterestSize)

	res := NewInterestZone()

	for x := gridCord.X - 1; x <= gridCord.X+1; x++ {
		for y := gridCord.Y - 1; y <= gridCord.Y+1; y++ {
			if cell, ok := g.im[Cord{X: x, Y: y}]; ok {
				maps.Copy(res.Snakes, cell.Snakes)
				maps.Copy(res.Apples, cell.Apples)
			}
		}
	}

	// if in zone of interest there is only player's snake
	// find the nearest one
	if len(res.Snakes) < 2 && len(g.snakes) > 1 {
		stepLength := 1
		dir := Right

		found := false

		for !found {
			for range 2 {
				for s := 0; s < stepLength; s++ {
					gridCord = gridCord.Apply(dir)

					if !gridCord.InBound(g.interestGridSizeX, g.interestGridSizeY) {
						continue
					}

					cell, ok := g.im[gridCord]
					if !ok {
						continue
					}
					_, ok = cell.Snakes[snakeID]

					if (len(cell.Snakes) == 1 && !ok) || (len(cell.Snakes) > 1) {
						for id, outOfInterestSnake := range g.im[gridCord].Snakes {
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
