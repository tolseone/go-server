package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"

	authgrpc "go-server/internal/clients/auth/grpc"
	"go-server/internal/config"
	"go-server/internal/router"
	"go-server/internal/schedule"
	"go-server/pkg/logging"

)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")

	cfg := config.GetConfig()
	router := router.GetRouter(cfg)

	loggerSlog := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	authClient, err := authgrpc.New(context.Background(), loggerSlog, cfg.Clients.Auth.Address, cfg.Clients.Auth.Timeout, cfg.Clients.Auth.RetriesCount)
	if err != nil {
		logger.Error("failed to init auth client", err)
		os.Exit(1)
	}

	authClient.IsAdmin(context.Background(), "f854a904-33ec-4c99-a893-e27ab8f9a83d")

	go schedule.ScheduleTask()

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	logger.Info("listen tcp")
	listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
