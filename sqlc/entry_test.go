package db

import (
	"context"
	"database/sql"
	"github.com/CrunchyBlue/Golang-Bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	account, _, _ := createRandomAccount()

	entry, arg, err := createRandomEntry(account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntry(t *testing.T) {
	account, _, _ := createRandomAccount()

	entry1, _, _ := createRandomEntry(account.ID)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	account, _, _ := createRandomAccount()

	entry1, _, _ := createRandomEntry(account.ID)

	arg := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: util.RandomInt(0, 1000),
	}

	entry2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, arg.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	account, _, _ := createRandomAccount()

	entry1, _, _ := createRandomEntry(account.ID)

	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestGetEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		account, _, _ := createRandomAccount()
		_, _, err := createRandomEntry(account.ID)
		require.NoError(t, err)
	}

	arg := GetEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.GetEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestGetEntriesForAccount(t *testing.T) {
	account, _, _ := createRandomAccount()

	for i := 0; i < 10; i++ {
		_, _, err := createRandomEntry(account.ID)
		require.NoError(t, err)
	}

	arg := GetEntriesForAccountParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.GetEntriesForAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
