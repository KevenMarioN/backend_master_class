package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kevenmarion/backend_master_class/db/mocks"
	db "github.com/kevenmarion/backend_master_class/db/sqlc"
	"github.com/kevenmarion/backend_master_class/util"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mocks.Store)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, mock.Anything).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, []db.Account{account})
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, mock.Anything).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mocks.Store) {
				store.On("GetAccount", mock.Anything, mock.Anything).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "InvalidID",
			accountID:  -1,
			buildStubs: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeMock := mocks.NewStore(t)
			tc.buildStubs(storeMock)
			server := NewServer(storeMock)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.LoadRouters()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListAccountApi(t *testing.T) {
	accounts := []db.Account{
		randomAccount(),
		randomAccount(),
		randomAccount(),
	}

	testCases := []struct {
		name   string
		filter struct {
			PageSize int32
			PageID   int32
		}
		buildStubs    func(store *mocks.Store)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			filter: struct {
				PageSize int32
				PageID   int32
			}{
				PageSize: 5,
				PageID:   1,
			},
			buildStubs: func(store *mocks.Store) {
				store.On("ListAccounts", mock.Anything, mock.Anything).Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, accounts)
			},
		},
		{
			name: "InvalidParams",
			filter: struct {
				PageSize int32
				PageID   int32
			}{
				PageSize: 20,
				PageID:   -1,
			},
			buildStubs: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			filter: struct {
				PageSize int32
				PageID   int32
			}{
				PageSize: 5,
				PageID:   1,
			},
			buildStubs: func(store *mocks.Store) {
				store.On("ListAccounts", mock.Anything, mock.Anything).Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeMock := mocks.NewStore(t)
			tc.buildStubs(storeMock)
			server := NewServer(storeMock)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts?offset=%v&limit=%v", tc.filter.PageID, tc.filter.PageSize)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.LoadRouters()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateAccountApi(t *testing.T) {
	account := randomAccount()
	accountJson, err := json.Marshal(db.CreateAccountParams{
		Owner:    account.Owner,
		Currency: account.Currency,
	})
	require.NoError(t, err)

	testCases := []struct {
		name          string
		json          []byte
		buildStubs    func(store *mocks.Store)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			json: accountJson,
			buildStubs: func(store *mocks.Store) {
				store.On("CreateAccount", mock.Anything, db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
				}).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, []db.Account{account})
			},
		},
		{
			name:       "InvalidParams",
			json:       []byte{},
			buildStubs: func(store *mocks.Store) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			json: accountJson,
			buildStubs: func(store *mocks.Store) {
				store.On("CreateAccount", mock.Anything, db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
				}).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeMock := mocks.NewStore(t)
			tc.buildStubs(storeMock)
			server := NewServer(storeMock)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(tc.json))
			require.NoError(t, err)

			server.LoadRouters()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
func randomAccount() db.Account {
	return db.Account{
		ID:        util.RandomInt(1, 1000),
		Owner:     util.RandomOwner(),
		Balance:   util.RandomMoney(),
		Currency:  util.RandomCurrency(),
		CreatedAt: time.Now(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	var (
		err      error
		response []byte
	)
	if len(accounts) > 1 {
		response, err = json.Marshal(accounts)
		require.NoError(t, err)
	} else {
		response, err = json.Marshal(accounts[0])
		require.NoError(t, err)

	}

	result, err := io.ReadAll(body)
	require.NoError(t, err)
	require.Equal(t, string(result), string(response))
}
