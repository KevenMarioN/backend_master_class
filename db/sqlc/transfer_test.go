package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, to_account, from_account Account) (transfer Tranfer) {
	arg := CreateTransferParams{
		Amount:        100,
		ToAccountID:   to_account.ID,
		FromAccountID: from_account.ID,
	}
	result, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, result.Amount, arg.Amount)
	require.Equal(t, result.ToAccountID, arg.ToAccountID)
	require.Equal(t, result.FromAccountID, arg.FromAccountID)

	require.NotZero(t, result.ID)
	require.NotZero(t, result.CreatedAt)
	return result
}

func TestCreateTransfer(t *testing.T) {
	account_to := createRandomAccount(t)
	account_from := createRandomAccount(t)

	createRandomTransfer(t, account_to, account_from)
}

func TestGetTransfer(t *testing.T) {
	account_to := createRandomAccount(t)
	account_from := createRandomAccount(t)

	transfer := createRandomTransfer(t, account_to, account_from)

	result, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, transfer.ID, result.ID)
	require.Equal(t, transfer.Amount, result.Amount)
	require.Equal(t, transfer.ToAccountID, result.ToAccountID)
	require.Equal(t, transfer.FromAccountID, result.FromAccountID)
	require.Equal(t, transfer.CreatedAt, result.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	account_1 := createRandomAccount(t)
	account_2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, account_1, account_2)
		createRandomTransfer(t, account_2, account_1)
	}

	arg := ListTransfersParams{
		FromAccountID: account_1.ID,
		ToAccountID:   account_1.ID,
		SizeLimit:     5,
		SizeOffset:    5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.ToAccountID == account_1.ID || transfer.FromAccountID == account_1.ID)
	}

}
