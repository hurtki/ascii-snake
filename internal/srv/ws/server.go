package ws

import (
	"context"
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hurtki/ascii-snake/internal/srv/app"
)

type Server struct {
	game *app.Game

	sm *SessionManager

	mu     sync.Mutex
	logger *slog.Logger
}

func NewServer(game *app.Game, logger *slog.Logger, sm *SessionManager) *Server {
	return &Server{
		sm:     sm,
		game:   game,
		logger: logger.With("service", "ws-server"),
	}
}

func (s *Server) HandleWS(conn *websocket.Conn, token string) {
	if !s.sm.SessionExists(context.TODO(), token) {
		s.logger.Warn("no session found for a new established ws connection, closing conn", "given_tok", token, "addr", conn.RemoteAddr())
		conn.Close()
		return
	}

	if s.sm.SessionHasConn(context.TODO(), token) {
		s.logger.Warn("session already has an established ws connection, closing a new", "given_tok", token, "addr", conn.RemoteAddr())
		conn.Close()
		return
	}

	s.sm.LinkConnectionToSession(context.TODO(), token, conn)

	s.logger.Info("linked new connection to an existing session", "tok", token, "addr", conn.RemoteAddr())

	go s.readLoop(conn, token)
	go s.WriteLoop(conn, token)
}

func (s *Server) readLoop(conn *websocket.Conn, token string) {
	s.logger.Info("started read loop", "tok", token)
	for {
		_, buf, err := conn.ReadMessage()
		if err != nil {
			s.logger.Error("can't read message, closing session", "err", err, "tok", token)
			s.sm.CloseSession(context.TODO(), token)
			return
		}

		dir, err := app.NewDirection(uint8(buf[0]))
		if err != nil {
			s.sm.CloseSession(context.TODO(), token)
			return
		}

		snakeID := s.sm.GetSessionPlayerID(context.TODO(), token)

		s.logger.Debug("Move", "direction", dir, "tok", token, "player_id", snakeID, "addr", conn.RemoteAddr())

		s.game.AddMove(snakeID, app.Move{Direction: dir})
	}
}

func (s *Server) WriteLoop(conn *websocket.Conn, token string) {
	s.logger.Info("started write loop", "tok", token)
	playerID := s.sm.GetSessionPlayerID(context.TODO(), token)
	for {
		interestZone := s.game.GetInterestZoneForSnakeAfterTick(playerID)
		payload := serializeInterestZone(interestZone)

		if !s.sm.SessionExists(context.TODO(), token) {
			return
		}

		err := conn.WriteMessage(websocket.BinaryMessage, payload)
		if err != nil {
			s.logger.Info("connection closed, closing session", "reason", "error writing interest zone payload", "err", err)
			s.sm.CloseSession(context.TODO(), token)
			return
		}
	}
}
