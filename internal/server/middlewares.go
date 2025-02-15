package server

import (
	"avitotech/internal/customErrors"
	"avitotech/internal/models"
	jwt2 "avitotech/pkg/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
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

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		slog.Info(fmt.Sprintf("--> [%s] \"%s\" [%d] %s", c.Request.Method, c.Request.URL, c.Writer.Status(), time.Since(start)))
	}
}
