package main

import (
	"github.com/Coderovshik/chat_server/internal/app"
	"github.com/Coderovshik/chat_server/internal/config"
)

func main() {
	cfg := config.New()
	app := app.New(cfg)

	app.Run()
}
