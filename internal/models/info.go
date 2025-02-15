package models

// InfoResponse struct for InfoResponse
type InfoResponse struct {
	Coins       int                     `json:"coins"`
	Inventory   []InfoResponseInventory `json:"inventory"`
	CoinHistory InfoResponseCoinHistory `json:"coinHistory"`
}

// InfoResponseCoinHistory struct for InfoResponseCoinHistory
type InfoResponseCoinHistory struct {
	Received []InfoResponseCoinHistoryReceived `json:"received"`
	Sent     []InfoResponseCoinHistorySent     `json:"sent"`
}

// InfoResponseCoinHistoryReceived struct for InfoResponseCoinHistoryReceived
type InfoResponseCoinHistoryReceived struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

// InfoResponseCoinHistorySent struct for InfoResponseCoinHistorySent
type InfoResponseCoinHistorySent struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

// InfoResponseInventory struct for InfoResponseInventory
type InfoResponseInventory struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

// NewInfoResponse instantiates a new InfoResponse object
func NewInfoResponse() *InfoResponse {
	this := InfoResponse{}
	return &this
}
