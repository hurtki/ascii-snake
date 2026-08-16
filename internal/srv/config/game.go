package config

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/DeanPDX/dotconfig"
)

type GameConfig struct {
	BaseSnakeLength int `env:"BASE_SNAKE_LENGTH,required"`
	XSize           int `env:"MAP_HEIGHT,required"`
	YSize           int `env:"MAP_WIDTH,required"`
	TickTime        time.Duration
	TickTimeStr     string `env:"TICK_DURATION,required"`
	// (2InterestSize+1)*(1InterestSize+1) is a square player is supposed
	// to see on client
	InterestSize int `env:"INTEREST_RADIUS,required"`

	// Measured in grid cells ( chunks )
	PlayerSpawnPaddingFromBorder int `env:"SPAWN_CHUNK_PADDING_FROM_BORDER,required"`
}

func LoadGameConfigFromEnv() (GameConfig, error) {
	var buf bytes.Buffer

	for _, e := range os.Environ() {
		// e = "KEY=VALUE"
		buf.WriteString(e)
		buf.WriteByte('\n')
	}

	cfg, err := dotconfig.FromReader[GameConfig](&buf)

	duration, err := time.ParseDuration(cfg.TickTimeStr)
	if err != nil {
		return cfg, fmt.Errorf("invalid duration formatting for TICK_DURATION, %w", err)
	}
	cfg.TickTime = duration

	return cfg, err
}
