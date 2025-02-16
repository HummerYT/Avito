package send_coin

import "context"

type sender interface {
	SendCoin(ctx context.Context, fromUser, toUser string, amount int64) error
}
