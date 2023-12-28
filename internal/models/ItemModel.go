package model

import (
	"context"
	"errors"
	"go-server/internal/repositories/db/postgresItem"

	"github.com/google/uuid"
)

type Item struct {
	ItemId      uuid.UUID `json:"item_id"`
	Name        string    `json:"name"`
	Rarity      string    `json:"rarity"`
	Description string    `json:"description,omitempty"`
}

type ModelItem struct {
	repositoryItem db.RepositoryItem
}

func NewModelItem(repo db.RepositoryItem) *ModelItem {
	return &ModelItem{repositoryItem: repo}
}

func (m *ModelItem) CreateItem(ctx context.Context, item Item) error {
	return m.repositoryItem.Create(ctx, item)
}

func (m *ModelItem) GetItemList(ctx context.Context) ([]Item, error) {
	entities, err := m.repositoryItem.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(entities))
	for i, entity := range entities {
		items[i] = entity.(Item)
	}

	return items, nil
}

func (m *ModelItem) GetItemByUUID(ctx context.Context, uuid string) (*Item, error) {
	entity, err := m.repositoryItem.FindOne(ctx, uuid)
	if err != nil {
		return nil, err
	}

	item, ok := entity.(Item)
	if !ok {
		return nil, errors.New("invalid entity type")
	}

	return &item, nil
}

func (m *ModelItem) UpdateItem(ctx context.Context, updatedItem *Item) error {
	return m.repositoryItem.Update(ctx, updatedItem)
}

func (m *ModelItem) DeleteItemByUUID(ctx context.Context, uuid string) error {
	return m.repositoryItem.Delete(ctx, uuid)
}
