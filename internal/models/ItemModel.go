package model

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-server/internal/repositories/db/postgresItem"
	"go-server/pkg/logging"
)

type Item struct {
	ItemId      uuid.UUID `json:"item_id"`
	Name        string    `json:"name"`
	Rarity      string    `json:"rarity"`
	Description string    `json:"description,omitempty"`
}

func (itm *Item) Save() (interface{}, error) {
	var data db.ItemData
	data.Name = itm.Name
	data.Rarity = itm.Rarity
	data.Description = itm.Description

	logger := logging.GetLogger()
	repo := db.NewRepository(logger)
	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	if itm.ItemId != uuid.Nil {
		return repo.Update(context.TODO(), data)
	} else {
		return repo.Create(context.TODO(), data)
	}
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
	repo := db.NewRepository(logger)
	data, err := repo.FindOne(context.TODO(), id)
	if err != nil {
		logger.Infof("Failed to load item: %v", err)
		return &Item{}, err
	}
	return &Item{
		data.ItemId,
		data.Name,
		data.Rarity,
		data.Description,
	}, nil

}

func LoadItems() ([]*Item, error) {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)
	data, err := repo.FindAll(context.TODO())
	if err != nil {
		logger.Infof("Failed to load items: %v", err)
		return []*Item{}, err
	}

	var itms []*Item
	for _, itm := range data {
		itms = append(itms, &Item{
			itm.ItemId,
			itm.Name,
			itm.Rarity,
			itm.Description,
		})
	}
	return itms, nil

}

func DeleteItem(id string) error {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)
	if err := repo.Delete(context.TODO(), id); err != nil {
		logger.Infof("Failed to delete item: %v", err)
		return err
	}
	return nil
}
