package router

import (
	"go-server/internal/config"
	"go-server/internal/controllers/handlers/handlerItem"
	// "go-server/internal/controllers/handlers/handlerTrade"
	// "go-server/internal/controllers/handlers/handlerUser"

	"github.com/julienschmidt/httprouter"
)

const (
	// tradesURL     = "/api/trades"
	// tradeURL      = "/api/trades/:uuid"
	// usertradesURL = "/api/users/:uuid/trades"
	// itemtradesURL = "/api/items/:uuid/trades"

	itemsURL = "/api/items"
	itemURL  = "/api/items/:uuid"

	// usersURL = "/api/users"
	// userURL  = "/api/users/:uuid"
)

func GetRouter(cfg *config.Config) *httprouter.Router {
	router := httprouter.New()

	// tradeController := handlerTrade.NewTradeController()
	itemController := handlerItem.NewItemController()
	// userController := handlerUser.NewUserController()

	// router.GET(itemtradesURL, tradeController.GetTradesByItemUUID)
	// router.GET(tradesURL, tradeController.GetTradeList)
	// router.POST(tradesURL, tradeController.CreateTrade)
	// router.DELETE(tradeURL, tradeController.DeleteTradeByUUID)
	// router.GET(tradeURL, tradeController.GetTradeByTradeUUID)
	// router.PUT(tradeURL, tradeController.UpdateTradeByUUID)
	// router.GET(usertradesURL, tradeController.GetTradesByUserUUID)

	router.GET(itemsURL, itemController.GetItemList)
	router.GET(itemURL, itemController.GetItemByUUID)
	router.POST(itemsURL, itemController.CreateItem)
	router.DELETE(itemURL, itemController.DeleteItemByUUID)

	// router.GET(usersURL, userController.GetUserList)
	// router.GET(userURL, userController.GetUserByUUID)
	// router.POST(usersURL, userController.CreateUser)
	// router.DELETE(userURL, userController.DeleteUserByUUID)
	// router.PUT(userURL, userController.UpdateUserByUUID)

	return router
}
