package main

import "log/slog"

func main() {
	if err := run(); err != nil {
		slog.Error("something failed to execute program", "error", err)
		return
	}
	slog.Info("All system offline")
}

func run() error {
	return nil
}