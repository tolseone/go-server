package model

import "github.com/google/uuid"

type Item struct {
	ItemId      uuid.UUID `json:"item_id"`
	Name        string    `json:"name"`
	Rarity      string    `json:"rarity"`
	Description string    `json:"description,omitempty"`
}
