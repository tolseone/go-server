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
		TradeID:        uuid.New(),
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
	repo := db.NewRepositoryTrade(logger)
	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	return repo.Create(context.TODO(), data)
}
