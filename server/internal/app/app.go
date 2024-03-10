package app

import (
	"log"
	"net/http"

	"github.com/Coderovshik/chat_server/internal/config"
	"github.com/Coderovshik/chat_server/internal/db"
	"github.com/Coderovshik/chat_server/internal/room"
	"github.com/Coderovshik/chat_server/internal/router"
	"github.com/Coderovshik/chat_server/internal/user"
	"github.com/gorilla/websocket"
)

type App struct {
	cfg    *config.Config
	router *router.Router
}

func New(cfg *config.Config) *App {
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("FATAL: failed to create database connection: %s", err.Error())
	}

	userRepo := user.NewRepository(database.GetDB())
	userService := user.NewService(userRepo, cfg)
	userHandler := user.NewHandler(userService, cfg)

	roomService := room.NewHub()
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	roomHandler := room.NewHandler(roomService, upgrader, cfg)

	r := router.NewRouter(userHandler, roomHandler)

	return &App{
		cfg:    cfg,
		router: r,
	}
}

func (a *App) Run() {
	log.Printf("server running %s", a.cfg.Addr())
	a.router.Run(a.cfg.Addr())
}
