package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

const AccessTokenTTL = time.Minute * 240

var AccessTokenSecret = []byte(os.Getenv("JWT_SECRET_TOKEN"))

type Claims struct {
	UserID    int64  `json:"userID"`
	Username  string `json:"username"`
	CompanyID uint   `json:"companyID,omitempty"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateStudentAccessToken(userID int64, username string, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(AccessTokenSecret)
}

func GenerateHrAccessToken(userID int64, username string, companyID uint, role string) (string, error) {
	claims := Claims{
		UserID:    userID,
		Username:  username,
		CompanyID: companyID,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(AccessTokenSecret)
}

func ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return AccessTokenSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
