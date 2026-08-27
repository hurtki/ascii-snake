package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hurtki/ascii-snake/internal/client/domain"
	"github.com/hurtki/ascii-snake/internal/client/game_ui"
)

type GameConnection struct {
	addr      string
	conn      *websocket.Conn
	gameUICfg game_ui.GameUICfg
}

func NewGameConnection(ctx context.Context, addr string) (*GameConnection, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s:3310/connect", addr), nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't run request: %w", err)
	}

	dto := ConnectResponseDTO{}
	body, _ := io.ReadAll(res.Body)

	// err = json.NewDecoder(res.Body).Decode(&dto)
	err = json.Unmarshal(body, &dto)
	if err != nil {
		return nil, fmt.Errorf("on http stage server responded with wrong json schema: %w, body: %s", err, string(body))
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, fmt.Sprintf("ws://%s:3310/ws?token=%s", addr, dto.Token), nil)
	if err != nil {
		return nil, fmt.Errorf("can't establish websocket connection")
	}

	return &GameConnection{
		addr: addr,
		conn: conn,
		gameUICfg: game_ui.GameUICfg{
			XSize:         dto.MapSizeX,
			YSize:         dto.MapSizeY,
			InterestSize:  dto.InterestSize,
			PlayerSnakeID: dto.PlayerID,
		},
	}, nil

}

type GameDrawer interface {
	DrawScreen(snakes []domain.Snake, apples []domain.Apple)
}

func (gc *GameConnection) StartDrawing(ctx context.Context, drawer GameDrawer) error {
	for ctx.Err() == nil {
		_, data, err := gc.conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("can't read message from conn: %w", err)
		}
		snakes, apples, err := deserializeBinaryMapResponse(data)
		if err != nil {
			return fmt.Errorf("can't deserialize message from server: %w", err)
		}
		drawer.DrawScreen(snakes, apples)
	}
	return nil
}

func (gc *GameConnection) GetGameUICfg() game_ui.GameUICfg {
	return gc.gameUICfg
}

type ConnectResponseDTO struct {
	Token        string `json:"token"`
	PlayerID     int    `json:"player_id"`
	MapSizeX     int    `json:"map_size_x"`
	MapSizeY     int    `json:"map_size_y"`
	InterestSize int    `json:"interest_size"`
}

func (gc *GameConnection) Close() error {
	return gc.conn.Close()
}

func (gc *GameConnection) SendMove(move domain.Direction) {
	gc.conn.WriteMessage(websocket.BinaryMessage, []byte{byte(move)})
}
