package models

import "time"

type InfoResponse struct {
	Coins        int64             `json:"coins"`
	Inventory    []InventoryItem   `json:"inventory"`
	Transactions []TransactionItem `json:"transactions"`
}

type InventoryItem struct {
	ItemType string `json:"item_type"`
	Quantity int64  `json:"quantity"`
}

type TransactionItem struct {
	FromUserID   string `json:"from_user"`
	FromUsername string
	ToUserName   string
	ToUserID     string
	Amount       int64     `json:"amount"`
	CreatedAt    time.Time `json:"created_at"`
}
