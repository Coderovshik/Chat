package main

import (
	"log"

	"github.com/Coderovshik/chat_server/internal/config"
	"github.com/Coderovshik/chat_server/internal/db"
)

func main() {
	cfg := config.New()
	m := db.NewMigrator(cfg)

	log.Println("migrating")
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
