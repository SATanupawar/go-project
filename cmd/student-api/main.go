package main

import (
	"fmt"

	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/satyampawar/go-project/internal/config"
	"github.com/satyampawar/go-project/internal/http/handlers/student"
	"github.com/satyampawar/go-project/internal/storage/sqlite"

	// Import SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	slog.Info("database connected")


	// server route
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetByID(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: router,
	}

	fmt.Printf("server started %s", cfg.HttpServer.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %s", err)
	}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to stop server: %s", slog.String("error", err.Error()))
	}

	slog.Info("server stopped")
	
}
