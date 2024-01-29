package model

import (
	"context"
	"fmt"
	"go-server/internal/repositories/db/postgresUser"
	"go-server/pkg/logging"

	"github.com/google/uuid"
)

type User struct {
	UserId   uuid.UUID `json:"user_id"`
	Username string    `json:"username,omitempty" validate:"required,min=3,max=100"`
	Email    string    `json:"email" validate:"required,email,min=6,max=100"`
}

func (usr *User) Save() (interface{}, error) {
	var data db.UserData
	data.Username = usr.Username
	data.Email = usr.Email

	logger := logging.GetLogger()
	repo := db.NewRepository(logger)
	if repo == nil {
		return nil, fmt.Errorf("failed to create repository")
	}

	if usr.UserId != uuid.Nil {
		data.UserId = usr.UserId
		return repo.Update(context.TODO(), data)
	} else {
		return repo.Create(context.TODO(), data)
	}
}

func NewUser(Username, Email string) *User {
	usr := new(User)
	usr.Username = Username
	usr.Email = Email
	return usr
}
func LoadUser(id string) (*User, error) {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)
	data, err := repo.FindOne(context.TODO(), id)
	if err != nil {
		logger.Infof("Failed to load User: %v", err)
		return &User{}, err
	}
	return &User{
		data.UserId,
		data.Username,
		data.Email,
	}, nil

}

func LoadUsers() ([]*User, error) {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)
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
		})
	}
	return usrs, nil

}
func LoadUserByEmail(email string) (*User, error) {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)
	data, err := repo.FindUserByEmail(context.TODO(), email)
	if err != nil {
		logger.Infof("Failed to load User by email: %v", err)
		return &User{}, err
	}
	return &User{
		data.UserId,
		data.Username,
		data.Email,
	}, nil
}

func DeleteUser(id string) error {
	logger := logging.GetLogger()
	repo := db.NewRepository(logger)
	if err := repo.Delete(context.TODO(), id); err != nil {
		logger.Infof("Failed to delete User: %v", err)
		return err
	}
	return nil
}
