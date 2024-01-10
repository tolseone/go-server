package model

type Trade struct {
	TradeId       string `json:"trade_id"`
	Status        string `json:"status,omitempty"`
	OfferedItem   *Item  `json:"offered_item,omitempty"`
	RequestedItem *Item  `json:"requested_item,omitempty"`
}
