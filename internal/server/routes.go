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

	r.POST("api/auth", s.AuthHandler)

	r.Use(AuthMiddleware(s.secretKey))

	r.GET("api/info", s.InfoHandler)
	r.POST("api/sendCoin", s.SendCoinHandler)

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

func (s *Server) InfoHandler(c *gin.Context) {
	userId, ok := c.Keys["userId"].(int)
	if !ok {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(customErrors.ErrInvalidRequest))
		return
	}
	resp, err := s.infoService.GetInfo(userId)
	if err != nil {
		slog.Error("Info handling Error:", err)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(customErrors.ErrISE))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) SendCoinHandler(c *gin.Context) {
	var req models.SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResp := models.NewErrorResponse(customErrors.ErrInvalidRequest)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	userId, ok := c.Keys["userId"].(int)
	if !ok {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(customErrors.ErrInvalidRequest))
		return
	}
	err := s.transactionService.SendCoin(userId, &req)
	if err != nil {
		slog.Error("SendCoin handling Error:", err)
		if errors.Is(err, customErrors.ErrNotEnoughCoins) {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse(customErrors.ErrNotEnoughCoins))
		} else if errors.Is(err, customErrors.ErrInvalidUsername) {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse(customErrors.ErrInvalidUsername))
		} else {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(customErrors.ErrISE))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Coins sent"})
}
