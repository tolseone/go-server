package router

import (
	"github.com/julienschmidt/httprouter"

	"go-server/internal/config"
	"go-server/internal/controllers/handlers"
	"go-server/internal/controllers/handlers/api-service"
	"go-server/internal/controllers/handlers/auth-service"
	"go-server/pkg/logging"

)

const (
	tradesURL     = "/api/trades"
	tradeURL      = "/api/trades/:uuid"
	usertradesURL = "/api/users/:uuid/trades"
	itemtradesURL = "/api/items/:uuid/trades"

	itemsURL = "/api/items"
	itemURL  = "/api/items/:uuid"

	usersURL = "/api/users"
	userURL  = "/api/users/:uuid"

	registerURL = "/api/register"
	loginURL    = "/api/login"
	logoutURL   = "/api/logout"
)

func GetRouter(cfg *config.Config) *httprouter.Router {
	router := httprouter.New()

	tradeHandler := handlerapi.NewTradeHandler()
	itemHandler := handlerapi.NewItemHandler()
	userHandler := handlerapi.NewUserHandler()
	authHandler := handlerauth.NewAuthHandler()

	router.GET(itemtradesURL, tradeHandler.GetTradesByItemUUID)
	router.GET(tradesURL, tradeHandler.GetTradeList)
	router.POST(tradesURL, tradeHandler.CreateTrade)
	router.DELETE(tradeURL, tradeHandler.DeleteTradeByUUID)
	router.GET(tradeURL, tradeHandler.GetTradeByTradeUUID)
	router.PUT(tradeURL, tradeHandler.UpdateTradeByUUID)
	router.GET(usertradesURL, tradeHandler.GetTradesByUserUUID)

	router.GET(itemsURL, middleware.AuthMiddleware(itemHandler.GetItemList, logging.GetLogger()))
	router.GET(itemURL, middleware.AuthMiddleware(itemHandler.GetItemByUUID, logging.GetLogger()))
	router.POST(itemsURL, middleware.AuthMiddleware(itemHandler.CreateItem, logging.GetLogger()))
	router.DELETE(itemURL, middleware.AuthMiddleware(itemHandler.DeleteItemByUUID, logging.GetLogger()))

	router.GET(usersURL, userHandler.GetUserList)
	router.GET(userURL, userHandler.GetUserByUUID)
	router.POST(usersURL, userHandler.CreateUser)
	router.DELETE(userURL, userHandler.DeleteUserByUUID)
	router.PUT(userURL, userHandler.UpdateUserByUUID)

	router.POST(registerURL, authHandler.RegisterUser)
	router.POST(loginURL, authHandler.LoginUser)
	router.POST(logoutURL, authHandler.LogoutUser)

	return router
}
