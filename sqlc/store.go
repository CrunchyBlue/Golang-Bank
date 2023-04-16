package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}

	return tx.Commit()
}

type TransferTxParams struct {
	SourceAccountID      int64 `json:"source_account_id"`
	DestinationAccountID int64 `json:"destination_account_id"`
	Amount               int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer           Transfer `json:"transfer"`
	SourceAccount      Account  `json:"source_account"`
	DestinationAccount Account  `json:"destination_account"`
	FromEntry          Entry    `json:"from_entry"`
	ToEntry            Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(
		ctx, func(q *Queries) error {
			var err error

			result.Transfer, err = q.CreateTransfer(
				ctx, CreateTransferParams{
					SourceAccountID:      params.SourceAccountID,
					DestinationAccountID: params.DestinationAccountID,
					Amount:               params.Amount,
				},
			)
			if err != nil {
				return err
			}

			result.FromEntry, err = q.CreateEntry(
				ctx, CreateEntryParams{
					AccountID: params.SourceAccountID,
					Amount:    -params.Amount,
				},
			)
			if err != nil {
				return err
			}

			result.ToEntry, err = q.CreateEntry(
				ctx, CreateEntryParams{
					AccountID: params.DestinationAccountID,
					Amount:    params.Amount,
				},
			)
			if err != nil {
				return err
			}

			// Order transaction queries to prevent deadlock
			if params.SourceAccountID < params.DestinationAccountID {
				result.SourceAccount, result.DestinationAccount, err = transfer(
					ctx, q, params.SourceAccountID, params.DestinationAccountID, params.Amount,
				)
			} else {
				result.DestinationAccount, result.SourceAccount, err = transfer(
					ctx, q, params.DestinationAccountID, params.SourceAccountID, -params.Amount,
				)
			}

			return nil
		},
	)

	return result, err
}

func transfer(
	ctx context.Context, q *Queries, firstAccountID int64, secondAccountID int64, amount int64,
) (firstAccount Account, secondAccount Account, err error) {
	firstAccount, err = q.UpdateAccountBalance(
		ctx, UpdateAccountBalanceParams{
			ID:     firstAccountID,
			Amount: -amount,
		},
	)
	if err != nil {
		return
	}

	secondAccount, err = q.UpdateAccountBalance(
		ctx, UpdateAccountBalanceParams{
			ID:     secondAccountID,
			Amount: amount,
		},
	)

	return
}
