package app

// adds move to queue that will be applied in tick time
// can't return error, but there is not guarantee, that move will be applied
// *for example if there were two, only first one will be applied
// not for tick time!
func (g *Game) AddMove(snakeID int, move Move) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g.addQueueMu.Lock()

	if _, ok := g.snakes[snakeID]; !ok {
		return
	}

	g.moves[snakeID] = move
	g.addQueueMu.Unlock()
}
