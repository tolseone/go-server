package item

type Item struct {
	ItemId      string `json:"item_id"`
	Name        string `json:"name"`
	Rarity      string `json:"rarity"`
	Description string `json:"description,omitempty"`
}
