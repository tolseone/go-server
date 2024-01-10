package handlerUser

import (
	"encoding/json"
	"go-server/internal/models"
	"go-server/pkg/logging"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type UserController struct {
	logger *logging.Logger
}

func NewUserController() *UserController {
	return &UserController{
		logger: logging.GetLogger(),
	}
}
func (h *UserController) GetUserList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Здесь вы можете вызвать соответствующий метод из вашего хранилища для получения списка пользователей.
	users, err := model.LoadUsers()
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
	// Раскодируйте тело запроса в нового пользователя.
	userName := params.ByName("Username")
	email := params.ByName("Email")

	var newUser = model.NewUser(userName, email)
	if newUser == nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вызовите соответствующий метод из вашего хранилища для создания пользователя.
	id, err := newUser.Save()
	if err != nil {
		h.logger.Errorf("ошибка при создании пользователя: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	newUser.UserId = id.(uuid.UUID)

	// Верните созданного пользователя в ответе.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
func (h *UserController) DeleteUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Извлеките идентификатор пользователя из параметров URL.
	userID := params.ByName("uuid")

	// Вызовите соответствующий метод из вашего хранилища для удаления пользователя по идентификатору.
	if err := model.DeleteUser(userID); err != nil {
		h.logger.Errorf("ошибка при удалении пользователя по UUID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Верните успешный ответ об удалении.
	w.WriteHeader(http.StatusNoContent)
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
