package app

// Tick time
// completeMoves adds moves to snakes that didn't register
// any move before tick
func (g *Game) completeMoves() {
	for snakeID, s := range g.snakes {
		if _, ok := g.moves[snakeID]; !ok {
			g.moves[snakeID] = Move{Direction: s.Body[0].Opposite()}
		}
	}
}
