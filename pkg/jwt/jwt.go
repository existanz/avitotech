package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTUtil struct {
	secretKey []byte
}

func NewJWTUtil(secretKey string) *JWTUtil {
	return &JWTUtil{secretKey: []byte(secretKey)}
}

func (j *JWTUtil) GenerateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // Токен действителен неделю
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}
