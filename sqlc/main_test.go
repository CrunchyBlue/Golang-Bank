package db

import (
	"context"
	"database/sql"
	"github.com/CrunchyBlue/Golang-Bank/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func createRandomAccount() (Account, CreateAccountParams, error) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomInt(0, 1000),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	return account, arg, err
}

func createRandomEntry(accountId int64) (Entry, CreateEntryParams, error) {
	arg := CreateEntryParams{
		AccountID: accountId,
		Amount:    util.RandomInt(0, 1000),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	return entry, arg, err
}

func createRandomTransfer(sourceAccountId int64, destinationAccountId int64) (Transfer, CreateTransferParams, error) {
	arg := CreateTransferParams{
		DestinationAccountID: destinationAccountId,
		SourceAccountID:      sourceAccountId,
		Amount:               util.RandomInt(0, 1000),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	return transfer, arg, err
}

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
