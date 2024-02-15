package model

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"go-server/internal/repositories/db"
	"go-server/pkg/logging"
)

func (usr *User) SaveByAdmin() (interface{}, error) {
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
	data.Role = usr.Role
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

func UpdateUserRole(id, role string) error {
	logger := logging.GetLogger()
	repo := db.NewRepositoryUser(logger)

	if repo == nil {
		return fmt.Errorf("failed to create repository")
	}

	return repo.UpdateUserRole(context.TODO(), id, role)
}
