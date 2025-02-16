package buy_item

import "context"

type buyer interface {
	BuyItem(ctx context.Context, userID, item string, cost int64) error
}
