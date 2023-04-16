package db

import (
	"context"
	"database/sql"
	"github.com/CrunchyBlue/Golang-Bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	account1, _, _ := createRandomAccount()
	account2, _, _ := createRandomAccount()

	transfer, params, err := createRandomTransfer(account1.ID, account2.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, params.DestinationAccountID, transfer.DestinationAccountID)
	require.Equal(t, params.SourceAccountID, transfer.SourceAccountID)
	require.Equal(t, params.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestGetTransfer(t *testing.T) {
	account1, _, _ := createRandomAccount()
	account2, _, _ := createRandomAccount()

	transfer1, _, _ := createRandomTransfer(account1.ID, account2.ID)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.DestinationAccountID, transfer2.DestinationAccountID)
	require.Equal(t, transfer1.SourceAccountID, transfer2.SourceAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	account1, _, _ := createRandomAccount()
	account2, _, _ := createRandomAccount()

	transfer1, _, _ := createRandomTransfer(account1.ID, account2.ID)

	params := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomMoney(),
	}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.DestinationAccountID, transfer2.DestinationAccountID)
	require.Equal(t, transfer1.SourceAccountID, transfer2.SourceAccountID)
	require.Equal(t, params.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	account1, _, _ := createRandomAccount()
	account2, _, _ := createRandomAccount()

	transfer1, _, _ := createRandomTransfer(account1.ID, account2.ID)

	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		account1, _, _ := createRandomAccount()
		account2, _, _ := createRandomAccount()
		createRandomTransfer(account1.ID, account2.ID)
	}

	params := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
