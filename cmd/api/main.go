package main

import (
	"log"
	"os"
	"time"

	"github.com/LincolnG4/Haku/internal/auth"
	"github.com/LincolnG4/Haku/internal/db"
	"github.com/LincolnG4/Haku/internal/store"
	"github.com/LincolnG4/Haku/internal/utils"
)

func main() {
	tokenSecret, ok := os.LookupEnv("AUTH_TOKEN_SECRET")
	if !ok {
		panic("missing AUTH_TOKEN_SECRET env ")
	}

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
		auth: authConfig{
			token: tokenConfig{
				secret:     tokenSecret,
				expiration: 24 * time.Hour,
				iss:        "Haku",
			},
		},
	}

	// Setup database connection
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	store := store.NewPostgresStorage(db)
	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	// Setup API server
	app := application{
		config:        cfg,
		store:         store,
		authenticator: jwtAuthenticator,
	}

	// Start server
	mux := app.mount()
	log.Fatal(app.run(mux))

}
