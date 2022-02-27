package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/hamdysherif/simplebank/db/mock"
	db "github.com/hamdysherif/simplebank/db/sqlc"
	"github.com/hamdysherif/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	account1 := db.Account{
		ID:       1,
		Owner:    util.RandomOwner(),
		Currency: util.AllowedCurrencies()[0],
		Balance:  500,
	}
	account2 := db.Account{
		ID:       util.RandomInt(1, 10),
		Owner:    util.RandomOwner(),
		Currency: util.AllowedCurrencies()[0],
		Balance:  300,
	}
	entry1 := db.Entry{AccountID: account1.ID, Amount: -5}
	entry2 := db.Entry{AccountID: account2.ID, Amount: 5}
	trans := db.Transfer{Amount: 5, FromAccountID: account1.ID, ToAccountID: account2.ID}

	testCases := []struct {
		name          string
		params        gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "Ok",
			params: gin.H{"from_account_id": account1.ID, "to_account_id": account2.ID, "amount": 5, "currency": account1.Currency},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), account1.ID).
					Times(1).
					Return(account1, nil)

				store.
					EXPECT().
					GetAccount(gomock.Any(), account2.ID).
					Times(1).
					Return(account2, nil)

				store.
					EXPECT().
					TransferTx(gomock.Any(), db.TransferParams{FromAccountID: account1.ID, ToAccountID: account2.ID, Amount: 5}).
					Times(1).
					Return(db.TransferResult{FromAccount: account1, ToAccount: account2, FromEntry: entry1, ToEntry: entry2, Transfer: trans}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				body, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)
				var r db.TransferResult
				err = json.Unmarshal(body, &r)
				require.NoError(t, err)

				require.Equal(t, r.FromAccount, account1)
				require.Equal(t, r.FromEntry, entry1)
				require.Equal(t, r.ToAccount, account2)
				require.Equal(t, r.ToEntry, entry2)
				require.Equal(t, r.Transfer, trans)
			},
		},
		{
			name:   "InternalServerError",
			params: gin.H{"from_account_id": account1.ID, "to_account_id": account2.ID, "amount": 5, "currency": account1.Currency},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), account1.ID).
					Times(1).
					Return(account1, nil)

				store.
					EXPECT().
					GetAccount(gomock.Any(), account2.ID).
					Times(1).
					Return(account2, nil)

				store.
					EXPECT().
					TransferTx(gomock.Any(), db.TransferParams{FromAccountID: account1.ID, ToAccountID: account2.ID, Amount: 5}).
					Times(1).
					Return(db.TransferResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "BadRequestAmount",
			params: gin.H{"from_account_id": account1.ID, "to_account_id": account2.ID, "amount": 0, "currency": account1.Currency},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					TransferTx(gomock.Any(), db.TransferParams{FromAccountID: account1.ID, ToAccountID: account2.ID, Amount: 5}).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "BadRequestFromAccount",
			params: gin.H{"from_account_id": account1.ID, "to_account_id": account2.ID, "amount": 10, "currency": account1.Currency},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), account1.ID).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

				store.
					EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "BadRequestCurrencyNotMatch",
			params: gin.H{"from_account_id": account1.ID, "to_account_id": account2.ID, "amount": 10, "currency": account1.Currency},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), account1.ID).
					Times(1).
					Return(db.Account{ID: util.RandomInt(5, 100), Owner: util.RandomOwner(), Currency: util.AllowedCurrencies()[1], Balance: 0}, nil)

				store.
					EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "BadRequestToAccount",
			params: gin.H{"from_account_id": account1.ID, "to_account_id": account2.ID, "amount": 10, "currency": account1.Currency},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), account1.ID).
					Times(1).
					Return(account1, nil)

				store.
					EXPECT().
					GetAccount(gomock.Any(), account2.ID).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)

				store.
					EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	url := "/transfers"
	store := mockdb.NewMockStore(ctrl)
	server := NewServer(store)

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			req, err := json.Marshal(tc.params)
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(req))
			recorder := httptest.NewRecorder()

			tc.buildStubs(store)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
