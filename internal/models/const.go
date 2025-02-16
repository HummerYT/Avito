package models

import "time"

const (
	AuthorizationToken = "Authorization"
)

var (
	DurationJwtToken = time.Hour * 24

	MinEntropyBits = 50

	PriceItem = map[string]int64{
		"t-shirt":    80,
		"cup":        20,
		"book":       50,
		"pen":        10,
		"powerbank":  200,
		"hoody":      300,
		"umbrella":   200,
		"socks":      10,
		"wallet":     50,
		"pink-hoody": 500,
	}
)
