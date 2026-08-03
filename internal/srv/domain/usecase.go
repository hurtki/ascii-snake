package domain

import (
	"context"
	"math/rand"

	"github.com/hurtki/ascii-snake/internal/srv/app"

	"golang.org/x/sync/singleflight"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const baseTokenSymbolsCount = 50

func genRandomString(length int) string {
	str := make([]rune, length)
	for i := range length {
		str[i] = rune(alphabet[int(rand.Int31())%len(alphabet)])
	}
	return string(str)
}

type GameUsecase struct {
	game *app.Game
	sm   SessionManager

	sg singleflight.Group
}

type SessionManager interface {
	Create(ctx context.Context, token string, playerID int) error
}

func NewGameUsecase(game *app.Game, sm SessionManager) *GameUsecase {
	return &GameUsecase{
		game: game,
		sm:   sm,
	}
}

func (u *GameUsecase) JoinRoom(ctx context.Context) (JoinOut, error) {
	playerID, err := u.game.AddPlayer()
	if err != nil {
		return JoinOut{}, err
	}
	token := genRandomString(baseTokenSymbolsCount)
	size := u.game.GetMapSize()

	u.sm.Create(ctx, token, playerID)

	return JoinOut{
		Token:    token,
		MapSize:  size,
		PlayerID: playerID,
	}, nil
}
