package api

import (
	db "github.com/CrunchyBlue/Golang-Bank/sqlc"
	"github.com/CrunchyBlue/Golang-Bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func generateMockUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

func generateMockAccounts(owner string, numAccounts int) []db.Account {
	var accounts []db.Account

	for i := 0; i < numAccounts; i++ {
		accounts = append(
			accounts, db.Account{
				ID:       util.RandomInt(1, 1000),
				Owner:    owner,
				Balance:  util.RandomInt(0, 1000),
				Currency: util.RandomCurrency(),
			},
		)
	}
	return accounts
}

func generateMockEntries(numEntries int, accountID int64) []db.Entry {
	var entries []db.Entry

	for i := 0; i < numEntries; i++ {
		entries = append(
			entries, db.Entry{
				ID:        util.RandomInt(1, 1000),
				AccountID: accountID,
				Amount:    util.RandomInt(-1000, 1000),
			},
		)
	}
	return entries
}

func generateMockTransfers(numTransfers int, sourceAccountID int64, destinationAccountID int64) []db.Transfer {
	var transfers []db.Transfer

	for i := 0; i < numTransfers; i++ {
		transfers = append(
			transfers, db.Transfer{
				ID:                   util.RandomInt(1, 1000),
				SourceAccountID:      sourceAccountID,
				DestinationAccountID: destinationAccountID,
				Amount:               util.RandomInt(0, 1000),
			},
		)
	}
	return transfers
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		AccessTokenSymmetricKey: util.RandomString(32),
		AccessTokenDuration:     time.Minute,
	}

	server, err := NewServer(store, config)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
