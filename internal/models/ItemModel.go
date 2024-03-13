package model

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"go-server/internal/repositories/db"
	"go-server/pkg/logging"

)

type Item struct {
	ItemId  uuid.UUID `json:"item_id"`
	Name    string    `json:"name" validate:"required,min=3,max=100"`
	Rarity  string    `json:"rarity" validate:"required,min=3,max=20"`
	Quality string    `json:"quality,omitempty" validate:"required,min=3,max=1000"`
}

func (itm *Item) Save() (interface{}, error) {
	logger := logging.GetLogger()
	repo := db.NewRepositoryItem(logger)

	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	var data db.ItemData
	data.Name = itm.Name
	data.Rarity = itm.Rarity
	data.Quality = itm.Quality

	if itm.ItemId != uuid.Nil {
		return repo.Update(context.TODO(), data)
	} else {
		return repo.Create(context.TODO(), data)
	}
}

func NewItem(name, rarity, quality string) *Item {
	return &Item{
		Name:    name,
		Rarity:  rarity,
		Quality: quality,
	}
}
func LoadItem(id string) (*Item, error) {
	logger := logging.GetLogger()
	repo := db.NewRepositoryItem(logger)

	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	data, err := repo.FindOne(context.TODO(), id)
	if err != nil {
		logger.Infof("Failed to load item: %v", err)
		return &Item{}, err
	}
	return &Item{
		data.ItemId,
		data.Name,
		data.Rarity,
		data.Quality,
	}, nil

}

func LoadItems() ([]*Item, error) {
	logger := logging.GetLogger()
	repo := db.NewRepositoryItem(logger)

	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

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
			itm.Quality,
		})
	}
	return itms, nil

}

func DeleteItem(id string) error {
	logger := logging.GetLogger()
	repo := db.NewRepositoryItem(logger)

	if repo == nil {
		return fmt.Errorf("failed to create repository")
	}

	if err := repo.Delete(context.TODO(), id); err != nil {
		logger.Infof("Failed to delete item: %v", err)
		return err
	}
	return nil
}

type ItemsResponse struct {
	Success   bool                  `json:"success"`
	Currency  string                `json:"currency"`
	Timestamp int64                 `json:"timestamp"`
	ItemsList map[string]ItemDetail `json:"items_list"`
}

type ItemDetail struct {
	Name          string    `json:"name"`
	Marketable    int       `json:"marketable"`
	Tradable      int       `json:"tradable"`
	ClassID       string    `json:"classid"`
	IconURL       string    `json:"icon_url"`
	Type          string    `json:"type"`
	Rarity        string    `json:"rarity"`
	RarityColor   string    `json:"rarity_color"`
	Quality       string    `json:"quality"`
	QualityColor  string    `json:"quality_color"`
	Price         ItemPrice `json:"price"`
	FirstSaleDate string    `json:"first_sale_date"`
}

type ItemPrice struct {
	Hours24 PriceDetail `json:"24_hours"`
	Days7   PriceDetail `json:"7_days"`
	Days30  PriceDetail `json:"30_days"`
	AllTime PriceDetail `json:"all_time"`
	OPSKins float64     `json:"opskins_average"`
}

type PriceDetail struct {
	Average      float64 `json:"average"`
	Median       float64 `json:"median"`
	Sold         string  `json:"sold"`
	StandardDev  string  `json:"standard_deviation"`
	LowestPrice  float64 `json:"lowest_price"`
	HighestPrice float64 `json:"highest_price"`
}

type ItemPartial struct {
	Name    string `json:"name"`
	Rarity  string `json:"rarity"`
	Quality string `json:"quality"`
}