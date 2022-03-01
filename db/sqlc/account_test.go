package db

import (
	"context"
	"testing"

	"github.com/hamdysherif/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {

	user := createRandomUser(t)
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomeBalance(),
		Currency: util.RandomCurrency(),
		UserID:   user.ID,
	}

	account, err := testQueries.CreateAccount(context.Background(), args)

	require.Nil(t, err)
	require.Equal(t, args.Balance, account.Balance)
	require.NotZero(t, account.ID)

	return account

}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	account, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Nil(t, err)
	require.Equal(t, account.ID, account1.ID)
	require.Equal(t, account.Balance, account1.Balance)
	require.Equal(t, account.Currency, account1.Currency)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomeBalance(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, account2.Balance, args.Balance)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	args := ListAccountsParams{
		Offset: 0,
		Limit:  5,
	}
	lstAccounts, err := testQueries.ListAccounts(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, lstAccounts, 5)

	for _, account := range lstAccounts {
		require.NotEmpty(t, account)
	}
}

func TestAddAccountBalance(t *testing.T) {
	account1 := createRandomAccount(t)

	args := AddAccountBalanceParams{
		ID:     account1.ID,
		Amount: util.RandomeBalance(),
	}

	account2, err := testQueries.AddAccountBalance(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, account2.Balance, args.Amount+account1.Balance)
}
