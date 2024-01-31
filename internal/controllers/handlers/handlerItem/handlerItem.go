package handlerItem

import (
	"go-server/internal/models"
	"go-server/pkg/logging"

	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type ItemController struct {
	logger    *logging.Logger
	validator *validator.Validate
}

func NewItemController() *ItemController {
	return &ItemController{
		logger:    logging.GetLogger(),
		validator: validator.New(),
	}
}

func (h *ItemController) GetItemList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	items, err := model.LoadItems()
	if err != nil {
		h.logger.Errorf("failed to get items: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	itemJSON, err := json.Marshal(items)
	if err != nil {
		h.logger.Errorf("ошибка при преобразовании пользователей в JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(itemJSON)
}

func (h *ItemController) GetItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	itemID := params.ByName("uuid")

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
	var newItem *model.Item

	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(newItem); err != nil {
		errors := err.(validator.ValidationErrors)
		for _, e := range errors {
			h.logger.Errorf("Validation error: %s", e)
		}
		http.Error(w, "Validation Error", http.StatusBadRequest)
		return
	}

	id, err := newItem.Save()
	if err != nil {
		h.logger.Fatalf("failed to create item: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	newItem.ItemId = id.(uuid.UUID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

func (h *ItemController) DeleteItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	itemID := params.ByName("uuid")

	if err := model.DeleteItem(itemID); err != nil {
		h.logger.Errorf("failed to delete item by UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Удаление предмета с UUID " + itemID + " прошло успешно"))
}
