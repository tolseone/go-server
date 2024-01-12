package handlerTrade

import (
	"encoding/json"
	"errors"
	"go-server/internal/models"
	"go-server/pkg/logging"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type TradeController struct {
	logger    *logging.Logger
	validator *validator.Validate
}

func NewTradeController() *TradeController {
	return &TradeController{
		logger:    logging.GetLogger(),
		validator: validator.New(),
	}
}

func (h *TradeController) GetTradesByItemUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *TradeController) GetTradeList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func (h *TradeController) CreateTrade(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var trade model.Trade

	if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(trade); err != nil {
		h.logger.Errorf("validation error: %v", err)
		http.Error(w, "Validation Error", http.StatusBadRequest)
		return
	}

	// Валидация данных о предметах
	if err := h.validateItems(trade.OfferedItems); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validateItems(trade.RequestedItems); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTrade := model.NewTrade(trade.UserID, trade.OfferedItems, trade.RequestedItems)
	if _, err := newTrade.Save(); err != nil {
		h.logger.Fatalf("failed to create trade: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTrade)
}

func (h *TradeController) DeleteTradeByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *TradeController) GetTradeByTradeUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *TradeController) UpdateTradeByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *TradeController) GetTradesByUserUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *TradeController) validateItems(items []*uuid.UUID) error {
	for _, itemID := range items {
		if _, err := uuid.Parse(itemID.String()); err != nil {
			return errors.New("invalid item UUID")
		}
	}
	return nil
}
