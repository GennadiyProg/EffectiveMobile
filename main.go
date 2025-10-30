package main

import (
	"effective_mobile/internal/config"
	"effective_mobile/internal/handler"
	"effective_mobile/internal/repository/postgres"
	"effective_mobile/internal/service"
	"effective_mobile/pkg/database"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// @title Subscription API
// @version 1.0
// @description This is a sample subscription management API
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg := config.Load()
	slog.Debug("Debug log enabled")

	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Could not connect to database: " + err.Error())
		panic(err)
	}

	subRepository := postgres.NewSubscriptionRepository(db)
	subService := service.NewSubscriptionService(subRepository)

	subHandler := handler.NewSubscriptionHandler(subService)

	// Router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/subscription", func(r chi.Router) {
			r.Get("/", subHandler.ListSubscription)
			r.Get("/sum", subHandler.GetSubscriptionsSum)
			r.Get("/{id}", subHandler.GetSubscription)

			r.Post("/", subHandler.CreateSubscription)
			r.Put("/{id}", subHandler.UpdateSubscription)
			r.Delete("/{id}", subHandler.DeleteSubscription)
		})
	})

	slog.Info("Server starting on port " + cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		slog.Error("Could not star server on port " + cfg.Port + ": " + err.Error())
		panic(err)
	}
}
