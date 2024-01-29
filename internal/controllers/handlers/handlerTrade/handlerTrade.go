package handlerTrade

import (
	"encoding/json"
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

func (h *TradeController) CreateTrade(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Извлекаем данные из тела запроса
	var tradeData model.Trade
	if err := json.NewDecoder(r.Body).Decode(&tradeData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация данных
	if err := h.validator.Struct(tradeData); err != nil {
		// Обработка ошибок валидации
		errors := err.(validator.ValidationErrors)
		for _, e := range errors {
			// Выводим ошибку в лог или обрабатываем ее по вашему усмотрению
			h.logger.Errorf("Validation error: %s", e)
		}
		http.Error(w, "Validation Error", http.StatusBadRequest)
		return
	}

	// Подставьте сюда свою логику работы с базой данных
	id, err := tradeData.Save()
	if err != nil {
		h.logger.Fatalf("failed to create trade: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tradeData.TradeID = id.(uuid.UUID)

	// Возвращаем успешный статус и информацию о созданной сделке
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tradeData)
}

func (h *TradeController) GetTradeList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	trades, err := model.LoadTradeList()
	if err != nil {
		h.logger.Errorf("failed to get trades: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Marshal trades to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trades)
}

func (h *TradeController) GetTradesByItemUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Extract item ID from URL
	itemID := params.ByName("item_id")

	// Call the corresponding method from your repository to fetch trades by item ID
	trades, err := model.LoadTradesByItemUUID(itemID)
	if err != nil {
		h.logger.Errorf("failed to get trades by item UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Marshal trades to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trades)
}
func (h *TradeController) DeleteTradeByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Extract trade ID from URL
	tradeID := params.ByName("trade_id")

	// Call the corresponding method from your repository to delete the trade by ID
	if err := model.DeleteTradeByID(tradeID); err != nil {
		h.logger.Errorf("failed to delete trade by ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return a successful response for the deletion
	w.WriteHeader(http.StatusNoContent)
}

func (h *TradeController) GetTradeByTradeUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Extract trade ID from URL
	tradeID := params.ByName("trade_id")

	// Call the corresponding method from your repository to fetch the trade by ID
	trade, err := model.LoadTradeByID(tradeID)
	if err != nil {
		h.logger.Errorf("failed to get trade by ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Marshal trade to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trade)
}

func (h *TradeController) UpdateTradeByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Extract trade ID from URL
	tradeID := params.ByName("trade_id")

	// Extract data from request body
	var updateData *model.Trade
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		h.logger.Errorf("failed to decode update data: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Call the corresponding method from your repository to update the trade by ID
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

	// Marshal updated trade to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTrade)
}

func (h *TradeController) GetTradesByUserUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Extract user ID from URL
	userID := params.ByName("user_id")

	// Call the corresponding method from your repository to fetch trades by user ID
	trades, err := model.LoadTradesByUserUUID(userID)
	if err != nil {
		h.logger.Errorf("failed to get trades by user UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Marshal trades to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trades)
}
