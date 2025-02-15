package server

import (
	"avitotech/internal/customErrors"
	"avitotech/internal/models"
	jwt2 "avitotech/pkg/jwt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(customErrors.ErrUnauthorized))
			c.Abort()
			return
		}
		jwtParser := jwt2.NewJWTUtil(secretKey)
		userId, err := jwtParser.ParseUserIdFromToken(authHeader)
		if err != nil {
			slog.Warn("Authorization error:", err)
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(customErrors.ErrUnauthorized))
			c.Abort()
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}
