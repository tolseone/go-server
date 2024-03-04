package handlerapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"

	itemgrpc "go-server/internal/clients/item/grpc"
	"go-server/internal/config"
	"go-server/internal/grpc-clients"
	"go-server/internal/models"
	"go-server/pkg/logging"

)

type ItemHandler struct {
	logger            *logging.Logger
	validator         *validator.Validate
	itemServiceClient *itemgrpc.Client
}

func NewItemHandler() *ItemHandler {
	config := config.GetConfig()
	itemServiceClient, err := clients.NewItemClient(context.Background(), config)
	if err != nil {
		panic("failed to create item client: " + err.Error())
	}
	return &ItemHandler{
		logger:            logging.GetLogger(),
		validator:         validator.New(),
		itemServiceClient: itemServiceClient,
	}
}

func (h *ItemHandler) GetItemList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	const op = "handlerapi.GetItemList"

	items, err := h.itemServiceClient.GetAllItems(context.TODO())
	if err != nil {
		h.logger.Errorf("failed to get items: %s: %s", op, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	itemJSON, err := json.Marshal(items)
	if err != nil {
		h.logger.Errorf("ошибка при преобразовании пользователей в JSON: %s: %s", op, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(itemJSON)
}

func (h *ItemHandler) GetItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	itemID := params.ByName("uuid")

	item, err := h.itemServiceClient.GetItem(context.TODO(), itemID)
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

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

	id, err := h.itemServiceClient.CreateItem(context.TODO(), newItem.Name, newItem.Rarity, newItem.Description)
	if err != nil {
		h.logger.Fatalf("failed to create item: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	newItem.ItemId = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

func (h *ItemHandler) DeleteItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	itemID := params.ByName("uuid")

	_, err := h.itemServiceClient.DeleteItem(context.TODO(), itemID)
	if err != nil {
		h.logger.Errorf("failed to delete item by UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Удаление предмета с UUID " + itemID + " прошло успешно"))
}
