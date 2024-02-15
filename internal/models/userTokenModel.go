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
	UserAgent      string    `json:"user_agent"`
	UserRole       string    `json:"user_role"`
}

type TokenClaims struct {
	jwt.StandardClaims
	UserId    uuid.UUID `json:"user_id"`
	UserAgent string    `json:"user_agent"`
	UserRole  string    `json:"user_role"`
}

func (t *Token) Save() (interface{}, error) {
	logger := logging.GetLogger()
	repo := db.NewRepositoryToken(logger)

	if repo == nil {
		logger.Fatal("failed to create repository")
		return nil, errors.New("failed to create repository")
	}

	var data db.TokenData
	data.UserID = t.UserID
	data.Token = t.Token
	data.ExpirationTime = t.ExpirationTime
	data.UserAgent = t.UserAgent
	data.UserRole = t.UserRole

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

func GenerateJWT(user *User, ua string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.UserId,
		ua,
		user.Role,
	})

	tokenString, err := token.SignedString([]byte(SingingKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetTokenByUserAgent(ua string) (*Token, error) {
	logger := logging.GetLogger()
	repo := db.NewRepositoryToken(logger)

	if repo == nil {
		logger.Fatal("failed to create repository")
		return nil, errors.New("failed to create repository")
	}

	data, err := repo.GetTokenByUA(context.TODO(), ua)
	if err != nil {
		logger.Tracef("Failed to get token by user agent: %v\n", err)
		return nil, err
	}

	return &Token{
		data.ID,
		data.UserID,
		data.Token,
		data.ExpirationTime,
		data.UserAgent,
		data.UserRole,
	}, nil

}

func FindTokensByUserID(userID uuid.UUID) ([]Token, error) {
	return nil, nil
}

func DeleteTokenByUserID(userID string) error {
	logger := logging.GetLogger()
	repo := db.NewRepositoryToken(logger)

	if repo == nil {
		logger.Fatal("failed to create repository")
		return errors.New("failed to create repository")
	}

	if err := repo.DeleteTokenByUserID(context.TODO(), userID); err != nil {
		return err
	}

	return nil
}

func ParseToken(accesstoken string) (*Token, error) {
	token, err := jwt.ParseWithClaims(accesstoken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("token claims are not of type")
	}

	parsedToken := &Token{
		UserID:         claims.UserId,
		Token:          accesstoken,
		ExpirationTime: time.Unix(claims.ExpiresAt, 0),
		UserAgent:      claims.UserAgent,
		UserRole:       claims.UserRole,
	}

	return parsedToken, nil
}
