package handlerUser

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"go-server/internal/models"
	"go-server/pkg/logging"

)

type UserHandler struct {
	logger    *logging.Logger
	validator *validator.Validate
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		logger:    logging.GetLogger(),
		validator: validator.New(),
	}
}

func (h *UserHandler) GetUserList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	users, err := model.LoadUsers()
	if err != nil {
		h.logger.Errorf("ошибка при получении списка пользователей: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userJSON, err := json.Marshal(users)
	if err != nil {
		h.logger.Errorf("ошибка при преобразовании пользователей в JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}

func (h *UserHandler) GetUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("uuid")

	user, err := model.LoadUser(userID)
	if err != nil {
		h.logger.Errorf("ошибка при получении пользователя по UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var newUser *model.User

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(newUser); err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	existingUser, err := model.LoadUserByEmail(newUser.Email)
	if err == nil && existingUser != nil {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	id, err := newUser.Save()
	if err != nil {
		h.logger.Errorf("ошибка при создании пользователя: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	newUser.UserId = id.(uuid.UUID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (h *UserHandler) DeleteUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("uuid")

	if err := model.DeleteUser(userID); err != nil {
		h.logger.Errorf("ошибка при удалении пользователя по UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Удаление пользователя с UUID " + userID + " прошло успешно"))
}

func (h *UserHandler) UpdateUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("uuid")

	var updatedUser *model.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		h.logger.Errorf("ошибка при декодировании тела запроса: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(updatedUser); err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		h.logger.Errorf("ошибка при парсинге UUID пользователя: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	updatedUser.UserId = parsedUserID

	_, err = updatedUser.Save()
	if err != nil {
		h.logger.Errorf("ошибка при обновлении пользователя по UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}
