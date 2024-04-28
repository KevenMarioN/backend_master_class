package db_test

import (
	"context"
	"testing"
	"time"

	db "github.com/kevenmarion/backend_master_class/db/sqlc"
	mocks "github.com/kevenmarion/backend_master_class/mocks/db/sqlc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var transferIDCount int64

func createRandomTransfer(t *testing.T, to_account, from_account db.Account) (transfer db.Transfer) {
	transferIDCount++
	var (
		arg = db.CreateTransferParams{
			Amount:        100,
			ToAccountID:   to_account.ID,
			FromAccountID: from_account.ID,
		}
		mockTransfer = db.Transfer{
			ID:            transferIDCount,
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
			CreatedAt:     time.Now(),
		}
		mockQueries = mocks.NewQuerier(t)
	)

	mockQueries.On("CreateTransfer", mock.Anything, mock.Anything).
		Return(mockTransfer, nil)

	result, err := mockQueries.CreateTransfer(context.Background(), arg)
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
	mockAccountTo := createRandomAccount(t)
	mockAccountFrom := createRandomAccount(t)

	createRandomTransfer(t, mockAccountTo, mockAccountFrom)
}

func TestGetTransfer(t *testing.T) {
	var (
		mockAccountTo   = createRandomAccount(t)
		mockAccountFrom = createRandomAccount(t)
		mockTransfer    = createRandomTransfer(t, mockAccountTo, mockAccountFrom)
		mockQueries     = mocks.NewQuerier(t)
	)
	mockQueries.On("GetTransfer", mock.Anything, mock.Anything).
		Return(mockTransfer, nil)

	result, err := mockQueries.GetTransfer(context.Background(), mockTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, mockTransfer.ID, result.ID)
	require.Equal(t, mockTransfer.Amount, result.Amount)
	require.Equal(t, mockTransfer.ToAccountID, result.ToAccountID)
	require.Equal(t, mockTransfer.FromAccountID, result.FromAccountID)
	require.Equal(t, mockTransfer.CreatedAt, result.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	var (
		account_1     = createRandomAccount(t)
		account_2     = createRandomAccount(t)
		mockQueries   = mocks.NewQuerier(t)
		mockTransfers = []db.Transfer{}
	)

	for i := 0; i < 10; i++ {
		mockTransfers = append(mockTransfers, createRandomTransfer(t, account_1, account_2))
		mockTransfers = append(mockTransfers, createRandomTransfer(t, account_2, account_1))
	}

	arg := db.ListTransfersParams{
		FromAccountID: account_1.ID,
		ToAccountID:   account_1.ID,
		SizeLimit:     5,
		SizeOffset:    5,
	}

	mockQueries.On("ListTransfers", mock.Anything, mock.Anything).
		Return(mockTransfers, nil)

	transfers, err := mockQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 20)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.ToAccountID == account_1.ID || transfer.FromAccountID == account_1.ID)
	}

}
