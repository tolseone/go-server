package model

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"go-server/internal/repositories/db"
	"go-server/pkg/logging"

)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

func (usr *User) CreateUser() (interface{}, error) {
	logger := logging.GetLogger()
	repo := db.NewRepositoryUser(logger)

	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	var data db.UserData
	var err error

	data.Username = usr.Username
	data.Email = usr.Email
	data.Password, err = generatePasswordHash(usr.Password)
	if err != nil {
		logger.Fatalf("Failed to generate password hash: %s", err.Error())
	}

	return repo.Create(context.TODO(), data)
}

func AuthenticateUser(email, password string) (*User, error) {
	user, err := LoadUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil // Authentication successful
}

func generatePasswordHash(pw string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
