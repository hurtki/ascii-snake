package app

import (
	// "sync"
	"time"

	"sync"

	"golang.org/x/sync/singleflight"
)

type GameConfig struct {
	BaseSnakeLength int
	XSize, YSize    int
	TickTime        time.Duration
	// (2InterestSize+1)*(1InterestSize+1) is a square player is supposed
	// to see on client
	InterestSize int

	// Measured in grid cells ( chunks )
	PlayerSpawnPaddingFromBorder int
}

type Game struct {
	snakes map[int]Snake
	apples map[Cord]struct{}

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
	moves   map[int]Move
	movesMu sync.Mutex

	// used on tick "stop the world"
	mu sync.RWMutex

	// After every tick one struct is sent to notify readers to check plot
	AfterTickCh chan struct{}

	// singleflight is used for distribution of a new map after tick
	// used in GetMapCopyAfterTick()
	sf singleflight.Group

	// interest grid map
	im                                   map[interestGridCord]interestGridCell
	interestGridSizeX, interestGridSizeY int
}

func InitGame(cfg GameConfig) *Game {
	return &Game{
		cfg:               cfg,
		snakes:            make(map[int]Snake),
		apples:            make(map[Cord]struct{}),
		moves:             make(map[int]Move),
		cntr:              1,
		AfterTickCh:       make(chan struct{}),
		im:                make(map[interestGridCord]interestGridCell),
		interestGridSizeX: (cfg.XSize + cfg.InterestSize - 1) / cfg.InterestSize,
		interestGridSizeY: (cfg.YSize + cfg.InterestSize - 1) / cfg.InterestSize,
	}
}

func (g *Game) Start() {
	t := time.NewTicker(g.cfg.TickTime)

	for {
		<-t.C

		// TICK TIME LOCK
		// STOP THE WORLD
		g.mu.Lock()

		g.updateInterestGridMap()

		g.completeMoves()

		g.applyMoves()

		clear(g.moves)

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

		g.SpawnApples()

		g.mu.Unlock()

		select {
		case g.AfterTickCh <- struct{}{}:
		default:
		}
	}
}
