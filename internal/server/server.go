package server

import (
	"avitotech/internal/service"
	"avitotech/pkg/jwt"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"avitotech/internal/database"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port      int
	secretKey string

	authService service.AuthService
	infoService service.InfoService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()
	secretKey := os.Getenv("JWT_SECRET")
	jwtUtil := jwt.NewJWTUtil(secretKey)
	NewServer := &Server{
		port:      port,
		secretKey: secretKey,

		authService: service.NewAuthService(db, jwtUtil),
		infoService: service.NewInfoService(db),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
