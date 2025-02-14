package server

import (
	"avitotech/internal/customErrors"
	"avitotech/internal/models"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // My frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/", s.BaseHandler)
	r.POST("api/auth", s.AuthHandler)

	return r
}

func (s *Server) BaseHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Handler is working"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) AuthHandler(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResp := models.NewErrorResponse(customErrors.ErrInvalidRequest)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	resp, err := s.authService.Authenticate(&req)
	if err != nil {
		slog.Error("Auth handling Error:", err)
		if errors.Is(err, customErrors.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(customErrors.ErrInvalidCredentials))
		} else {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(customErrors.ErrISE))
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}
