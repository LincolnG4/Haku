package main

import (
	"log"

	"github.com/LincolnG4/Haku/internal/db"
	"github.com/LincolnG4/Haku/internal/store"
	"github.com/LincolnG4/Haku/internal/utils"
)

func main() {
	// Setup Config
	cfg := config{
		addr: utils.GetEnvString("ADDR", ":8080"),
		db: dbConfig{
			addr:         utils.GetEnvString("DB_CONNECTION_STRING", ""),
			maxOpenConns: utils.GetEnvInt("DB_MAX_OPEN_CONNECTIONS", 30),
			maxIdleConns: utils.GetEnvInt("DB_MAX_IDLE_CONNECTIONS", 30),
			maxIdleTime:  utils.GetEnvString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: utils.GetEnvString("ENV", "development"),
	}

	// Setup database connection
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	store := store.NewPostgresStorage(db)

	// Setup API server
	app := application{
		config: cfg,
		store:  store,
	}

	// Start server
	mux := app.mount()
	log.Fatal(app.run(mux))

}
