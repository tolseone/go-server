package user

import (
	"go-server/internal/handlers"
	"go-server/pkg/logging"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var _ handlers.Handler = &handler{} // подсказка - если в методах что-то не так, то оно падает

const (
	usersURL = "/api/users"
	userURL  = "/api/users/:uuid"
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
	router.GET(usersURL, h.GetUserList)
	router.GET(userURL, h.GetUserByUUID)
	router.POST(usersURL, h.CreateUser)
	router.DELETE(userURL, h.DeleteUserByUUID)
	router.PUT(userURL, h.UpdateUserByUUID)
	router.PATCH(userURL, h.PartiallyUpdateUserByUUID)
}
func (h *handler) GetUserList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("This is list of users"))
}
func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("This is user by UUID"))
}
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(201)
	w.Write([]byte("This is creating user"))
}
func (h *handler) DeleteUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("This is delete user by UUID"))
}
func (h *handler) UpdateUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("This is update user by UUID"))
}
func (h *handler) PartiallyUpdateUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("This is partially update user by UUID"))
}
