package entities

type InventoryItem struct {
	ItemType string `json:"item_type,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}
