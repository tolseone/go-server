package storage

import (
	"go-server/internal/models/itemModel"
	_ "go-server/internal/models/tradeModel"
	_ "go-server/internal/models/userModel"
)

type ItemRepository interface {
	GetAllItems() ([]item.Item, error)
	GetItemByID(itemID string) (*item.Item, error)
	CreateItem(newItem *item.Item) error
	DeleteItemByID(itemID string) error
}

/*
type TradeRepository interface {
	GetAllTrades() ([]trade.Trade, error)
	GetTradeByID(tradeID string) (*trade.Trade, error)
	CreateTrade(newTrade *trade.Trade) error
	DeleteTrade(tradeID string) error
	UpdateTrade(tradeID string, newTrade *trade.Trade) error
	GetTradeItems(tradeID string) ([]item.Item, error)
	AddItemToTrade(tradeID string, newItem *item.Item) error
	RemoveItemFromTrade(tradeID string, itemID string) error
}

type UserRepository interface {
	GetAllUsers() ([]user.User, error)
	CreateUser(newUser *user.User) error
	GetUserByID(userID string) (*user.User, error)
	DeleteUserByID(userID string) error
	UpdateUser(userID string, newUser *user.User) error
}
*/
