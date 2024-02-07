package handlerauth

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"

	"go-server/internal/models"
	"go-server/pkg/logging"
)

type AuthHandler struct {
	logger    *logging.Logger
	validator *validator.Validate
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		logger:    logging.GetLogger(),
		validator: validator.New(),
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var newUser model.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := newUser.Save()
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	newUser.UserId = userID.(uuid.UUID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var credentials model.LoginInput

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := model.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		h.logger.Tracef("Failed to authenticate user: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := model.GenerateJWT(user)
	if err != nil {
		h.logger.Tracef("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	w.Write([]byte(token))
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) LogoutUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
