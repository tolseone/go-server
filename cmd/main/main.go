package main

import (
	"context"
	"fmt"
	"go-server/internal/config"
	"go-server/internal/item"
	"go-server/internal/user"
	us "go-server/internal/user/db"
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
	repository := us.NewRepository(PostgreSQLClient, logger)
	logger.Info("connected to user repository")

	// func Create()
	// newUser := user.User{
	// 	Username: "Dima",
	// 	Email:    "dima@mail.ru",
	// }
	// if err := repository.Create(context.TODO(), &newUser); err != nil {
	// 	logger.Fatalf("Failed to create user: %v", err)
	// }
	// logger.Infof("Created user: %v", newUser)

	// func FindOne()
	one, err := repository.FindOne(context.TODO(), "380ec643-c806-49ec-89bc-c0bf3c581e55")
	if err != nil {
		logger.Fatalf("Failed to find user: %v", err)
	}
	logger.Infof("Found user: %v", one)

	// func FindAll()
	all, err := repository.FindAll(context.TODO())
	if err != nil {
		logger.Fatalf("Failed to find all users: %v", err)
	}

	for _, usr := range all {
		logger.Infof("%v", usr)
	}

	logger.Info("register item handler")
	handler := item.NewHandler(logger)
	handler.Reqister(router)
	handler = user.NewHandler(logger)
	handler.Reqister(router)

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
