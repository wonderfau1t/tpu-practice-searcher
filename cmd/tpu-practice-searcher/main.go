package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"net/http"
	"os"
	"tpu-practice-searcher/internal/config"
	"tpu-practice-searcher/internal/http-server/handlers"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/logger"
	"tpu-practice-searcher/internal/storage/postgresql"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Info("Config successfully loaded")

	db, err := postgresql.SetupStorage(cfg.Storage)
	if err != nil {
		log.Error("Failed to set storage")
		os.Exit(1)
	}
	log.Info("Storage successfully set")

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middlewares.AuthMiddleware)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"message": "Hello world!"})
	})
	router.Get("/auth", handlers.Auth(log, db))
	router.Get("/register", handlers.RegisterStudent(log, db))

	http.ListenAndServe("0.0.0.0:8000", router)
}
