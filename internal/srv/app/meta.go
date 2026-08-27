package app

import "time"

func (g *Game) GetXSize() int              { return g.cfg.XSize }
func (g *Game) GetYSize() int              { return g.cfg.YSize }
func (g *Game) GetTickTime() time.Duration { return g.cfg.TickTime }
func (g *Game) GetInterestSize() int       { return g.cfg.InterestSize }
func (g *Game) GetBaseSnakeLength() int    { return g.cfg.BaseSnakeLength }
