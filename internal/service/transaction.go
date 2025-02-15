package service

import (
	"avitotech/internal/customErrors"
	"avitotech/internal/database"
	"avitotech/internal/models"
)

type TransactionService interface {
	SendCoin(userID int, req *models.SendCoinRequest) error
}

type transactionService struct {
	db database.Service
}

func NewTransactionService(db database.Service) *transactionService {
	return &transactionService{
		db: db,
	}
}

func (s *transactionService) SendCoin(userID int, req *models.SendCoinRequest) error {
	toUser, err := s.db.GetUserByName(req.ToUser)
	if err != nil || toUser == nil {
		return customErrors.ErrInvalidUsername
	}
	err = s.db.SendCoin(userID, toUser.ID, req.Amount)
	if err != nil {
		return err
	}
	return nil
}
