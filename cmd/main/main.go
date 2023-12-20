package main

import (
	"context"
	"fmt"
	"go-server/internal/config"
	"go-server/internal/item"
	"go-server/internal/user"
	"go-server/pkg/client/postgresql"
	"go-server/pkg/logging"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	logger.Info("register item handler")
	handler := item.NewHandler(logger)
	handler.Reqister(router)
	handler = user.NewHandler(logger)
	handler.Reqister(router)

	client, err := postgresql.NewClient(context.Background(), 3, cfg.Storage)
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer client.Close()

	logger.Info("Connected to PostgreSQL!")

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	logger.Info("listen tcp")
	listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port)) // lookback interface
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
