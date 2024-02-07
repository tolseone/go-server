package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"go-server/internal/repositories/db"
	"go-server/pkg/logging"
)

const (
	singingKey = "JdjJw74DFjdnbr32Aggkde"
	tokenTTL   = 12 * time.Hour
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
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
	// Load user by email from the database
	user, err := LoadUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Compare the stored hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err // Passwords don't match
	}

	return user, nil // Authentication successful
}

func GenerateJWT(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.UserId,
	})

	tokenString, err := token.SignedString([]byte(singingKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func generatePasswordHash(pw string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ParseToken(accesstoken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accesstoken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return uuid.Nil, errors.New("unexpected signing method")
		}
		return []byte(singingKey), nil

	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.Nil, errors.New("token cliams are not of type")
	}

	return claims.UserId, nil
}
