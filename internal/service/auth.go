package service

import (
	"avitotech/internal/customErrors"
	"avitotech/internal/database"
	"avitotech/internal/entities"
	"avitotech/internal/models"
	"avitotech/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	Authenticate(req *models.AuthRequest) (*models.AuthResponse, error)
}
type authService struct {
	db      database.Service
	jwtUtil *jwt.JWTUtil
}

func NewAuthService(db database.Service, jwtUtil *jwt.JWTUtil) *authService {
	return &authService{
		db:      db,
		jwtUtil: jwtUtil,
	}
}
func (s *authService) Authenticate(req *models.AuthRequest) (*models.AuthResponse, error) {
	user, err := s.db.GetUserByName(req.Username)
	if err != nil {
		return nil, err
	}

	// Если пользователь не существует, создаем нового
	if user == nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		user = &entities.User{
			Username:  req.Username,
			Password:  string(hashedPassword),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.db.AddUser(user); err != nil {
			return nil, err
		}
	} else {
		// Проверяем пароль существующего пользователя
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return nil, customErrors.ErrInvalidCredentials
		}
	}

	// Генерируем JWT токен
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token}, nil
}
