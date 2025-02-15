package service

import "avitotech/internal/database"

type ShopService interface {
	BuyItem(userId int, itemType string) error
}

type shopService struct {
	db database.Service
}

func NewShopService(db database.Service) *shopService {
	return &shopService{
		db: db,
	}
}

func (s *shopService) BuyItem(userId int, itemType string) error {
	err := s.db.BuyItem(userId, itemType)
	if err != nil {
		return err
	}
	return nil
}
