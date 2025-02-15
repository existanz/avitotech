package service

import (
	"avitotech/internal/database"
	"avitotech/internal/models"
)

type InfoService interface {
	GetInfo(userId int) (*models.InfoResponse, error)
}
type infoService struct {
	db database.Service
}

func NewInfoService(db database.Service) *infoService {
	return &infoService{
		db: db,
	}
}

func (s *infoService) GetInfo(userId int) (*models.InfoResponse, error) {
	response := models.NewInfoResponse()

	// Получаем количество монет
	coins, err := s.db.GetCoinsByUserID(userId)
	if err != nil {
		return nil, err
	}
	response.Coins = coins

	// Получаем инвентарь
	inventoryItems, err := s.db.GetInventoryByUserID(userId)
	if err != nil {
		return nil, err
	}
	for _, item := range inventoryItems {
		response.Inventory = append(response.Inventory, models.InfoResponseInventory{
			Type:     item.ItemType,
			Quantity: item.Quantity,
		})
	}

	// Получаем историю транзакций
	transactions, err := s.db.GetTransactionsByUserID(userId)
	if err != nil {
		return nil, err
	}
	for _, transaction := range transactions {
		if transaction.ToUserID == userId {
			response.CoinHistory.Received = append(response.CoinHistory.Received, models.InfoResponseCoinHistoryReceived{
				FromUser: s.db.GetUserNameById(transaction.FromUserID),
				Amount:   transaction.Amount,
			})
		}
		if transaction.FromUserID == userId {
			response.CoinHistory.Sent = append(response.CoinHistory.Sent, models.InfoResponseCoinHistorySent{
				ToUser: s.db.GetUserNameById(transaction.ToUserID),
				Amount: transaction.Amount,
			})
		}
	}

	return response, nil
}
