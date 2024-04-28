package db_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	db "github.com/kevenmarion/backend_master_class/db/sqlc"
	mocks "github.com/kevenmarion/backend_master_class/mocks/db/sqlc"
	"github.com/kevenmarion/backend_master_class/util"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var accountIDCount int64

func createRandomAccount(t *testing.T) db.Account {
	accountIDCount++
	var (
		arg = db.CreateAccountParams{
			Owner:    util.RandomOwner(),
			Balance:  util.RandomMoney(),
			Currency: util.RandomCurrency(),
		}
		mockQueries = mocks.NewQuerier(t)
		mockAccount = db.Account{
			ID:        accountIDCount,
			Owner:     arg.Owner,
			Balance:   arg.Balance,
			Currency:  arg.Currency,
			CreatedAt: time.Now(),
		}
	)

	mockQueries.On("CreateAccount", mock.Anything, mock.Anything).
		Return(mockAccount, nil)
	result, err := mockQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Owner, result.Owner)
	require.Equal(t, arg.Balance, result.Balance)
	require.Equal(t, arg.Currency, result.Currency)

	require.NotZero(t, result.ID)
	require.NotZero(t, result.CreatedAt)

	return result
}

func TestCreateAccount(t *testing.T) {
	_ = createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	var (
		mockQueries = mocks.NewQuerier(t)
		mockAccount = createRandomAccount(t)
	)
	mockQueries.On("GetAccount", mock.Anything, mock.Anything).
		Return(mockAccount, nil)
	result, err := mockQueries.GetAccount(context.Background(), mockAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, mockAccount.ID, result.ID)
	require.Equal(t, mockAccount.Owner, result.Owner)
	require.Equal(t, mockAccount.Balance, result.Balance)
	require.Equal(t, mockAccount.Currency, result.Currency)
	require.WithinDuration(t, mockAccount.CreatedAt, result.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	var (
		mockAccount = createRandomAccount(t)
		arg         = db.UpdateAccountParams{
			ID:      mockAccount.ID,
			Balance: 100,
		}
		mockQueries        = mocks.NewQuerier(t)
		mockAccountUpdated = db.Account{
			ID:        mockAccount.ID,
			Owner:     mockAccount.Owner,
			Balance:   arg.Balance,
			Currency:  mockAccount.Currency,
			CreatedAt: mockAccount.CreatedAt,
		}
	)
	mockQueries.On("UpdateAccount", mock.Anything, mock.Anything).
		Return(mockAccountUpdated, nil)
	result, err := mockQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, mockAccount.ID, result.ID)
	require.Equal(t, mockAccount.Owner, result.Owner)
	require.Equal(t, arg.Balance, result.Balance)
	require.Equal(t, mockAccount.Currency, result.Currency)
	require.WithinDuration(t, mockAccount.CreatedAt, result.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	var (
		mockAccount = createRandomAccount(t)
		mockQueries = mocks.NewQuerier(t)
	)

	mockQueries.
		On("DeleteAccount", mock.Anything, mock.Anything).Return(nil).
		On("GetAccount", mock.Anything, mock.Anything).
		Return(db.Account{}, sql.ErrNoRows)

	err := mockQueries.DeleteAccount(context.Background(), mockAccount.ID)
	require.NoError(t, err)

	result, err := mockQueries.GetAccount(context.Background(), mockAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, result)
}

func TestListAccounts(t *testing.T) {
	var (
		lastAccount  db.Account
		mockQueries  = mocks.NewQuerier(t)
		mockAccounts = []db.Account{}
	)

	for i := 0; i < 5; i++ {
		lastAccount = createRandomAccount(t)
		mockAccounts = append(mockAccounts, lastAccount)
	}

	mockQueries.On("ListAccounts", mock.Anything, mock.Anything).
		Return(mockAccounts, nil)

	arg := db.ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}
	accounts, err := mockQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	require.NotEmpty(t, accounts[len(accounts)-1])
	require.Equal(t, lastAccount.Owner, accounts[len(accounts)-1].Owner)
}
