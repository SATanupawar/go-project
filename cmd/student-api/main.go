package main

import (
	"fmt"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"log/slog"
	"context"
	"time"
	"github.com/satyampawar/go-project/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	// server route
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to student api"))
	})

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
