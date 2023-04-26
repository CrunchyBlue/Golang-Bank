package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/CrunchyBlue/Golang-Bank/db/mock"
	db "github.com/CrunchyBlue/Golang-Bank/sqlc"
	"github.com/CrunchyBlue/Golang-Bank/token"
	"github.com/CrunchyBlue/Golang-Bank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := generateMockUser(t)
	account := generateMockAccounts(user.Username, 1)[0]

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(
					db.Account{}, sql.ErrNoRows,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(
					db.Account{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
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

				server := newTestServer(t, store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/account/%d", tc.accountID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)

				tc.setupAuth(t, request, server.tokenMaker)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestGetAccountsAPI(t *testing.T) {
	user, _ := generateMockUser(t)
	accounts := generateMockAccounts(user.Username, 10)

	testCases := []struct {
		name          string
		pageSize      int
		pageNumber    int
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			pageSize:   10,
			pageNumber: 1,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Times(1).Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesAccounts(t, recorder.Body, accounts)
			},
		},
		{
			name:       "InternalError",
			pageSize:   10,
			pageNumber: 1,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Times(1).Return(
					[]db.Account{}, sql.ErrConnDone,
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
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Times(0)
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

				server := newTestServer(t, store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/accounts?page_size=%d&page_number=%d", tc.pageSize, tc.pageNumber)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)

				tc.setupAuth(t, request, server.tokenMaker)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestCreateAccountAPI(t *testing.T) {
	user, _ := generateMockUser(t)
	owner := user.Username
	currency := util.RandomCurrency()

	account := db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Currency: currency,
		Balance:  0,
	}

	testCases := []struct {
		name          string
		owner         string
		currency      string
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			owner:    owner,
			currency: currency,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesCreatedAccount(t, recorder.Body, account)
			},
		},
		{
			name:     "InternalError",
			owner:    owner,
			currency: currency,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(
					db.Account{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
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

				server := newTestServer(t, store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprint("/account")

				jsonAccount := fmt.Sprintf(`{"owner": "%s", "currency": "%s"}`, tc.owner, tc.currency)
				jsonBody := []byte(jsonAccount)
				bodyReader := bytes.NewReader(jsonBody)

				request, err := http.NewRequest(http.MethodPost, url, bodyReader)
				require.NoError(t, err)

				tc.setupAuth(t, request, server.tokenMaker)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestUpdateAccountAPI(t *testing.T) {
	user, _ := generateMockUser(t)
	account := generateMockAccounts(user.Username, 1)[0]

	owner := util.RandomOwner()
	currency := util.RandomCurrency()
	balance := util.RandomInt(-1000, 1000)

	expectedAccount := db.Account{
		ID:        account.ID,
		Owner:     owner,
		Currency:  currency,
		Balance:   balance,
		CreatedAt: account.CreatedAt,
	}

	testCases := []struct {
		name          string
		accountID     int64
		owner         string
		currency      string
		balance       int64
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			owner:     owner,
			currency:  currency,
			balance:   balance,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Times(1).Return(
					expectedAccount, nil,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesUpdatedAccount(t, recorder.Body, account, owner, currency, balance)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			owner:     owner,
			currency:  currency,
			balance:   balance,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Times(1).Return(
					db.Account{}, sql.ErrNoRows,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			owner:     owner,
			currency:  currency,
			balance:   balance,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Times(1).Return(
					db.Account{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			owner:     owner,
			currency:  currency,
			balance:   balance,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Times(0)
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

				server := newTestServer(t, store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/account/%d", tc.accountID)

				jsonAccount := fmt.Sprintf(
					`{"owner": "%s", "currency": "%s", "balance": %d}`, tc.owner, tc.currency, tc.balance,
				)
				jsonBody := []byte(jsonAccount)
				bodyReader := bytes.NewReader(jsonBody)

				request, err := http.NewRequest(http.MethodPut, url, bodyReader)
				require.NoError(t, err)

				tc.setupAuth(t, request, server.tokenMaker)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestUpdateAccountBalanceAPI(t *testing.T) {
	user, _ := generateMockUser(t)
	account := generateMockAccounts(user.Username, 1)[0]

	amount := util.RandomInt(-1000, 1000)

	expectedAccount := db.Account{
		ID:        account.ID,
		Owner:     account.Owner,
		Currency:  account.Currency,
		Balance:   account.Balance + amount,
		CreatedAt: account.CreatedAt,
	}

	testCases := []struct {
		name          string
		accountID     int64
		owner         string
		currency      string
		amount        int64
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			amount:    amount,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAccountBalance(gomock.Any(), gomock.Any()).Times(1).Return(
					expectedAccount, nil,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesUpdatedAccount(
					t, recorder.Body, account, account.Owner, account.Currency, account.Balance+amount,
				)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			amount:    amount,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAccountBalance(gomock.Any(), gomock.Any()).Times(1).Return(
					db.Account{}, sql.ErrNoRows,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			amount:    amount,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAccountBalance(gomock.Any(), gomock.Any()).Times(1).Return(
					db.Account{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			amount:    amount,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAccountBalance(gomock.Any(), gomock.Any()).Times(0)
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

				server := newTestServer(t, store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/account/%d/balance", tc.accountID)

				jsonAccount := fmt.Sprintf(
					`{"amount": %d}`, tc.amount,
				)
				jsonBody := []byte(jsonAccount)
				bodyReader := bytes.NewReader(jsonBody)

				request, err := http.NewRequest(http.MethodPut, url, bodyReader)
				require.NoError(t, err)

				tc.setupAuth(t, request, server.tokenMaker)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestDeleteAccountAPI(t *testing.T) {
	user, _ := generateMockUser(t)
	account := generateMockAccounts(user.Username, 1)[0]

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(
					sql.ErrNoRows,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(
					sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Any()).Times(0)
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

				server := newTestServer(t, store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/account/%d", tc.accountID)
				request, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)

				tc.setupAuth(t, request, server.tokenMaker)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func requireBodyMatchesAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var fetchedAccount db.Account
	err = json.Unmarshal(data, &fetchedAccount)
	require.NoError(t, err)
	require.Equal(t, account, fetchedAccount)
}

func requireBodyMatchesAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var fetchedAccounts []db.Account
	err = json.Unmarshal(data, &fetchedAccounts)
	require.NoError(t, err)

	for i, _ := range accounts {
		require.Equal(t, accounts[i], fetchedAccounts[i])
	}
}

func requireBodyMatchesCreatedAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var createdAccount db.Account
	err = json.Unmarshal(data, &createdAccount)
	require.NoError(t, err)

	require.Equal(t, account.Owner, createdAccount.Owner)
	require.Equal(t, account.Currency, createdAccount.Currency)
	require.Equal(t, account.Balance, createdAccount.Balance)
}

func requireBodyMatchesUpdatedAccount(
	t *testing.T, body *bytes.Buffer, account db.Account, owner string, currency string, balance int64,
) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var updatedAccount db.Account
	err = json.Unmarshal(data, &updatedAccount)
	require.NoError(t, err)

	require.Equal(t, account.ID, updatedAccount.ID)
	require.Equal(t, account.CreatedAt, updatedAccount.CreatedAt)

	require.Equal(t, owner, updatedAccount.Owner)
	require.Equal(t, currency, updatedAccount.Currency)
	require.Equal(t, balance, updatedAccount.Balance)
}
