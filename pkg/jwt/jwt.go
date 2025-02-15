package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
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

func (j *JWTUtil) ParseUserIdFromToken(tokenString string) (int, error) {
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return -1, fmt.Errorf("bearer token is required")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil || !token.Valid {
		return -1, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return -1, fmt.Errorf("invalid token claims")
	}
	userId, ok := claims["user_id"].(float64)
	if !ok {
		return -1, fmt.Errorf("invalid user ID in token")
	}
	return int(userId), nil
}
