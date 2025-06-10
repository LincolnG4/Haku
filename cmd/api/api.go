package main

import (
	"log"
	"net/http"
	"time"

	"github.com/LincolnG4/Haku/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/v1", func(router chi.Router) {
		router.Get("/health", app.healthCheckHandler)

		// Pipeline
		router.Route("/pipelines", func(router chi.Router) {
			router.Post("/", app.createPipelineHandler)
		})

		// Task
		router.Route("/Tasks", func(router chi.Router) {
			router.Post("/", app.createPipelineHandler)
		})

	})

	return router
}

func (a *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         a.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", a.config.addr)
	return srv.ListenAndServe()
}
