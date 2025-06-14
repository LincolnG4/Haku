package main

import (
	"log"

	"github.com/LincolnG4/Haku/internal/db"
	"github.com/LincolnG4/Haku/internal/store"
	"github.com/LincolnG4/Haku/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	addr := utils.GetEnvString("DB_CONNECTION_STRING", "")
	conn, err := db.New(addr, 30, 30, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewPostgresStorage(conn)
	db.Seed(store)
}
