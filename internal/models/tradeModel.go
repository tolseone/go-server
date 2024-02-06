package model

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"go-server/internal/repositories/db/postgresTrade"
	"go-server/pkg/logging"

)

type Trade struct {
	TradeID        uuid.UUID `json:"trade_id"`
	UserID         uuid.UUID `json:"user_id" validate:"required"`
	Status         string    `json:"status"`
	Date           time.Time `json:"date"`
	OfferedItems   []*Item   `json:"offered_items" validate:"required"`
	RequestedItems []*Item   `json:"requested_items" validate:"required"`
}

// TradeItem is structure of item in trade.
type TradeItem struct {
	ID         uuid.UUID `json:"id"`
	TradeID    uuid.UUID `json:"trade_id"`
	ItemID     uuid.UUID `json:"item_id"`
	ItemStatus string    `json:"item_status"` // can be "offered" или "requested"
}

func NewTrade(userID uuid.UUID, offeredItems, requestedItems []*Item) *Trade {
	return &Trade{
		UserID:         userID,
		Status:         "pending",
		OfferedItems:   offeredItems,
		RequestedItems: requestedItems,
	}
}

func (t *Trade) Save() (interface{}, error) {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)

	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	var data db.TradeData
	data.UserID = t.UserID
	data.Status = t.Status
	data.Date = t.Date
	data.OfferedItems = make([]db.TradeItem, len(t.OfferedItems))
	data.RequestedItems = make([]db.TradeItem, len(t.RequestedItems))

	for i, item := range t.OfferedItems {
		data.OfferedItems[i] = db.TradeItem{
			ItemID:     item.ItemId,
			ItemStatus: "offered",
		}
	}

	for i, item := range t.RequestedItems {
		data.RequestedItems[i] = db.TradeItem{
			ItemID:     item.ItemId,
			ItemStatus: "requested",
		}
	}

	if t.TradeID != uuid.Nil {
		return repo.Update(context.TODO(), data)
	} else {
		return repo.Create(context.TODO(), data)
	}
}

func LoadTradeList() ([]*Trade, error) {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)

	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	data, err := repo.FindAll(context.TODO())
	if err != nil {
		logger.Infof("Failed to load trades: %v", err)
		return []*Trade{}, err
	}

	var trades []*Trade
	for _, tradeData := range data {
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
			tradeData.Status,
			tradeData.Date,
			offeredItems,
			requestedItems,
		})
	}
	return trades, nil
}

func LoadTradeByID(tradeID string) (*Trade, error) {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)

	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	data, err := repo.FindOne(context.TODO(), tradeID)
	if err != nil {
		logger.Infof("Failed to load trade by ID: %v", err)
		return &Trade{}, err
	}

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
		data.Status,
		data.Date,
		offeredItems,
		requestedItems,
	}, nil
}

func DeleteTradeByID(tradeID string) error {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)

	if repo == nil {
		return fmt.Errorf("failed to create repository")
	}

	if err := repo.Delete(context.TODO(), tradeID); err != nil {
		logger.Infof("Failed to delete trade: %v", err)
		return err
	}
	return nil
}

func LoadTradesByItemUUID(itemID string) ([]*Trade, error) {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)

	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	tradeData, err := repo.FindByItemUUID(context.TODO(), itemID)
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
			trade.Status,
			trade.Date,
			offeredItems,
			requestedItems,
		})
	}
	return trades, nil
}

func LoadTradesByUserUUID(userID string) ([]*Trade, error) {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)

	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	tradeData, err := repo.GetTradesByUserUUID(context.TODO(), userID)
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
			trade.Status,
			trade.Date,
			offeredItems,
			requestedItems,
		})
	}
	return trades, nil
}

func loadItems(tradeItems []db.TradeItem) ([]*Item, error) {
	logger := logging.GetLogger()
	var items []*Item

	for _, tradeItem := range tradeItems {
		item, err := LoadItem(tradeItem.ItemID.String())
		if err != nil {
			logger.Infof("Failed to load item: %v", err)
			return []*Item{}, err
		}
		items = append(items, item)
	}

	return items, nil
}
