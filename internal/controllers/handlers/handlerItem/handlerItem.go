package handler

import (
	"encoding/json"
	"go-server/internal/controllers"
	"go-server/internal/models"
	"go-server/internal/repositories"
	"go-server/pkg/logging"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var _ handlers.Handler = &handler{} // подсказка - если в методах что-то не так, то оно падает

const (
	itemsURL = "/api/items"
	itemURL  = "/api/items/:uuid"
)

type handler struct {
	logger   *logging.Logger
	repoItem storage.Repository
}

func NewHandler(logger *logging.Logger, repoItem storage.Repository) handlers.Handler {
	return &handler{
		logger:   logger,
		repoItem: repoItem,
	}
}
func (h *handler) Reqister(router *httprouter.Router) {
	router.GET(itemsURL, h.GetItemList)
	router.GET(itemURL, h.GetItemByUUID)
	router.POST(itemsURL, h.CreateItem)
	router.DELETE(itemURL, h.DeleteItemByUUID)
}
func (h *handler) GetItemList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	items, err := h.repoItem.FindAll(r.Context())
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
func (h *handler) GetItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	itemID := params.ByName("uuid")

	// Call the corresponding method from your repository to fetch the item by ID
	item, err := h.repoItem.FindOne(r.Context(), itemID)
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
func (h *handler) CreateItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Извлекаем данные из тела запроса
	var newItem model.Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Создаем пользователя с извлеченными данными

	// Подставьте сюда свою логику работы с базой данных
	if err := h.repoItem.Create(r.Context(), &newItem); err != nil {
		h.logger.Fatalf("failed to create item: %v", err)
		return
	}

	// Возвращаем успешный статус и информацию о созданном предмете
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}
func (h *handler) DeleteItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	itemID := params.ByName("uuid")

	// Call the corresponding method from your repository to delete the item by ID
	if err := h.repoItem.Delete(r.Context(), itemID); err != nil {
		h.logger.Errorf("failed to delete item by UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return a successful response for the deletion
	w.WriteHeader(http.StatusNoContent)
}
