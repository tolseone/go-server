package main

import (
	"context"
	"fmt"
	"go-server/internal/config"
	user "go-server/internal/controllers/handlers/handlerUser"
	item "go-server/internal/controllers/handlers/handleritem"
	"go-server/internal/models"
	it "go-server/internal/repositories/db/postgresItem"
	us "go-server/internal/repositories/db/postgresUser"
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

	PostgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	logger.Info("connected to PostgreSQL")

	// repositoryUser := us.NewRepository(PostgreSQLClient, logger)
	// logger.Info("connected to user repository")

	repositoryItem := it.NewRepository(PostgreSQLClient, logger)
	logger.Info("connected to item repository")

	modelItem := model.NewModelItem(repositoryItem)

	handler := item.NewHandler(logger, modelItem)
	handler.Reqister(router)
	logger.Info("registred item handler")

	// handler = user.NewHandler(logger, modelUser)
	// handler.Reqister(router)
	// logger.Info("registred user handler")

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
