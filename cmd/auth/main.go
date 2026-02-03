package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/meis1kqt/sso/internal/app"
	"github.com/meis1kqt/sso/internal/config"
)

func main() {

	cfg := config.MustLoad()
	
	var logger *slog.Logger

	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	slog.SetDefault(logger)


	application := app.New(logger, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCServer.MustRun()


	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sing := <- stop

	logger.Info("sing", slog.String("signal", sing.String()))

	application.GRPCServer.Stop()

	logger.Info("application stopped")

}
