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

	transfer, arg, err := createRandomTransfer(account1.ID, account2.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.DestinationAccountID, transfer.DestinationAccountID)
	require.Equal(t, arg.SourceAccountID, transfer.SourceAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

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

	arg := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomMoney(),
	}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.DestinationAccountID, transfer2.DestinationAccountID)
	require.Equal(t, transfer1.SourceAccountID, transfer2.SourceAccountID)
	require.Equal(t, arg.Amount, transfer2.Amount)
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

func TestGetTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		account1, _, _ := createRandomAccount()
		account2, _, _ := createRandomAccount()
		_, _, err := createRandomTransfer(account1.ID, account2.ID)
		require.NoError(t, err)
	}

	arg := GetTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.GetTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestGetOutboundTransfersForAccount(t *testing.T) {
	account1, _, _ := createRandomAccount()

	for i := 0; i < 10; i++ {
		account2, _, _ := createRandomAccount()
		_, _, err := createRandomTransfer(account1.ID, account2.ID)
		require.NoError(t, err)
	}

	arg := GetOutboundTransfersForAccountParams{
		SourceAccountID: account1.ID,
		Limit:           5,
		Offset:          5,
	}

	transfers, err := testQueries.GetOutboundTransfersForAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.SourceAccountID)
	}
}

func TestGetInboundTransfersForAccount(t *testing.T) {
	account2, _, _ := createRandomAccount()

	for i := 0; i < 10; i++ {
		account1, _, _ := createRandomAccount()
		_, _, err := createRandomTransfer(account1.ID, account2.ID)
		require.NoError(t, err)
	}

	arg := GetInboundTransfersForAccountParams{
		DestinationAccountID: account2.ID,
		Limit:                5,
		Offset:               5,
	}

	transfers, err := testQueries.GetInboundTransfersForAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, account2.ID, transfer.DestinationAccountID)
	}
}
