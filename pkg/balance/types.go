package balance

import (
	"time"
)

type User struct {
	ID      int64   `json:"id"`
	Balance float64 `json:"balance"`
}

type BalanceHistory struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	BalanceFrom float64   `json:"balance_from"`
	BalanceTo   float64   `json:"balance_to"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}
