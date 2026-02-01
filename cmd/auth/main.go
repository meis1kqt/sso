package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/meis1kqt/sso/internal/config"
)

func main() {

	cfg := config.MustLoad()
	
	var logger *slog.Logger

	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	slog.SetDefault(logger)

}
