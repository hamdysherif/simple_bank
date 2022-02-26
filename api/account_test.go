package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/hamdysherif/simplebank/db/mock"
	db "github.com/hamdysherif/simplebank/db/sqlc"
	"github.com/hamdysherif/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchResponse(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Badrequest",
			accountID: int64(0),
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)
	server := NewServer(store)

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			tc.buildStubs(store)
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()
	args := createAccountRequest{
		Owner:    account.Owner,
		Currency: account.Currency,
	}

	testCases := []struct {
		name          string
		params        createAccountRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			params: args,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), db.CreateAccountParams{Owner: args.Owner, Currency: args.Currency, Balance: 0}).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "BadRequestCurrencyNotProvided",
			params: createAccountRequest{Owner: "Ali", Currency: ""},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "BadRequestCurrencyNotAllowed",
			params: createAccountRequest{Owner: "Ali", Currency: "ANY"},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "BadRequestOwner",
			params: createAccountRequest{Owner: "", Currency: "LE"},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			params: createAccountRequest{Owner: "Ali", Currency: "LE"},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), db.CreateAccountParams{Owner: "Ali", Currency: "LE", Balance: 0}).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	server := NewServer(store)
	url := "/accounts"

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			var reqBody bytes.Buffer
			err := json.NewEncoder(&reqBody).Encode(tc.params)
			require.NoError(t, err)
			tc.buildStubs(store)

			request, err := http.NewRequest(http.MethodPost, url, &reqBody)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 10),
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
		Balance:  util.RandomeBalance(),
	}
}

func requireBodyMatchResponse(t *testing.T, body *bytes.Buffer, account db.Account) {
	b, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var getAccount db.Account
	err = json.Unmarshal(b, &getAccount)
	require.NoError(t, err)
	require.Equal(t, getAccount, account)
}
