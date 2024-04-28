package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kevenmarion/backend_master_class/util"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	result, err := testQueries.CreateAccount(context.Background(), arg)
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
	account := createRandomAccount(t)
	result, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, account.ID, result.ID)
	require.Equal(t, account.Owner, result.Owner)
	require.Equal(t, account.Balance, result.Balance)
	require.Equal(t, account.Currency, result.Currency)
	require.WithinDuration(t, account.CreatedAt, result.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: 100,
	}
	result, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, account.ID, result.ID)
	require.Equal(t, account.Owner, result.Owner)
	require.Equal(t, arg.Balance, result.Balance)
	require.Equal(t, account.Currency, result.Currency)
	require.WithinDuration(t, account.CreatedAt, result.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	result, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, result)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  10,
		Offset: 0,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
