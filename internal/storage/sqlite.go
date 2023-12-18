package storage

import (
	"go-server/internal/models/itemModel"
	_ "go-server/internal/models/tradeModel"
	_ "go-server/internal/models/userModel"
	"gorm.io/gorm"
)

type SQLiteRepository struct {
	db *gorm.DB
}

func NewSQLiteRepository(db *gorm.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) GetAllItems() ([]item.Item, error) {
	var items []item.Item
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *SQLiteRepository) GetItemByID(itemID string) (*item.Item, error) {
	var item item.Item
	if err := r.db.First(&item, "item_id = ?", itemID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *SQLiteRepository) CreateItem(newItem *item.Item) error {
	return r.db.Create(newItem).Error
}

func (r *SQLiteRepository) DeleteItem(itemID string) error {
	return r.db.Delete(&item.Item{}, "item_id = ?", itemID).Error
}
