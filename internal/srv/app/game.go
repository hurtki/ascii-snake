package app

import (
	// "sync"
	"time"

	"sync"

	"golang.org/x/sync/singleflight"
)

type GameConfig struct {
	BaseSnakeLength int
	xSize, ySize    int
	TickTime        time.Duration
}

type Game struct {
	snakes map[int]Snake

	cfg GameConfig

	// ID for new player
	// increment after adding a new one
	// starts with 1, cause 0 is zero value
	cntr int

	// add a new player requests
	// slice of callback functions
	addQueue   []func(int, error)
	addQueueMu sync.Mutex

	// apply move queue
	// slice of moves that came to server
	movesQueue   []Move
	movesQueueMu sync.Mutex

	// used on tick "stop the world"
	mu sync.RWMutex

	// After every tick one struct is sent to notify readers to check plot
	AfterTickCh chan struct{}

	// singleflight is used for distribution of a new map after tick
	// used in GetMapCopyAfterTick()
	sf singleflight.Group
}

func InitGame(xSize, ySize int) *Game {
	return &Game{
		cntr:        1,
		AfterTickCh: make(chan struct{}),
	}
}

func (g *Game) Start() {
	t := time.NewTicker(g.cfg.TickTime)

	for {
		<-t.C

		// TICK TIME LOCK
		// STOP THE WORLD
		g.mu.Lock()

		moves := DeduplicateMoves(g.movesQueue)
		g.applyMoves(moves)

		// clear moves queue
		g.movesQueue = g.movesQueue[:0]

		for i, callback := range g.addQueue {
			id, ok := g.createPlayer()
			if !ok {
				for j := i; i < len(g.addQueue); i++ {
					g.addQueue[j](0, ErrNoPlaceOnPlot)
				}
				break
			}
			callback(id, nil)
		}
		g.addQueue = g.addQueue[:0]

		g.mu.Unlock()

		select {
		case g.AfterTickCh <- struct{}{}:
		default:
		}
	}
}

func (g *Game) GetMapCopyAfterTick() [][]Cell {
	res, _, _ := g.sf.Do("", func() (any, error) {
		<-g.AfterTickCh
		g.mu.RLock()

		res := make([][]Cell, len(g.plot))
		for i := range res {
			res[i] = make([]Cell, len(g.plot[i]))
			copy(res[i], g.plot[i])
		}

		g.mu.RUnlock()
		return res, nil
	})
	return res.([][]Cell)
}
