/*
 * Сервис по обмену вещами Steam
 *
 * API for exchanging virtual items
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package trade

import (
	"go-server/internal/models/itemModel"
)

type Trade struct {
	TradeId       string     `json:"trade_id"`
	Status        string     `json:"status,omitempty"`
	OfferedItem   *item.Item `json:"offered_item,omitempty"`
	RequestedItem *item.Item `json:"requested_item,omitempty"`
}

type Trades struct {
	Trades []Trade `json:"trades"`
}
