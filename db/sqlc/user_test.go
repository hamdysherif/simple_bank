package db

import (
	"context"
	"testing"

	"github.com/hamdysherif/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.GenerateHashedPassowrd("secret")
	require.NoError(t, err)
	args := CreateUserParams{
		FullName:       util.RandomOwner(),
		Username:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.Email, user.Email)
	require.NotZero(t, user.ID)

	return user

}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user, err := testQueries.GetUser(context.Background(), user1.ID)

	require.NoError(t, err)
	require.Equal(t, user.ID, user1.ID)
	require.Equal(t, user1.FullName, user.FullName)
	require.Equal(t, user1.Username, user.Username)
	require.Equal(t, user1.Email, user.Email)
}
