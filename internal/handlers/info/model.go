package info

import "AvitoTask/internal/models"

type Output struct {
	Coins       int64             `json:"coins"`
	Inventory   []InvOutput       `json:"inventory"`
	CoinHistory CoinHistoryOutput `json:"coinHistory"`
}

type InvOutput struct {
	Type     string `json:"type"`
	Quantity int64  `json:"quantity"`
}

type CoinHistoryOutput struct {
	Received []ReceivedItem `json:"received"`
	Sent     []SentItem     `json:"sent"`
}

type ReceivedItem struct {
	FromUser string `json:"fromUser"`
	Amount   int64  `json:"amount"`
}

type SentItem struct {
	ToUser string `json:"toUser"`
	Amount int64  `json:"amount"`
}

func ConvertInfoResponse(infoResp models.InfoResponse, currentUserID, username string) Output {
	out := Output{
		Coins:     infoResp.Coins,
		Inventory: make([]InvOutput, 0, len(infoResp.Inventory)),
		CoinHistory: CoinHistoryOutput{
			Received: make([]ReceivedItem, 0),
			Sent:     make([]SentItem, 0),
		},
	}

	for _, inv := range infoResp.Inventory {
		out.Inventory = append(out.Inventory, InvOutput{
			Type:     inv.ItemType,
			Quantity: inv.Quantity,
		})
	}

	for _, tx := range infoResp.Transactions {
		switch {
		case tx.ToUserID == currentUserID:
			out.CoinHistory.Received = append(out.CoinHistory.Received, ReceivedItem{
				FromUser: username,
				Amount:   tx.Amount,
			})

		case tx.FromUserID == currentUserID:
			out.CoinHistory.Sent = append(out.CoinHistory.Sent, SentItem{
				ToUser: tx.Username,
				Amount: tx.Amount,
			})
		}
	}

	return out
}
