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
	port int

	authService service.AuthService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()
	jwtUtil := jwt.NewJWTUtil(os.Getenv("JWT_SECRET"))
	NewServer := &Server{
		port: port,

		authService: service.NewAuthService(db, jwtUtil),
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
