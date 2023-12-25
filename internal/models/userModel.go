package model

import "github.com/google/uuid"

type User struct {
	UserId   uuid.UUID `json:"user_id"`
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email"`
}
