package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aidenfine/recruiter-rater/internal/api"
	"github.com/aidenfine/recruiter-rater/internal/config"
	"github.com/aidenfine/recruiter-rater/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func main() {
	config.LoadConfig()
	cfg := config.Get()

	db, err := db.Connect(cfg)
	if err != nil {
		fmt.Println("Failed to connect to db")
		panic(err)
	}

	defer db.Close()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1Router := api.NewV1Router(db)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from api"))
	})
	r.Route("/api", func(r chi.Router) {
		r.Mount("/v1", v1Router.V1Router())
	})
	server := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	go func() {
		slog.Info("Server is starting on " + cfg.Port)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()
	// handle shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
	slog.Info("Server exited")

}
