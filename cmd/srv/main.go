package main

import (
	"log/slog"
	"os"

	"github.com/hurtki/ascii-snake/internal/srv/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logger.Info("starting")

	gameConfig, err := config.LoadGameConfigFromEnv()
	if err != nil {
		logger.Error("can't load game config", "err", err)
		return
	}

	logger.Info("time tick", "time_str", gameConfig.TickTime.String())

	// game := app.InitGame(gameConfig)
	//
	// go game.Start()
	//
	// sessionManager := ws.NewSessionManager()
	//
	// usecase := domain.NewGameUsecase(game, sessionManager)
	//
	// // wsHandler := ws.NewServer(game, logger, sessionManager)
	//
	// joinHandler := http_handlers.NewJoinHandler(usecase)
	//
	//	var upgrader = websocket.Upgrader{
	//		CheckOrigin: func(r *http.Request) bool {
	//			return true
	//		},
	//		ReadBufferSize:  1024,
	//		WriteBufferSize: 1024,
	//	}
	//
	//	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	//		token := r.URL.Query().Get("token")
	//
	//		_, err := upgrader.Upgrade(w, r, nil)
	//		if err != nil {
	//			fmt.Printf("Upgrade error: %v\n", err)
	//			return // Ответ со статусом ошибки отправится автоматически
	//		}
	//		fmt.Println("token got:", token)
	//
	//		// wsHandler.HandleWS(conn, token)
	//	})
	//
	// http.HandleFunc("GET /room", joinHandler.Join)
	//
	// // go wsHandler.WriteLoop()
	//
	// http.ListenAndServe(":3310", nil)
}
