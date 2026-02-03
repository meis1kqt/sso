package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authrpc "github.com/meis1kqt/sso/internal/grpc/auth"
	"google.golang.org/grpc"
)


type App struct {
	log *slog.Logger
	gRPC *grpc.Server
	port int
}

func New(log *slog.Logger,  port int ) *App {
	gRPCServer := grpc.NewServer()

	authrpc.Register(gRPCServer)

	return &App{
		log: log,
		gRPC: gRPCServer,
		port: port,
	} 
}

func (a *App) MustRun(){
	if err := a.Run(); err != nil {
		slog.Error("we have error run", "error", err)
		panic(err)
	}
}

func (a *App) Run() error {
	
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))

	if err != nil {
		slog.Error("run error", "message", err)
		return err
	}

	slog.Info("gRPC server started:","port", a.port)

	if err := a.gRPC.Serve(l); err != nil {
		slog.Error("we have error", "error", err)
		return err
	}
	return nil
}


func (a *App) Stop(){
	slog.Info("gRPC will stop")
	a.gRPC.GracefulStop()
}