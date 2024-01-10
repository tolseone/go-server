package handlerTrade

import (
	"encoding/json"
	"go-server/internal/models"
	"go-server/pkg/logging"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type TradeController struct {
	logger *logging.Logger
}

func NewTradeController() *TradeController {
	return &TradeController{}
}

func (h *TradeController) GetTradesByItemUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	items, err := model.GetItemList(r.Context())
	if err != nil {
		h.logger.Errorf("failed to get items: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Marshal the items to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}
func (h *TradeController) GetTradeList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	itemID := params.ByName("uuid")

	// Call the corresponding method from your repository to fetch the item by ID
	item, err := model.LoadItem(itemID)
	if err != nil {
		h.logger.Errorf("failed to get item by UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Marshal the item to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}
func (h *TradeController) CreateTrade(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Извлекаем данные из тела запроса
	var newItem = model.NewItem() // ??? что прописывать
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Создаем пользователя с извлеченными данными

	// Подставьте сюда свою логику работы с базой данных
	if err := model.CreateItem(r.Context(), &newItem); err != nil {
		h.logger.Fatalf("failed to create item: %v", err)
		return
	}

	// Возвращаем успешный статус и информацию о созданном предмете
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}
func (h *TradeController) DeleteTradeByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	itemID := params.ByName("uuid")

	// Call the corresponding method from your repository to delete the item by ID
	if err := h.modelItem.DeleteItemByUUID(r.Context(), itemID); err != nil {
		h.logger.Errorf("failed to delete item by UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return a successful response for the deletion
	w.WriteHeader(http.StatusNoContent)
}
func (h *TradeController) GetTradeByTradeUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}
func (h *TradeController) UpdateTradeByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}
func (h *TradeController) GetTradesByUserUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}
