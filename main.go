package main

import (
	"encurtador-go/api"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	if err := run(); err != nil {
		slog.Error("something failed to execute program", "error", err)
		return
	}
	slog.Info("All system offline")
}

func run() error {
	// In-memory database
	db := make(map[string]string)
	handler := api.NewHandler(db)	

	serve := http.Server{
		ReadTimeout: 10 * time.Second,
		IdleTimeout: time.Minute,
		WriteTimeout: 10 * time.Second,
		Addr: ":8080",
		Handler: handler,
	}

	if err := serve.ListenAndServe(); err != nil {
		return err
	}
	return nil
}