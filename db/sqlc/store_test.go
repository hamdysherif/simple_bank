package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	testStore := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	fmt.Println("Balance Before:", fromAccount.Balance, toAccount.Balance)
	n := 10
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferResult)

	args := TransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        amount,
	}

	for i := 0; i < n; i++ {
		go func(i int) {
			result, err := testStore.TransferTx(context.Background(), args)
			fmt.Println("Balance IN TX:", result.FromAccount.Balance, result.ToAccount.Balance, i)
			errs <- err
			results <- result
		}(i)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotZero(t, result.FromEntry.ID)
		require.Equal(t, result.FromEntry.Amount, -amount)
		_, fErr := testStore.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, fErr)

		require.NotZero(t, result.ToEntry.ID)
		require.Equal(t, result.ToEntry.Amount, amount)
		_, tErr := testStore.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, tErr)

		require.NotZero(t, result.Transfer.ID)
		require.Equal(t, result.Transfer.FromAccountID, args.FromAccountID)
		require.Equal(t, result.Transfer.ToAccountID, args.ToAccountID)
		_, trErr := testStore.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, trErr)

		// check account balance
		fromDiff := fromAccount.Balance - result.FromAccount.Balance
		toDiff := result.ToAccount.Balance - toAccount.Balance
		require.Equal(t, fromDiff, toDiff)
	}

	fAccount, _ := testStore.GetAccount(context.Background(), fromAccount.ID)
	tAccount, _ := testStore.GetAccount(context.Background(), toAccount.ID)

	// balalnce after all operations
	require.Equal(t, fAccount.Balance, fromAccount.Balance-int64(n)*amount)
	require.Equal(t, tAccount.Balance, toAccount.Balance+int64(n)*amount)
}

func TestTransferTxDeadLock(t *testing.T) {
	testStore := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("Balance Before:", account1.Balance, account2.Balance)
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID
		if i%2 == 0 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferParams{FromAccountID: fromAccountId, ToAccountID: toAccountId, Amount: amount})
			fmt.Println("Balance IN TX:", result.FromAccount.Balance, result.ToAccount.Balance)
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	fAccount, _ := testStore.GetAccount(context.Background(), account1.ID)
	tAccount, _ := testStore.GetAccount(context.Background(), account2.ID)

	// balalnce after all operations
	require.Equal(t, fAccount.Balance, account1.Balance)
	require.Equal(t, tAccount.Balance, account2.Balance)
}

func TestTransferTxPure(t *testing.T) {
	testStore := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	fmt.Println("Balance Beforee:", fromAccount.Balance, toAccount.Balance)
	n := 10
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferResult)

	args := TransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        amount,
	}

	for i := 0; i < n; i++ {
		go func(i int) {
			result, err := testStore.TransferTxPure(context.Background(), args)
			fmt.Println("Balance IN TX:", result.FromAccount.Balance, result.ToAccount.Balance, i)
			errs <- err
			results <- result
		}(i)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotZero(t, result.FromAccount.ID)
		// check account balance
		fromDiff := fromAccount.Balance - result.FromAccount.Balance
		toDiff := result.ToAccount.Balance - toAccount.Balance
		require.Equal(t, fromDiff, toDiff)

	}

	fAccount, _ := testStore.GetAccount(context.Background(), fromAccount.ID)
	tAccount, _ := testStore.GetAccount(context.Background(), toAccount.ID)

	// balalnce after all operations
	require.Equal(t, fAccount.Balance, fromAccount.Balance-int64(n)*amount)
	require.Equal(t, tAccount.Balance, toAccount.Balance+int64(n)*amount)
}
