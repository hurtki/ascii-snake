package app

import "math/rand"

// CreatePlayer finds place for new snake in horizontal or vertical position
// Initializes cells for snake
// Tick time
func (g *Game) createPlayer() (int, bool) {
	xInnerSize := g.interestGridSizeX - 2*g.cfg.PlayerSpawnPaddingFromBorder
	total :=
		xInnerSize *
			(g.interestGridSizeY - 2*g.cfg.PlayerSpawnPaddingFromBorder)

	indices := make([]int, total)

	for i := range total {
		indices[i] = i
	}

	rand.Shuffle(total, func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	// now let's find a blank chunk of the map
	// to create a new snake there
	for _, idx := range indices {
		gridCord := interestGridCord{
			X: idx % xInnerSize,
			Y: idx / xInnerSize,
		}
		gridCord.X += g.cfg.PlayerSpawnPaddingFromBorder
		gridCord.Y += g.cfg.PlayerSpawnPaddingFromBorder

		_, ok := g.im[gridCord]
		if ok {
			continue
		}

		newSnakeHeadCord := Cord{
			X: gridCord.X*g.cfg.InterestSize + int(g.cfg.InterestSize/2),
			Y: gridCord.Y*g.cfg.InterestSize + (g.cfg.BaseSnakeLength - 1),
		}
		newBody := make([]Direction, g.cfg.BaseSnakeLength-1)
		for i := range g.cfg.BaseSnakeLength - 1 {
			newBody[i] = Left
		}
		newSnake := Snake{Cord: newSnakeHeadCord, Body: newBody}
		g.snakes[g.cntr] = newSnake

		// updaing interest map after creating player
		cell := newInterestGridCell()
		cell.EnsureSnake(g.cntr, newSnake)
		g.im[gridCord] = cell

		g.cntr++

		return g.cntr - 1, true
	}

	return 0, false
}
