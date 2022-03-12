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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/hamdysherif/simplebank/db/mock"
	db "github.com/hamdysherif/simplebank/db/sqlc"
	"github.com/hamdysherif/simplebank/util"
	"github.com/stretchr/testify/assert"
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
				assert.Equal(t, http.StatusOK, recorder.Code)
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
				assert.Equal(t, http.StatusNotFound, recorder.Code)
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
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
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
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)
	server := NewTestServer(t, store)

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			tc.buildStubs(store)

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			authorizationType := "bearer"
			token, err := server.tokenMaker.CreateToken("hamdy", time.Hour)
			assert.NoError(t, err)

			request.Header.Set("authorization", fmt.Sprintf("%s %s", authorizationType, token))

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		params        gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			params: gin.H{"owner": account.Owner, "currency": account.Currency},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), db.CreateAccountParams{
						Owner:    account.Owner,
						Currency: account.Currency,
						Balance:  0}).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchResponse(t, recorder.Body, account)
			},
		},
		{
			name:   "BadRequestCurrencyNotProvided",
			params: gin.H{"owner": "Ali", "currency": ""},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "BadRequestCurrencyNotAllowed",
			params: gin.H{"owner": "Ali", "currency": "ANY"},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "BadRequestOwner",
			params: gin.H{"owner": "", "currency": "LE"},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			params: gin.H{"owner": "Ali", "currency": "LE"},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), db.CreateAccountParams{Owner: "Ali", Currency: "LE", Balance: 0}).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	server := NewTestServer(t, store)
	url := "/accounts"

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			body, err := json.Marshal(tc.params)
			assert.NoError(t, err)
			tc.buildStubs(store)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			assert.NoError(t, err)

			authorizationType := "bearer"
			token, err := server.tokenMaker.CreateToken("hamdy", time.Hour)
			assert.NoError(t, err)

			request.Header.Set("authorization", fmt.Sprintf("%s %s", authorizationType, token))

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
	assert.NoError(t, err)

	var getAccount db.Account
	err = json.Unmarshal(b, &getAccount)
	assert.NoError(t, err)
	assert.Equal(t, getAccount, account)
}
