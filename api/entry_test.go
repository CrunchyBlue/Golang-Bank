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

func TestGetEntryAPI(t *testing.T) {
	accountID := util.RandomInt(1, 1000)
	entry := generateMockEntries(1, accountID)[0]

	testCases := []struct {
		name          string
		entryID       int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			entryID: entry.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetEntry(gomock.Any(), gomock.Eq(entry.ID)).Times(1).Return(entry, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesEntry(t, recorder.Body, entry)
			},
		},
		{
			name:    "NotFound",
			entryID: entry.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetEntry(gomock.Any(), gomock.Eq(entry.ID)).Times(1).Return(
					db.Entry{}, sql.ErrNoRows,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			entryID: entry.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetEntry(gomock.Any(), gomock.Eq(entry.ID)).Times(1).Return(
					db.Entry{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "BadRequest",
			entryID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetEntry(gomock.Any(), gomock.Any()).Times(0)
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

				url := fmt.Sprintf("/entry/%d", tc.entryID)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestGetEntriesAPI(t *testing.T) {
	accountID := util.RandomInt(1, 1000)
	entries := generateMockEntries(10, accountID)

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
				store.EXPECT().GetEntries(gomock.Any(), gomock.Any()).Times(1).Return(entries, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesEntries(t, recorder.Body, entries)
			},
		},
		{
			name:       "InternalError",
			pageSize:   10,
			pageNumber: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetEntries(gomock.Any(), gomock.Any()).Times(1).Return(
					[]db.Entry{}, sql.ErrConnDone,
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
				store.EXPECT().GetEntries(gomock.Any(), gomock.Any()).Times(0)
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

				url := fmt.Sprintf("/entries?page_size=%d&page_number=%d", tc.pageSize, tc.pageNumber)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestGetEntriesForAccountAPI(t *testing.T) {
	accountID := util.RandomInt(1, 1000)
	entries := generateMockEntries(10, accountID)

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
			accountID:  accountID,
			pageSize:   10,
			pageNumber: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetEntriesForAccount(gomock.Any(), gomock.Any()).Times(1).Return(entries, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesEntries(t, recorder.Body, entries)
			},
		},
		{
			name:       "InternalError",
			accountID:  accountID,
			pageSize:   10,
			pageNumber: 1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetEntriesForAccount(gomock.Any(), gomock.Any()).Times(1).Return(
					[]db.Entry{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "BadRequest",
			accountID:  accountID,
			pageSize:   -1,
			pageNumber: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetEntriesForAccount(gomock.Any(), gomock.Any()).Times(0)
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

				url := fmt.Sprintf("/entries/%d?page_size=%d&page_number=%d", tc.accountID, tc.pageSize, tc.pageNumber)
				request, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func TestCreateEntryAPI(t *testing.T) {
	accountID := util.RandomInt(1, 1000)
	amount := util.RandomInt(-1000, 1000)

	entry := db.Entry{
		ID:        util.RandomInt(1, 1000),
		AccountID: accountID,
		Amount:    amount,
	}

	testCases := []struct {
		name          string
		accountID     int64
		amount        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: accountID,
			amount:    amount,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(1).Return(entry, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesCreatedEntry(t, recorder.Body, entry)
			},
		},
		{
			name:      "InternalError",
			accountID: accountID,
			amount:    amount,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(1).Return(
					db.Entry{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
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

				url := fmt.Sprint("/entry")

				jsonEntry := fmt.Sprintf(`{"account_id": %d, "amount": %d}`, tc.accountID, tc.amount)
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

func TestUpdateEntryAPI(t *testing.T) {
	accountID := util.RandomInt(1, 1000)
	entry := generateMockEntries(1, accountID)[0]

	amount := util.RandomInt(-1000, 1000)

	expectedEntry := db.Entry{
		ID:        entry.ID,
		AccountID: entry.AccountID,
		Amount:    amount,
		CreatedAt: entry.CreatedAt,
	}

	testCases := []struct {
		name          string
		entryID       int64
		amount        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			entryID: expectedEntry.ID,
			amount:  amount,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateEntry(gomock.Any(), gomock.Any()).Times(1).Return(
					expectedEntry, nil,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchesUpdatedEntry(
					t, recorder.Body, entry, amount,
				)
			},
		},
		{
			name:    "NotFound",
			entryID: -1,
			amount:  amount,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateEntry(gomock.Any(), gomock.Any()).Times(1).Return(
					db.Entry{}, sql.ErrNoRows,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			entryID: expectedEntry.ID,
			amount:  amount,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateEntry(gomock.Any(), gomock.Any()).Times(1).Return(
					db.Entry{}, sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "BadRequest",
			entryID: 0,
			amount:  amount,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateEntry(gomock.Any(), gomock.Any()).Times(0)
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

				url := fmt.Sprintf("/entry/%d", tc.entryID)

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

func TestDeleteEntryAPI(t *testing.T) {
	accountID := util.RandomInt(1, 1000)
	entry := generateMockEntries(1, accountID)[0]

	testCases := []struct {
		name          string
		entryID       int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			entryID: entry.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteEntry(gomock.Any(), gomock.Eq(entry.ID)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:    "NotFound",
			entryID: entry.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteEntry(gomock.Any(), gomock.Eq(entry.ID)).Times(1).Return(
					sql.ErrNoRows,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			entryID: entry.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteEntry(gomock.Any(), gomock.Eq(entry.ID)).Times(1).Return(
					sql.ErrConnDone,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "BadRequest",
			entryID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteEntry(gomock.Any(), gomock.Any()).Times(0)
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

				url := fmt.Sprintf("/entry/%d", tc.entryID)
				request, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			},
		)
	}
}

func requireBodyMatchesEntry(t *testing.T, body *bytes.Buffer, entry db.Entry) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var fetchedEntry db.Entry
	err = json.Unmarshal(data, &fetchedEntry)
	require.NoError(t, err)
	require.Equal(t, entry, fetchedEntry)
}

func requireBodyMatchesEntries(t *testing.T, body *bytes.Buffer, entries []db.Entry) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var fetchedEntries []db.Entry
	err = json.Unmarshal(data, &fetchedEntries)
	require.NoError(t, err)

	for i, _ := range entries {
		require.Equal(t, entries[i], fetchedEntries[i])
	}
}

func requireBodyMatchesCreatedEntry(t *testing.T, body *bytes.Buffer, entry db.Entry) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var createdEntry db.Entry
	err = json.Unmarshal(data, &createdEntry)
	require.NoError(t, err)

	require.Equal(t, entry.AccountID, createdEntry.AccountID)
	require.Equal(t, entry.Amount, createdEntry.Amount)
}

func requireBodyMatchesUpdatedEntry(
	t *testing.T, body *bytes.Buffer, entry db.Entry, amount int64,
) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var updatedEntry db.Entry
	err = json.Unmarshal(data, &updatedEntry)
	require.NoError(t, err)

	require.Equal(t, entry.ID, updatedEntry.ID)
	require.Equal(t, entry.AccountID, updatedEntry.AccountID)
	require.Equal(t, entry.CreatedAt, updatedEntry.CreatedAt)

	require.Equal(t, amount, updatedEntry.Amount)
}
