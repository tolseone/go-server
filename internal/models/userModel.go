package model

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"go-server/internal/repositories/db"
	"go-server/pkg/logging"

)

type User struct {
	UserId   uuid.UUID `json:"user_id"`
	Username string    `json:"username,omitempty" validate:"required,min=3,max=100"`
	Email    string    `json:"email" validate:"required,email,min=6,max=100"`
	Password string    `json:"password" validate:"required,min=6,max=100"`
	Role     string    `json:"role,omitempty"`
}

func (usr *User) Save() (interface{}, error) {
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
	data.Role = "user"
	if err != nil {
		logger.Fatalf("Failed to generate password hash: %s", err.Error())
	}

	if usr.UserId != uuid.Nil {
		data.UserId = usr.UserId
		return repo.Update(context.TODO(), data)
	} else {
		return repo.Create(context.TODO(), data)
	}
}

func NewUser(username, email, password string) *User {
	return &User{
		Username: username,
		Email:    email,
		Password: password,
		Role:     "user",
	}
}

func LoadUser(id string) (*User, error) {
	logger := logging.GetLogger()
	repo := db.NewRepositoryUser(logger)

	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	data, err := repo.FindOne(context.TODO(), id)
	if err != nil {
		logger.Infof("Failed to load User: %v", err)
		return &User{}, err
	}
	return &User{
		data.UserId,
		data.Username,
		data.Email,
		data.Password,
		data.Role,
	}, nil

}

func LoadUsers() ([]*User, error) {
	logger := logging.GetLogger()
	repo := db.NewRepositoryUser(logger)

	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	data, err := repo.FindAll(context.TODO())
	if err != nil {
		logger.Infof("Failed to load Users: %v", err)
		return []*User{}, err
	}

	var usrs []*User
	for _, usr := range data {
		usrs = append(usrs, &User{
			usr.UserId,
			usr.Username,
			usr.Email,
			usr.Password,
			usr.Role,
		})
	}
	return usrs, nil

}

func LoadUserByEmail(email string) (*User, error) {
	logger := logging.GetLogger()
	repo := db.NewRepositoryUser(logger)

	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	data, err := repo.FindUserByEmail(context.TODO(), email)
	if err != nil {
		logger.Infof("Failed to load User by email: %v", err)
		return &User{}, err
	}
	return &User{
		data.UserId,
		data.Username,
		data.Email,
		data.Password,
		data.Role,
	}, nil
}

func DeleteUser(id string) error {
	logger := logging.GetLogger()
	repo := db.NewRepositoryUser(logger)

	if repo == nil {
		return fmt.Errorf("failed to create repository")
	}

	if err := repo.Delete(context.TODO(), id); err != nil {
		logger.Infof("Failed to delete User: %v", err)
		return err
	}
	return nil
}
