// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package db

import (
	"time"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Entry struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	// Can be negative or positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Transfer struct {
	ID                   int64 `json:"id"`
	SourceAccountID      int64 `json:"source_account_id"`
	DestinationAccountID int64 `json:"destination_account_id"`
	// Must be positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}