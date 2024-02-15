package handleradmin

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

	id, err := newUser.SaveByAdmin()
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

func (h *AdminHandler) GetUserList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

func (h *AdminHandler) GetUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

func (h *AdminHandler) UpdateUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

	_, err = updatedUser.SaveByAdmin()
	if err != nil {
		h.logger.Errorf("ошибка при обновлении пользователя по UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

func (h *AdminHandler) UpdateUserRoleByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := params.ByName("uuid")

	var updatedUserRole struct {
		Role string `json:"role" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updatedUserRole); err != nil {
		h.logger.Errorf("ошибка при декодировании тела запроса: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(updatedUserRole); err != nil {
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

	err = model.UpdateUserRole(parsedUserID.String(), updatedUserRole.Role)
	if err != nil {
		h.logger.Errorf("ошибка при обновлении роли пользователя по UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Successfully updated user role")

	w.WriteHeader(http.StatusOK)
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
