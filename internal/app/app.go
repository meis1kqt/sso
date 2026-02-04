package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/meis1kqt/sso/internal/app/grpc"
	"github.com/meis1kqt/sso/internal/services/auth"
	"github.com/meis1kqt/sso/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App 
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		slog.Error("cannot connect", "error", err)
		panic(err)
	}

	authService := auth.New(log, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{GRPCServer: grpcApp}
}