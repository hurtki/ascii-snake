package ws

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
)

type sessionEntry struct {
	PlayerID int
	Conn     *websocket.Conn

	CloseOnce sync.Once
}

type SessionManager struct {
	// key: session token
	sessions map[string]*sessionEntry

	mu sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*sessionEntry),
	}
}

// Create adds a new session with given token and playerID
// New session doesn't have a WS connection
func (sm *SessionManager) Create(ctx context.Context, token string, playerID int) error {
	sm.mu.Lock()

	sm.sessions[token] = &sessionEntry{PlayerID: playerID}

	sm.mu.Unlock()

	return nil
}

// LinkConnectionToSession links given connection to already existing session
func (sm *SessionManager) LinkConnectionToSession(ctx context.Context, token string, conn *websocket.Conn) error {

	sm.mu.Lock()

	sess := sm.sessions[token]
	sess.Conn = conn
	sm.sessions[token] = sess

	sm.mu.Unlock()

	return nil
}

func (sm *SessionManager) SessionExists(ctx context.Context, token string) bool {
	sm.mu.RLock()
	_, ok := sm.sessions[token]
	sm.mu.RUnlock()
	return ok
}

func (sm *SessionManager) SessionHasConn(ctx context.Context, token string) bool {
	sm.mu.RLock()
	sess, _ := sm.sessions[token]
	sm.mu.RUnlock()
	return sess.Conn != nil
}

func (sm *SessionManager) CloseSession(ctx context.Context, token string) {
	sm.mu.Lock()

	sess, ok := sm.sessions[token]
	if !ok {
		return
	}
	sess.CloseOnce.Do(func() {
		sess.Conn.Close()
	})

	delete(sm.sessions, token)

	sm.mu.Unlock()
}

func (sm *SessionManager) GetSessionPlayerID(ctx context.Context, token string) int {
	sm.mu.RLock()

	sess := sm.sessions[token]

	sm.mu.RUnlock()

	return sess.PlayerID
}
func (sm *SessionManager) GetAllConns() []*websocket.Conn {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	conns := make([]*websocket.Conn, 0, len(sm.sessions))
	for _, entry := range sm.sessions {
		if entry.Conn != nil {
			conns = append(conns, entry.Conn)
		}
	}

	return conns
}
