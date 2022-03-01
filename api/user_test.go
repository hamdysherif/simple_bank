package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/hamdysherif/simplebank/db/mock"
	db "github.com/hamdysherif/simplebank/db/sqlc"
	"github.com/hamdysherif/simplebank/util"
	"github.com/stretchr/testify/require"
)

type eqUserParamMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqUserParamMatcher) Matches(x interface{}) bool {
	arg := x.(db.CreateUserParams)

	matched := util.CheckHashedPassword(arg.HashedPassword, e.password)
	if !matched {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(arg, e.arg)
}

func (e eqUserParamMatcher) String() string {
	return fmt.Sprintf("is equal to %v with password: %v", e.arg, e.password)
}

// EqUserParam returns a matcher that matches on equality for user params.
func EqUserParam(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqUserParamMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	password := "secred"
	user := randomUser(password)

	testCases := []struct {
		name          string
		params        gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{

			name:   "OK",
			params: gin.H{"username": user.Username, "email": user.Email, "full_name": user.FullName, "password": password},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Email:          user.Email,
					Username:       user.Username,
					FullName:       user.FullName,
					HashedPassword: user.HashedPassword,
				}
				store.
					EXPECT().
					CreateUser(gomock.Any(), EqUserParam(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchResponseUser(t, recorder.Body, user)
			},
		},
		{
			name:   "StatusBadRequest_Username",
			params: gin.H{"username": "", "email": user.Email, "full_name": user.FullName, "password": password},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "StatusInternalServerError",
			params: gin.H{"username": user.Username, "email": user.Email, "full_name": user.FullName, "password": password},
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
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
	url := "/users"

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			body, err := json.Marshal(tc.params)
			require.NoError(t, err)
			tc.buildStubs(store)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
func randomUser(password string) db.User {
	hashedPassowrd, _ := util.GenerateHashedPassowrd(password)
	return db.User{
		ID:             util.RandomInt(1, 10),
		Email:          util.RandomEmail(),
		FullName:       util.RandomOwner(),
		HashedPassword: hashedPassowrd,
		Username:       util.RandomOwner(),
	}
}

func requireBodyMatchResponseUser(t *testing.T, body *bytes.Buffer, user db.User) {
	b, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var getUser db.User
	err = json.Unmarshal(b, &getUser)
	require.NoError(t, err)
	require.Equal(t, getUser, user)
}
