package server

import (
	"avitotech/internal/customErrors"
	"avitotech/internal/models"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.New()
	r.Use(LoggerMiddleware())
	r.Use(gin.Recovery())
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
	r.GET("api/buy/:item", s.BuyItemHandler)

	return r
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
		slog.Error("Auth handling", "Error", err)
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
		slog.Error("Info handling", "Error", err)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(customErrors.ErrISE))
		return
	}
	if resp == nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(customErrors.ErrNotFound))
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) SendCoinHandler(c *gin.Context) {
	var req models.SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("SendCoin handling", "Error", err)
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
		slog.Error("SendCoin handling", "Error", err)
		if errors.Is(err, customErrors.ErrNotEnoughCoins) {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse(customErrors.ErrNotEnoughCoins))
		} else if errors.Is(err, customErrors.ErrInvalidUsername) {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse(customErrors.ErrInvalidUsername))
		} else {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(customErrors.ErrISE))
		}
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) BuyItemHandler(c *gin.Context) {
	itemType := c.Param("item")
	if itemType == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(fmt.Errorf("item type is required")))
		return
	}

	userId, ok := c.Keys["userId"].(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(customErrors.ErrUnauthorized))
		return
	}

	if err := s.shopService.BuyItem(userId, itemType); err != nil {
		slog.Error("BuyItem handling", "Error", err)
		if errors.Is(err, customErrors.ErrNotEnoughCoins) || errors.Is(err, customErrors.ErrNotFound) {
			c.JSON(http.StatusBadRequest, models.NewErrorResponse(err))
		} else {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(customErrors.ErrISE))
		}
		return
	}

	c.Status(http.StatusOK)
}
