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

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func createRandomAccount() (Account, CreateAccountParams, error) {
	params := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)

	return account, params, err
}

func createRandomEntry(accountId int64) (Entry, CreateEntryParams, error) {
	params := CreateEntryParams{
		AccountID: accountId,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), params)

	return entry, params, err
}

func createRandomTransfer(sourceAccountId int64, destinationAccountId int64) (Transfer, CreateTransferParams, error) {
	params := CreateTransferParams{
		DestinationAccountID: destinationAccountId,
		SourceAccountID:      sourceAccountId,
		Amount:               util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), params)

	return transfer, params, err
}

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
