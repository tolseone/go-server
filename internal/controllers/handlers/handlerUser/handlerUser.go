package handler

import (
	"encoding/json"
	"go-server/internal/controllers"
	model "go-server/internal/models"
	"go-server/internal/repositories"
	"go-server/pkg/logging"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

var _ handlers.Handler = &handler{} // подсказка - если в методах что-то не так, то оно падает

const (
	usersURL = "/api/users"
	userURL  = "/api/users/:uuid"
)

type handler struct {
	logger   *logging.Logger
	repoUser storage.Repository
}

func NewHandler(logger *logging.Logger, repoUser storage.Repository) handlers.Handler {
	return &handler{
		logger:   logger,
		repoUser: repoUser,
	}
}
func (h *handler) Reqister(router *httprouter.Router) {
	router.GET(usersURL, h.GetUserList)
	router.GET(userURL, h.GetUserByUUID)
	router.POST(usersURL, h.CreateUser)
	router.DELETE(userURL, h.DeleteUserByUUID)
	router.PUT(userURL, h.UpdateUserByUUID)
}
func (h *handler) GetUserList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Здесь вы можете вызвать соответствующий метод из вашего хранилища для получения списка пользователей.
	users, err := h.repoUser.FindAll(r.Context())
	if err != nil {
		h.logger.Errorf("ошибка при получении списка пользователей: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Преобразуйте пользователей в JSON и отправьте ответ.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Извлеките идентификатор пользователя из параметров URL.
	userID := params.ByName("uuid")

	// Вызовите соответствующий метод из вашего хранилища для получения пользователя по идентификатору.
	user, err := h.repoUser.FindOne(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("ошибка при получении пользователя по UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Преобразуйте пользователя в JSON и отправьте ответ.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Раскодируйте тело запроса в нового пользователя.
	var newUser model.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		h.logger.Errorf("ошибка при декодировании тела запроса: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Вызовите соответствующий метод из вашего хранилища для создания пользователя.
	if err := h.repoUser.Create(r.Context(), &newUser); err != nil {
		h.logger.Errorf("ошибка при создании пользователя: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Верните созданного пользователя в ответе.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
func (h *handler) DeleteUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Извлеките идентификатор пользователя из параметров URL.
	userID := params.ByName("uuid")

	// Вызовите соответствующий метод из вашего хранилища для удаления пользователя по идентификатору.
	if err := h.repoUser.Delete(r.Context(), userID); err != nil {
		h.logger.Errorf("ошибка при удалении пользователя по UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Верните успешный ответ об удалении.
	w.WriteHeader(http.StatusNoContent)
}
func (h *handler) UpdateUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Извлеките идентификатор пользователя из параметров URL.
	userID := params.ByName("uuid")

	// Раскодируйте тело запроса в обновленного пользователя.
	var updatedUser model.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		h.logger.Errorf("ошибка при декодировании тела запроса: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Парсим идентификатор пользователя из строки в формате UUID.
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		h.logger.Errorf("ошибка при парсинге UUID пользователя: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Устанавливаем идентификатор пользователя в соответствии с параметром URL.
	updatedUser.UserId = parsedUserID

	// Вызовите соответствующий метод из вашего хранилища для обновления пользователя.
	if err := h.repoUser.Update(r.Context(), &updatedUser); err != nil {
		h.logger.Errorf("ошибка при обновлении пользователя по UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Верните успешный ответ об обновлении.
	w.WriteHeader(http.StatusNoContent)
}
