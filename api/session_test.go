package api

import (
	"bytes"
	"encoding/json"
	mockdb "github.com/CrunchyBlue/Golang-Bank/db/mock"
	db "github.com/CrunchyBlue/Golang-Bank/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRenewAccessTokenAPI(t *testing.T) {
	user, _ := generateMockUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore, refreshToken string)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore, refreshToken string) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).Return(
					db.Session{
						ID:           uuid.UUID{},
						Username:     user.Username,
						RefreshToken: refreshToken,
						UserAgent:    "",
						ClientIp:     "",
						IsBlocked:    false,
						ExpiresAt:    time.Now().Add(time.Minute),
						CreatedAt:    time.Now(),
					}, nil,
				)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BlockedSession",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore, refreshToken string) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).Return(
					db.Session{
						ID:           uuid.UUID{},
						Username:     user.Username,
						RefreshToken: refreshToken,
						UserAgent:    "",
						ClientIp:     "",
						IsBlocked:    true,
						ExpiresAt:    time.Now().Add(-time.Minute),
						CreatedAt:    time.Now(),
					}, nil,
				)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidSessionUser",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore, refreshToken string) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).Return(
					db.Session{
						ID:           uuid.UUID{},
						Username:     "InvalidSessionUser",
						RefreshToken: refreshToken,
						UserAgent:    "",
						ClientIp:     "",
						IsBlocked:    true,
						ExpiresAt:    time.Now().Add(-time.Minute),
						CreatedAt:    time.Now(),
					}, nil,
				)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "TokenMismatch",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore, refreshToken string) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).Return(
					db.Session{
						ID:           uuid.UUID{},
						Username:     user.Username,
						RefreshToken: "TokenMismatch",
						UserAgent:    "",
						ClientIp:     "",
						IsBlocked:    false,
						ExpiresAt:    time.Now().Add(time.Minute),
						CreatedAt:    time.Now(),
					}, nil,
				)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredSession",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore, refreshToken string) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).Return(
					db.Session{
						ID:           uuid.UUID{},
						Username:     user.Username,
						RefreshToken: refreshToken,
						UserAgent:    "",
						ClientIp:     "",
						IsBlocked:    false,
						ExpiresAt:    time.Now().Add(-time.Minute),
						CreatedAt:    time.Now(),
					}, nil,
				)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

				server := newTestServer(t, store)
				recorder := httptest.NewRecorder()

				refreshToken, _, err := server.tokenMaker.CreateToken(
					user.Username, server.config.RefreshTokenDuration,
				)

				tc.buildStubs(store, refreshToken)

				refreshPayload := gin.H{
					"refresh_token": refreshToken,
				}

				// Marshal body data to JSON
				data, err := json.Marshal(refreshPayload)
				require.NoError(t, err)

				url := "/session/renew"
				request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(recorder)
			},
		)
	}
}
