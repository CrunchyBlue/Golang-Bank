package api

import (
	db "github.com/CrunchyBlue/Golang-Bank/sqlc"
	"github.com/CrunchyBlue/Golang-Bank/util"
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

func generateMockAccounts(numAccounts int) []db.Account {
	var accounts []db.Account

	for i := 0; i < numAccounts; i++ {
		accounts = append(
			accounts, db.Account{
				ID:       util.RandomInt(1, 1000),
				Owner:    util.RandomOwner(),
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

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
