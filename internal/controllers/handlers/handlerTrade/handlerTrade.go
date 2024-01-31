package handlerTrade

import (
	"go-server/internal/models"
	"go-server/pkg/logging"

	"encoding/json"
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

func (h *TradeController) CreateTrade(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var newTrade *model.Trade

	if err := json.NewDecoder(r.Body).Decode(&newTrade); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(newTrade); err != nil {
		errors := err.(validator.ValidationErrors)
		for _, e := range errors {
			h.logger.Errorf("Validation error: %s", e)
		}
		http.Error(w, "Validation Error", http.StatusBadRequest)
		return
	}

	// Save to DB
	id, err := newTrade.Save()
	if err != nil {
		h.logger.Fatalf("failed to create trade: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	newTrade.TradeID = id.(uuid.UUID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTrade)
}

func (h *TradeController) GetTradeList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	trades, err := model.LoadTradeList()
	if err != nil {
		h.logger.Errorf("failed to get trades: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trades)
}

func (h *TradeController) GetTradesByItemUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	itemID := params.ByName("item_id")

	trades, err := model.LoadTradesByItemUUID(itemID)
	if err != nil {
		h.logger.Errorf("failed to get trades by item UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trades)
}

func (h *TradeController) DeleteTradeByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tradeID := params.ByName("uuid")

	if err := model.DeleteTradeByID(tradeID); err != nil {
		h.logger.Errorf("failed to delete trade by ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TradeController) GetTradeByTradeUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tradeID := params.ByName("trade_id")

	trade, err := model.LoadTradeByID(tradeID)
	if err != nil {
		h.logger.Errorf("failed to get trade by ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trade)
}

func (h *TradeController) UpdateTradeByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tradeID := params.ByName("trade_id")

	var updateData *model.Trade
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		h.logger.Errorf("failed to decode update data: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := model.UpdateTradeByID(tradeID, updateData.OfferedItems, updateData.RequestedItems); err != nil {
		h.logger.Errorf("failed to update trade by UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	updatedTrade, err := model.LoadTradeByID(tradeID)
	if err != nil {
		h.logger.Errorf("failed to get trade by ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTrade)
}

func (h *TradeController) GetTradesByUserUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("user_id")

	trades, err := model.LoadTradesByUserUUID(userID)
	if err != nil {
		h.logger.Errorf("failed to get trades by user UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trades)
}
