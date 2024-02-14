package model

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"go-server/internal/repositories/db"
	"go-server/pkg/logging"

)

const (
	SingingKey = "JdjJw74DFjdnbr32Aggkde"
	tokenTTL   = 12 * time.Hour
)

type Token struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Token          string    `json:"token"`
	ExpirationTime time.Time `json:"expiration_time"`
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
}

func (t *Token) Save() (interface{}, error) {
	logger := logging.GetLogger()
	repo := db.NewRepositoryToken(logger)

	if repo == nil {
		logger.Fatal("failed to create repository")
	}

	var data db.TokenData
	data.UserID = t.UserID
	data.Token = t.Token
	data.ExpirationTime = t.ExpirationTime

	if t.ID != uuid.Nil {
		data.ID = t.ID
		return repo.Update(context.TODO(), data)
	} else {
		return repo.Create(context.TODO(), data)
	}
}

func CheckAndDeleteExpiredTokens() {
	logger := logging.GetLogger()
	repo := db.NewRepositoryToken(logger)

	if repo == nil {
		logger.Fatal("failed to create repository")
	}

	now := time.Now()

	expiredTokens, err := repo.GetExpiredTokens(context.TODO(), now)
	if err != nil {
		logger.Tracef("Error getting expired tokens: %v\n", err)
	}

	for _, token := range expiredTokens {
		err := repo.DeleteToken(context.TODO(), token.ID)
		if err != nil {
			logger.Tracef("Error deleting expired token %s: %v\n", token.ID, err)
			continue
		}
		logger.Tracef("Expired token deleted: %s\n", token.ID)
	}
}

func GenerateJWT(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.UserId,
	})

	tokenString, err := token.SignedString([]byte(SingingKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func FindTokensByUserID(userID uuid.UUID) ([]Token, error) {
	return nil, nil
}

func DeleteExpiredTokens() error {
	return nil
}

func ParseToken(accesstoken string) (*Token, error) {
	token, err := jwt.ParseWithClaims(accesstoken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SingingKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return nil, errors.New("token claims are not of type")
	}

	parsedToken := &Token{
		UserID:         claims.UserId,
		Token:          accesstoken,
		ExpirationTime: time.Unix(claims.ExpiresAt, 0),
	}

	return parsedToken, nil
}
