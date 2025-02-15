package entities

type Transaction struct {
	FromUserID int `json:"from_user_id,omitempty"`
	ToUserID   int `json:"to_user_id,omitempty"`
	Amount     int `json:"amount,omitempty"`
}
