package handlerItem

import (
	"encoding/json"
	"go-server/internal/models"
	"go-server/pkg/logging"
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
	// Вызов метода для получения списка предметов.
	items, err := model.LoadItems()
	if err != nil {
		h.logger.Errorf("failed to get items: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Marshal the items to JSON
	itemJSON, err := json.Marshal(items)
	if err != nil {
		h.logger.Errorf("ошибка при преобразовании пользователей в JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(itemJSON)
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
	// Создание объекта нового предмета.
	var newItem *model.Item

	// Декодирование данных из запроса в новый предмет.
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация данных, которые вводит пользователь.
	if err := h.validator.Struct(newItem); err != nil {
		errors := err.(validator.ValidationErrors)
		for _, e := range errors {
			// Выводим ошибку в лог или обрабатываем - OPTIONALLY
			h.logger.Errorf("Validation error: %s", e)
		}
		http.Error(w, "Validation Error", http.StatusBadRequest)
		return
	}

	// Вызов метода для создания предмета.
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
	// Извлекаем uuid предмета из параметров URL
	itemID := params.ByName("uuid")

	// Вызов метода для удаления предмета по идентификатору.
	if err := model.DeleteItem(itemID); err != nil {
		h.logger.Errorf("failed to delete item by UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return a successful response for the deletion
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Удаление предмета с UUID " + itemID + " прошло успешно"))
}
