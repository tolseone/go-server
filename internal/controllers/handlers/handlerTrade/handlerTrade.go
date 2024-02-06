package handlerTrade

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"go-server/internal/models"
	"go-server/pkg/logging"
)

type TradeHandler struct {
	logger    *logging.Logger
	validator *validator.Validate
}

func NewTradeHandler() *TradeHandler {
	return &TradeHandler{
		logger:    logging.GetLogger(),
		validator: validator.New(),
	}
}

func (h *TradeHandler) CreateTrade(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

func (h *TradeHandler) GetTradeList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

func (h *TradeHandler) GetTradesByItemUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	itemIDStr := params.ByName("uuid")

	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		h.logger.Errorf("failed to parse itemID: %v", err)
		http.Error(w, "Invalid itemID", http.StatusBadRequest)
		return
	}

	trades, err := model.LoadTradesByItemUUID(itemID.String())
	if err != nil {
		h.logger.Errorf("failed to get trades by item UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trades)
}

func (h *TradeHandler) DeleteTradeByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tradeIDStr := params.ByName("uuid")

	tradeID, err := uuid.Parse(tradeIDStr)
	if err != nil {
		h.logger.Errorf("failed to parse tradeID: %v", err)
		http.Error(w, "Invalid TradeID", http.StatusBadRequest)
		return
	}

	if err := model.DeleteTradeByID(tradeID.String()); err != nil {
		h.logger.Errorf("failed to delete trade by ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TradeHandler) GetTradeByTradeUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tradeIDStr := params.ByName("uuid")

	tradeID, err := uuid.Parse(tradeIDStr)
	if err != nil {
		h.logger.Errorf("failed to parse tradeID: %v", err)
		http.Error(w, "Invalid TradeID", http.StatusBadRequest)
		return
	}

	trade, err := model.LoadTradeByID(tradeID.String())
	if err != nil {
		h.logger.Errorf("failed to get trade by ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trade)
}

func (h *TradeHandler) UpdateTradeByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tradeID := params.ByName("uuid")

	var updateData *model.Trade
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		h.logger.Errorf("failed to decode update data: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	_, err := updateData.Save()
	if err != nil {
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

func (h *TradeHandler) GetTradesByUserUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userIDStr := params.ByName("uuid")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Errorf("failed to parse tradeID: %v", err)
		http.Error(w, "Invalid TradeID", http.StatusBadRequest)
		return
	}

	trades, err := model.LoadTradesByUserUUID(userID.String())
	if err != nil {
		h.logger.Errorf("failed to get trades by user UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trades)
}
