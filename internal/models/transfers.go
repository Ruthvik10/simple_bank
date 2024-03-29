package models

import "time"

type Transfer struct {
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransferStore interface {
	CreateTransfer(t *Transfer) error
}
