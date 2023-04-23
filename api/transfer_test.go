package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/CrunchyBlue/Golang-Bank/db/mock"
	db "github.com/CrunchyBlue/Golang-Bank/sqlc"
	"github.com/CrunchyBlue/Golang-Bank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTransferAPI(t *testing.T) {
	sourceAccountID := util.RandomInt(1, 1000)
	destinationAccountID := util.RandomInt(1, 1000)
	transfer := generateMockTransfers(1, sourceAccountID, destinationAccountID)[0]

	testCases := []struct {
		name          string
		transferID    int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			transferID: transfer.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTransfer(gomock.Any(), gomock.Eq(transfer.ID)).Times(1).Return(transfer, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesTransfer(t, recorder.Body, transfer)
			},
		},
		{
			name:       "NotFound",
			transferID: transfer.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTransfer(gomock.Any(), gomock.Eq(transfer.ID)).Times(1).Return(
					db.Transfer{}, sql.ErrNoRows,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			transferID: transfer.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTransfer(gomock.Any(), gomock.Eq(transfer.ID)).Times(1).Return(
					db.Transfer{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "BadRequest",
			transferID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTransfer(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store := mockdb.NewMockStore(ctrl)
				tc.buildStubs(store)

				server := NewServer(store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/transfer/%d", tc.transferID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestGetTransfersAPI(t *testing.T) {
	sourceAccountID := util.RandomInt(1, 1000)
	destinationAccountID := util.RandomInt(1, 1000)
	transfers := generateMockTransfers(10, sourceAccountID, destinationAccountID)

	testCases := []struct {
		name          string
		pageSize      int
		pageNumber    int
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			pageSize:   10,
			pageNumber: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTransfers(gomock.Any(), gomock.Any()).Times(1).Return(transfers, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesTransfers(t, recorder.Body, transfers)
			},
		},
		{
			name:       "InternalError",
			pageSize:   10,
			pageNumber: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTransfers(gomock.Any(), gomock.Any()).Times(1).Return(
					[]db.Transfer{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "BadRequest",
			pageSize:   -1,
			pageNumber: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTransfers(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store := mockdb.NewMockStore(ctrl)
				tc.buildStubs(store)

				server := NewServer(store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/transfers?page_size=%d&page_number=%d", tc.pageSize, tc.pageNumber)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestGetOutboundTransfersForAccountAPI(t *testing.T) {
	sourceAccountID := util.RandomInt(1, 1000)
	destinationAccountID := util.RandomInt(1, 1000)
	transfers := generateMockTransfers(10, sourceAccountID, destinationAccountID)

	testCases := []struct {
		name          string
		accountID     int64
		pageSize      int
		pageNumber    int
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			accountID:  sourceAccountID,
			pageSize:   10,
			pageNumber: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetOutboundTransfersForAccount(gomock.Any(), gomock.Any()).Times(1).Return(
					transfers, nil,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesTransfers(t, recorder.Body, transfers)
			},
		},
		{
			name:       "InternalError",
			accountID:  sourceAccountID,
			pageSize:   10,
			pageNumber: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetOutboundTransfersForAccount(gomock.Any(), gomock.Any()).Times(1).Return(
					[]db.Transfer{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "BadRequest",
			accountID:  sourceAccountID,
			pageSize:   -1,
			pageNumber: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetOutboundTransfersForAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store := mockdb.NewMockStore(ctrl)
				tc.buildStubs(store)

				server := NewServer(store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf(
					"/transfers/%d/outbound?page_size=%d&page_number=%d", tc.accountID, tc.pageSize, tc.pageNumber,
				)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestGetInboundTransfersForAccountAPI(t *testing.T) {
	sourceAccountID := util.RandomInt(1, 1000)
	destinationAccountID := util.RandomInt(1, 1000)
	transfers := generateMockTransfers(10, sourceAccountID, destinationAccountID)

	testCases := []struct {
		name          string
		accountID     int64
		pageSize      int
		pageNumber    int
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			accountID:  destinationAccountID,
			pageSize:   10,
			pageNumber: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetInboundTransfersForAccount(gomock.Any(), gomock.Any()).Times(1).Return(transfers, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesTransfers(t, recorder.Body, transfers)
			},
		},
		{
			name:       "InternalError",
			accountID:  destinationAccountID,
			pageSize:   10,
			pageNumber: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetInboundTransfersForAccount(gomock.Any(), gomock.Any()).Times(1).Return(
					[]db.Transfer{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "BadRequest",
			accountID:  destinationAccountID,
			pageSize:   -1,
			pageNumber: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetInboundTransfersForAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store := mockdb.NewMockStore(ctrl)
				tc.buildStubs(store)

				server := NewServer(store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf(
					"/transfers/%d/inbound?page_size=%d&page_number=%d", tc.accountID, tc.pageSize, tc.pageNumber,
				)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestCreateTransferAPI(t *testing.T) {
	accounts := generateMockAccounts(2)

	accounts[0].Currency = util.USD
	accounts[1].Currency = util.USD

	amount := util.RandomInt(1, 1000)
	currency := util.USD

	testCases := []struct {
		name                 string
		sourceAccountID      int64
		destinationAccountID int64
		amount               int64
		currency             string
		buildStubs           func(store *mockdb.MockStore)
		checkResponse        func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:                 "OK",
			sourceAccountID:      accounts[0].ID,
			destinationAccountID: accounts[1].ID,
			amount:               amount,
			currency:             currency,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accounts[0].ID)).Times(1).Return(accounts[0], nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accounts[1].ID)).Times(1).Return(accounts[1], nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:                 "InternalError",
			sourceAccountID:      accounts[0].ID,
			destinationAccountID: accounts[1].ID,
			amount:               amount,
			currency:             currency,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accounts[0].ID)).Times(1).Return(accounts[0], nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accounts[1].ID)).Times(1).Return(accounts[1], nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).Return(
					db.TransferTxResult{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateTransfer(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store := mockdb.NewMockStore(ctrl)
				tc.buildStubs(store)

				server := NewServer(store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprint("/transfer")

				jsonEntry := fmt.Sprintf(
					`{"source_account_id": %d, "destination_account_id": %d, "amount": %d, "currency": "%s"}`,
					tc.sourceAccountID,
					tc.destinationAccountID, tc.amount, tc.currency,
				)
				jsonBody := []byte(jsonEntry)
				bodyReader := bytes.NewReader(jsonBody)

				request, err := http.NewRequest(http.MethodPost, url, bodyReader)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestUpdateTransferAPI(t *testing.T) {
	sourceAccountID := util.RandomInt(1, 1000)
	destinationAccountID := util.RandomInt(1, 1000)
	transfer := generateMockTransfers(1, sourceAccountID, destinationAccountID)[0]

	amount := util.RandomInt(1, 1000)

	expectedTransfer := db.Transfer{
		ID:                   transfer.ID,
		SourceAccountID:      transfer.SourceAccountID,
		DestinationAccountID: transfer.DestinationAccountID,
		Amount:               amount,
		CreatedAt:            transfer.CreatedAt,
	}

	testCases := []struct {
		name          string
		transferID    int64
		amount        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			transferID: expectedTransfer.ID,
			amount:     amount,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateTransfer(gomock.Any(), gomock.Any()).Times(1).Return(
					expectedTransfer, nil,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesUpdatedTransfer(
					t, recorder.Body, transfer, amount,
				)
			},
		},
		{
			name:       "NotFound",
			transferID: expectedTransfer.ID,
			amount:     amount,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateTransfer(gomock.Any(), gomock.Any()).Times(1).Return(
					db.Transfer{}, sql.ErrNoRows,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			transferID: expectedTransfer.ID,
			amount:     amount,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateTransfer(gomock.Any(), gomock.Any()).Times(1).Return(
					db.Transfer{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "BadRequest",
			transferID: 0,
			amount:     amount,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateTransfer(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store := mockdb.NewMockStore(ctrl)
				tc.buildStubs(store)

				server := NewServer(store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/transfer/%d", tc.transferID)

				jsonAccount := fmt.Sprintf(
					`{"amount": %d}`, tc.amount,
				)
				jsonBody := []byte(jsonAccount)
				bodyReader := bytes.NewReader(jsonBody)

				request, err := http.NewRequest(http.MethodPut, url, bodyReader)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestDeleteTransferAPI(t *testing.T) {
	sourceAccountID := util.RandomInt(1, 1000)
	destinationAccountID := util.RandomInt(1, 1000)
	transfer := generateMockTransfers(1, sourceAccountID, destinationAccountID)[0]

	testCases := []struct {
		name          string
		transferID    int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			transferID: transfer.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteTransfer(gomock.Any(), gomock.Eq(transfer.ID)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "NotFound",
			transferID: transfer.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteTransfer(gomock.Any(), gomock.Eq(transfer.ID)).Times(1).Return(
					sql.ErrNoRows,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:       "InternalError",
			transferID: transfer.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteTransfer(gomock.Any(), gomock.Eq(transfer.ID)).Times(1).Return(
					sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "BadRequest",
			transferID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteTransfer(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(
			tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store := mockdb.NewMockStore(ctrl)
				tc.buildStubs(store)

				server := NewServer(store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/transfer/%d", tc.transferID)
				request, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func requireBodyMatchesTransfer(t *testing.T, body *bytes.Buffer, transfer db.Transfer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var fetchedTransfer db.Transfer
	err = json.Unmarshal(data, &fetchedTransfer)
	require.NoError(t, err)
	require.Equal(t, transfer, fetchedTransfer)
}

func requireBodyMatchesTransfers(t *testing.T, body *bytes.Buffer, transfers []db.Transfer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var fetchedTransfers []db.Transfer
	err = json.Unmarshal(data, &fetchedTransfers)
	require.NoError(t, err)

	for i, _ := range transfers {
		require.Equal(t, transfers[i], fetchedTransfers[i])
	}
}

func requireBodyMatchesUpdatedTransfer(
	t *testing.T, body *bytes.Buffer, transfer db.Transfer, amount int64,
) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var updatedTransfer db.Transfer
	err = json.Unmarshal(data, &updatedTransfer)
	require.NoError(t, err)

	require.Equal(t, transfer.ID, updatedTransfer.ID)
	require.Equal(t, transfer.SourceAccountID, updatedTransfer.SourceAccountID)
	require.Equal(t, transfer.DestinationAccountID, updatedTransfer.DestinationAccountID)
	require.Equal(t, transfer.CreatedAt, updatedTransfer.CreatedAt)

	require.Equal(t, amount, updatedTransfer.Amount)
}
