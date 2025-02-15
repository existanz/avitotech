package models

// SendCoinRequest struct for SendCoinRequest
type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}
