package handleradmin

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"

	"go-server/pkg/logging"

)

type AdminHandler struct {
	logger    *logging.Logger
	validator *validator.Validate
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		logger:    logging.GetLogger(),
		validator: validator.New(),
	}
}

func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) GetUserList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) GetUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) UpdateUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) UpdateUserRoleByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) DeleteUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) CreateItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) GetItemList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) GetItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) UpdateItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) DeleteItemByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) CreateTrade(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) GetTradeList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) GetTradeByTradeUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) UpdateTradeByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}

func (h *AdminHandler) DeleteTradeByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
}
