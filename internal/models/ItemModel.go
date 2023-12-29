package model

import (
	"context"
	"errors"
	"go-server/internal/repositories/db/postgresItem"
	"go-server/pkg/logging"

	"github.com/google/uuid"
)

type Item struct {
	ItemId      uuid.UUID `json:"item_id"`
	Name        string    `json:"name"`
	Rarity      string    `json:"rarity"`
	Description string    `json:"description,omitempty"`
}

func NewItem(Name, Rarity, Description string) *Item {
	itm := new(Item)
	itm.Name = Name
	itm.Rarity = Rarity
	itm.Description = Description
	return itm
}
func LoadItem(id string) (*Item, error) {
	logger := logging.GetLogger()
	repo := db.NewRepository(db.RepositoryItem.Client, logger)
	data, err := repo.FindOne(context.TODO(), id)
	if err != nil {
		logger.Infof("Failed to load item: %v", err)
		return &Item{}, err
	}

	itm := NewItem(data.(Item).Name, data.(Item).Rarity, data.(Item).Description)
	return itm, nil

}
func (itm *Item) save() {
	var data *db.ItemData
	data.Name = itm.Name
	data.Rarity = itm.Rarity
	data.Description = itm.Description

	logger := logging.GetLogger()
	repo := db.NewRepository(db.RepositoryItem.Client, logger)
	if itm.ItemId != uuid.Nil {
		repo.Update(context.TODO(), data)
	} else {
		repo.Create(context.TODO(), data)
	}
}
