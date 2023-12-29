package router

import (
	"github.com/julienschmidt/httprouter"
	"go-server/internal/controllers/handlers/handlerItem"
	"go-server/internal/controllers/handlers/handlerTrade"
	"go-server/internal/controllers/handlers/handlerUser"
)

const (
	tradesURL = "/api/trades"
	tradeURL  = "/api/trades/:uuid"

	itemsURL = "/api/items"
	itemURL  = "/api/items/:uuid"

	usersURL = "/api/users"
	userURL  = "/api/users/:uuid"
)

func getRouter() *httprouter.Router {
	router := httprouter.New()

	tradeController := handlerTrade.NewTradeController()
	itemController := handlerItem.NewItemController()
	userController := handlerUser.NewUserController()

	router.GET(tradesURL, tradeController.GettradeList)
	router.GET(tradeURL, tradeController.GettradeByUUID)
	router.POST(tradesURL, tradeController.Createtrade)
	router.DELETE(tradeURL, tradeController.DeletetradeByUUID)

	router.GET(itemsURL, itemController.GetItemList)
	router.GET(itemURL, itemController.GetItemByUUID)
	router.POST(itemsURL, itemController.CreateItem)
	router.DELETE(itemURL, itemController.DeleteItemByUUID)

	router.GET(usersURL, userController.GetUserList)
	router.GET(userURL, userController.GetUserByUUID)
	router.POST(usersURL, userController.CreateUser)
	router.DELETE(userURL, userController.DeleteUserByUUID)
	router.PUT(userURL, userController.UpdateUserByUUID)

	return router
}
