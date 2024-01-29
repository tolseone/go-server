package model

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-server/internal/repositories/db/postgresTrade"
	"go-server/pkg/logging"
)

type Trade struct {
	TradeID        uuid.UUID `json:"trade_id"`
	UserID         uuid.UUID `json:"user_id"`
	OfferedItems   []*Item   `json:"offered_items"`
	RequestedItems []*Item   `json:"requested_items"`
}

func NewTrade(userID uuid.UUID, offeredItems, requestedItems []*Item) *Trade {
	return &Trade{
		UserID:         userID,
		OfferedItems:   offeredItems,
		RequestedItems: requestedItems,
	}
}

func (t *Trade) Save() (interface{}, error) {
	var data db.TradeData
	data.UserID = t.UserID
	data.OfferedItems = make([]uuid.UUID, len(t.OfferedItems))
	data.RequestedItems = make([]uuid.UUID, len(t.RequestedItems))

	for i, item := range t.OfferedItems {
		data.OfferedItems[i] = item.ItemId
	}

	for i, item := range t.RequestedItems {
		data.RequestedItems[i] = item.ItemId
	}

	logger := logging.GetLogger()
	repo := db.NewRepository(logger)
	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	return repo.Create(context.TODO(), data)
}

func LoadTradeList() ([]*Trade, error) {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)
	data, err := repo.FindAll(context.TODO())
	if err != nil {
		logger.Infof("Failed to load trades: %v", err)
		return []*Trade{}, err
	}

	var trades []*Trade
	for _, tradeData := range data {
		// Используем функцию loadItems для OfferedItems и RequestedItems
		offeredItems, err := loadItems(tradeData.OfferedItems)
		if err != nil {
			logger.Infof("Failed to load offered items for trade: %v", err)
			return []*Trade{}, err
		}

		requestedItems, err := loadItems(tradeData.RequestedItems)
		if err != nil {
			logger.Infof("Failed to load requested items for trade: %v", err)
			return []*Trade{}, err
		}

		trades = append(trades, &Trade{
			tradeData.TradeID,
			tradeData.UserID,
			offeredItems,
			requestedItems,
		})
	}
	return trades, nil
}

func LoadTradeByID(tradeID string) (*Trade, error) {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)
	data, err := repo.FindOne(context.TODO(), tradeID)
	if err != nil {
		logger.Infof("Failed to load trade by ID: %v", err)
		return &Trade{}, err
	}

	// Используем функцию loadItems для OfferedItems и RequestedItems
	offeredItems, err := loadItems(data.OfferedItems)
	if err != nil {
		logger.Infof("Failed to load offered items for trade: %v", err)
		return &Trade{}, err
	}

	requestedItems, err := loadItems(data.RequestedItems)
	if err != nil {
		logger.Infof("Failed to load requested items for trade: %v", err)
		return &Trade{}, err
	}

	return &Trade{
		data.TradeID,
		data.UserID,
		offeredItems,
		requestedItems,
	}, nil
}

func DeleteTradeByID(tradeID string) error {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)
	if err := repo.Delete(context.TODO(), tradeID); err != nil {
		logger.Infof("Failed to delete trade: %v", err)
		return err
	}
	return nil
}

func LoadTradesByItemUUID(itemID string) ([]*Trade, error) {
	logger := logging.GetLogger()
	repoTrade := db.NewRepository(logger)
	tradeData, err := repoTrade.FindByItemUUID(context.TODO(), itemID)
	if err != nil {
		logger.Infof("Failed to load trades by item UUID: %v", err)
		return []*Trade{}, err
	}

	var trades []*Trade
	for _, trade := range tradeData {
		offeredItems, err := loadItems(trade.OfferedItems)
		if err != nil {
			logger.Infof("Failed to load items for trade: %v", err)
			return []*Trade{}, err
		}

		requestedItems, err := loadItems(trade.RequestedItems)
		if err != nil {
			logger.Infof("Failed to load items for trade: %v", err)
			return []*Trade{}, err
		}

		trades = append(trades, &Trade{
			trade.TradeID,
			trade.UserID,
			offeredItems,
			requestedItems,
		})
	}
	return trades, nil
}

func UpdateTradeByID(tradeID string, offeredItems, requestedItems []*Item) error {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)

	// Преобразуем OfferedItems и RequestedItems в списки uuid.UUID
	offeredItemIDs := make([]uuid.UUID, len(offeredItems))
	requestedItemIDs := make([]uuid.UUID, len(requestedItems))

	for i, item := range offeredItems {
		offeredItemIDs[i] = item.ItemId
	}

	for i, item := range requestedItems {
		requestedItemIDs[i] = item.ItemId
	}

	// Вызываем метод Update репозитория
	if err := repo.Update(context.TODO(), tradeID, offeredItemIDs, requestedItemIDs); err != nil {
		logger.Infof("Failed to update trade by ID: %v", err)
		return err
	}

	return nil
}

func LoadTradesByUserUUID(userID string) ([]*Trade, error) {
	logger := logging.GetLogger()
	repoTrade := db.NewRepository(logger)
	tradeData, err := repoTrade.FindByUserUUID(context.TODO(), userID)
	if err != nil {
		logger.Infof("Failed to load trades by user UUID: %v", err)
		return []*Trade{}, err
	}

	var trades []*Trade
	for _, trade := range tradeData {
		offeredItems, err := loadItems(trade.OfferedItems)
		if err != nil {
			logger.Infof("Failed to load items for trade: %v", err)
			return []*Trade{}, err
		}

		requestedItems, err := loadItems(trade.RequestedItems)
		if err != nil {
			logger.Infof("Failed to load items for trade: %v", err)
			return []*Trade{}, err
		}

		trades = append(trades, &Trade{
			trade.TradeID,
			trade.UserID,
			offeredItems,
			requestedItems,
		})
	}
	return trades, nil
}

func loadItems(itemIDs []uuid.UUID) ([]*Item, error) {
	logger := logging.GetLogger()
	var items []*Item

	// Используем метод LoadItem из modelItem для каждого itemID
	for _, itemID := range itemIDs {
		item, err := LoadItem(itemID.String())
		if err != nil {
			logger.Infof("Failed to load item: %v", err)
			return []*Item{}, err
		}
		items = append(items, item)
	}

	return items, nil
}
