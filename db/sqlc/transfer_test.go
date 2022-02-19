package db

import (
	"context"
	"testing"

	"github.com/hamdysherif/simplebank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	args := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomeBalance(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.Amount, transfer.Amount)
	return transfer
}
func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	trans := createRandomTransfer(t)
	fTrans, err := testQueries.GetTransfer(context.Background(), trans.ID)

	require.NoError(t, err)
	require.NotEmpty(t, fTrans)
	require.Equal(t, trans.ID, fTrans.ID)
}

func TestListTransfers(t *testing.T) {
	args := ListTransfersParams{
		Offset: 2,
		Limit:  7,
	}
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)

	assert.NoError(t, err)
	assert.Len(t, transfers, 7)

	for _, trans := range transfers {
		assert.NotEmpty(t, trans)
	}
}
