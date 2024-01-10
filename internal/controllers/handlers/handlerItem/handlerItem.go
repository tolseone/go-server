package handlerItem

import (
	"encoding/json"
	"go-server/internal/models"
	"go-server/pkg/logging"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type ItemController struct {
	logger *logging.Logger
}

func NewItemController() *ItemController {
	return &ItemController{
		logger: logging.GetLogger(),
	}
}

func (h *ItemController) GetItemList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	items, err := model.LoadItems()
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
func (h *ItemController) GetItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
func (h *ItemController) CreateItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Извлекаем данные из тела запроса
	itemName := params.ByName("Name")
	itemRarity := params.ByName("Rarity")
	itemDescription := params.ByName("Description")

	var newItem = model.NewItem(itemName, itemRarity, itemDescription)
	if newItem == nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Подставьте сюда свою логику работы с базой данных
	id, err := newItem.Save()
	if err != nil {
		h.logger.Fatalf("failed to create item: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	newItem.ItemId = id.(uuid.UUID)

	// Возвращаем успешный статус и информацию о созданном предмете
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}
func (h *ItemController) DeleteItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	itemID := params.ByName("uuid")

	// Call the corresponding method from your repository to delete the item by ID
	if err := model.DeleteItem(itemID); err != nil {
		h.logger.Errorf("failed to delete item by UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return a successful response for the deletion
	w.WriteHeader(http.StatusNoContent)
}
