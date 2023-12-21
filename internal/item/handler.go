package item

import (
	"go-server/internal/handlers"
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
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}
func (h *handler) Reqister(router *httprouter.Router) {
	router.GET(itemsURL, h.GetItemList)
	router.GET(itemURL, h.GetItemByUUID)
	router.POST(itemsURL, h.CreateItem)
	router.DELETE(itemURL, h.DeleteItemByUUID)
}
func (h *handler) GetItemList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("This is list of items"))
}
func (h *handler) GetItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("This is item by UUID"))
}
func (h *handler) CreateItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(201)
	w.Write([]byte("This is creating item"))
}
func (h *handler) DeleteItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("This is delete item by UUID"))
}
