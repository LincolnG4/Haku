package main

import (
	"net/http"
	"time"

	"github.com/LincolnG4/Haku/internal/auth"
	"github.com/LincolnG4/Haku/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	authenticator auth.Authenticator
	logger        *zap.SugaredLogger
}

type config struct {
	addr string `validate:"required"`
	db   dbConfig
	env  string
	auth authConfig
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type basicConfig struct {
	user     string
	password string
}

type tokenConfig struct {
	secret     string
	expiration time.Duration
	iss        string
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
			router.Use(app.AuthTokenMiddleware)
			router.Post("/", app.createPipelineHandler)

			router.Route("/{pipelineID}", func(router chi.Router) {
				router.Use(app.pipelineContextMiddleware)
				// Routes
				router.Get("/", app.getPipelineHandler)
				router.Patch("/", app.updatePipelineHandler)
				router.Delete("/", app.deletePipelineHandler)
			})
		})

		// Task
		// router.Route("/Tasks", func(router chi.Router) {
		// 	router.Route("/{taskID}", func(router chi.Router) {
		// 		router.Post("/", app.createPipelineHandler)
		// 	})
		// })

		// User
		router.Route("/users", func(router chi.Router) {
			router.Route("/{userID}", func(router chi.Router) {
				router.Use(app.AuthTokenMiddleware)
				router.Get("/", app.getUserHandler)
			})

		})

		// Organizations
		router.Route("/organizations", func(router chi.Router) {
			router.Use(app.AuthTokenMiddleware)

			router.Post("/", app.createOrganizationHandler)
			router.Route("/{organizationID}", func(router chi.Router) {
				router.Use(app.organizationContextMiddleware)

				router.Get("/", app.getOrganizationHandler)
			})

		})

		router.Route("/auth", func(router chi.Router) {
			router.Post("/user", app.registerUserHandler)
			router.Post("/token", app.createTokenHandler)
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

	a.logger.Info("server has started at", "addr", a.config.addr)
	return srv.ListenAndServe()
}
