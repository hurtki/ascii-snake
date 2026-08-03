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
		logger: logger,
	}
}

func (s *Server) HandleWS(conn *websocket.Conn, token string) {
	if !s.sm.SessionExists(context.TODO(), token) {
		s.logger.Error("no session found, closing conn", "tok", token, "addr", conn.RemoteAddr())
		conn.Close()
		return
	}

	if s.sm.SessionHasConn(context.TODO(), token) {
		s.logger.Error("conn already exists, closing new one", "tok", token, "addr", conn.RemoteAddr())
		conn.Close()
		return
	}

	s.sm.LinkConnectionToSession(context.TODO(), token, conn)

	go s.readLoop(conn, token)
}

func (s *Server) readLoop(conn *websocket.Conn, token string) {
	for {
		_, buf, err := conn.ReadMessage()
		if err != nil {
			s.logger.Error("can't read message, closing session", "err", err, "tok", token)
			s.sm.CloseSession(context.TODO(), token)
			return
		}

		motion, err := app.NewDirection(uint8(buf[0]))
		if err != nil {
			s.sm.CloseSession(context.TODO(), token)
			return
		}
		s.logger.Debug("Move", "direction", motion, "tok", token, "addr", conn.RemoteAddr())

		playerID := s.sm.GetSessionPlayerID(context.TODO(), token)

		s.game.AddMove(app.Move{PlayerID: playerID, Direction: motion})
	}
}

func (s *Server) WriteLoop() {
	for {
		plot := s.game.GetMapCopyAfterTick()
		serializedPlot := SerializePlot(plot)

		for _, c := range s.sm.GetAllConns() {
			c.WriteMessage(websocket.BinaryMessage, serializedPlot)
		}
	}
}
