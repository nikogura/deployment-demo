package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nikogura/deployment-demo/pkg/demo"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := demo.LoadConfig()

	err := demo.Run(ctx, cfg, logger)
	if err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
