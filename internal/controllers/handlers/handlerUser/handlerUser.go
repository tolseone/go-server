package handlerUser

import (
	"encoding/json"
	"fmt"
	"go-server/internal/models"
	"go-server/pkg/logging"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type UserController struct {
	logger    *logging.Logger
	validator *validator.Validate
}

func NewUserController() *UserController {
	return &UserController{
		logger:    logging.GetLogger(),
		validator: validator.New(),
	}
}
func (h *UserController) GetUserList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Вызов метода для получения списка пользователей.
	users, err := model.LoadUsers()
	if err != nil {
		h.logger.Errorf("ошибка при получении списка пользователей: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Преобразование пользователей в JSON.
	userJSON, err := json.Marshal(users)
	if err != nil {
		h.logger.Errorf("ошибка при преобразовании пользователей в JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отправка ответа.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}
func (h *UserController) GetUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Извлеките идентификатор пользователя из параметров URL.
	userID := params.ByName("uuid")

	// Вызовите соответствующий метод из вашего хранилища для получения пользователя по идентификатору.
	user, err := model.LoadUser(userID)
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
func (h *UserController) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Создание объекта нового пользователя.
	var newUser *model.User

	// Декодирование данных из запроса в нового пользователя.
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация данных, которые вводит пользователь.
	if err := h.validator.Struct(newUser); err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	// Проверка на существование пользователя с такими данными.
	existingUser, err := model.LoadUserByEmail(newUser.Email)
	if err == nil && existingUser != nil {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	// Вызов метода для создания пользователя.
	id, err := newUser.Save()
	if err != nil {
		h.logger.Errorf("ошибка при создании пользователя: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	newUser.UserId = id.(uuid.UUID)

	// Возвращаем созданного пользователя в ответе.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
func (h *UserController) DeleteUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Извлекаем uuid пользователя из параметров URL
	userID := params.ByName("uuid")

	// Вызов метода для удаления пользователя по идентификатору.
	if err := model.DeleteUser(userID); err != nil {
		h.logger.Errorf("ошибка при удалении пользователя по UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Верните успешный ответ об удалении.
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Удаление пользователя с UUID " + userID + " прошло успешно"))
}
func (h *UserController) UpdateUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Извлеките идентификатор пользователя из параметров URL.
	userID := params.ByName("uuid")

	// Раскодируйте тело запроса в обновленного пользователя.
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

	// Парсим идентификатор пользователя из строки в формате UUID.
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		h.logger.Errorf("ошибка при парсинге UUID пользователя: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Устанавливаем идентификатор пользователя в соответствии с параметром URL.
	updatedUser.UserId = parsedUserID // всё дошло успешно

	// Вызовите соответствующий метод из вашего хранилища для обновления пользователя.
	_, err = updatedUser.Save()
	if err != nil {
		h.logger.Errorf("ошибка при обновлении пользователя по UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Верните успешный ответ об обновлении.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // which status should i use?
	json.NewEncoder(w).Encode(updatedUser)
}
