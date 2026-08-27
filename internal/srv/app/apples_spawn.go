package app

import "math/rand"

// keeps len(g.snakes) on the field and if needed spawn new ones near
// the snake ( in the same chunk )
// Tick time
func (g *Game) SpawnApples1() {
	applesToSpawn := (len(g.snakes) - len(g.apples))

	// using range by map to take ~~~random snake and give it an apple into its
	// interest range

	for _, snake := range g.snakes {
		if applesToSpawn < 1 {
			return
		}
		gridCord := snake.Cord.ToInterestGrid(g.cfg.InterestSize)
		g.trySpawnAppleInGridCell(gridCord)
	}
}

// Keeps MAX ~(sizeX*sizeY)/(interestSize*interestSize) apples on the field
// one apple on every "loaded" chunk
// Tick time
func (g *Game) SpawnApples2() {
	for gridCord, cell := range g.im {
		if len(cell.Apples) > 0 {
			continue
		}
		g.trySpawnAppleInGridCell(gridCord)
	}
}

func (g *Game) trySpawnAppleInGridCell(gridCord interestGridCord) bool {
	gridCell, ok := g.im[gridCord]
	if !ok {
		gridCell = newInterestGridCell()
	}

	xRandChunkCord := rand.Int() % g.cfg.InterestSize
	yRandChunkCord := rand.Int() % g.cfg.InterestSize

	newAppleCord := Cord{
		X: gridCord.X*g.cfg.InterestSize + xRandChunkCord,
		Y: gridCord.Y*g.cfg.InterestSize + yRandChunkCord,
	}

	if !newAppleCord.InBound(g.cfg.XSize, g.cfg.YSize) {
		return false
	}

	captured := false
	for _, s := range gridCell.Snakes {
		if s.Contains(newAppleCord) {
			captured = true
			break
		}
	}
	if !captured {
		gridCell.EnsureApple(newAppleCord)
		g.im[gridCord] = gridCell
		g.apples[newAppleCord] = struct{}{}
		return true
	}
	return false
}
