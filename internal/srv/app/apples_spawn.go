package app

import "math/rand"

// Tick time
func (g *Game) SpawnApples() {
	applesToSpawn := (len(g.snakes) - len(g.apples))

	// using range by map to take ~~~random snake and give it an apple into its
	// interest range

	for _, snake := range g.snakes {
		if applesToSpawn < 1 {
			return
		}
		gridCord := snake.Cord.ToInterestGrid(g.cfg.InterestSize)
		gridCell := g.im[gridCord]

		xRandChunkCord := rand.Int() % g.cfg.InterestSize
		yRandChunkCord := rand.Int() % g.cfg.InterestSize

		newAppleCord := Cord{
			X: gridCord.X*g.cfg.InterestSize + xRandChunkCord,
			Y: gridCord.Y*g.cfg.InterestSize + yRandChunkCord,
		}

		// check if a new cord is captured by some snake
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
		}
	}
}
