package db

import (
	"context"
	"testing"

	"github.com/kevenmarion/backend_master_class/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.Amount, arg.Amount)
	require.Equal(t, entry.AccountID, arg.AccountID)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)

	result, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, entry.ID, result.ID)
	require.Equal(t, entry.Amount, result.Amount)
	require.Equal(t, entry.AccountID, result.AccountID)
	require.Equal(t, entry.CreatedAt, result.CreatedAt)
}

func TestListEntries(t *testing.T) {
	var lastEntry Entry
	for i := 0; i < 10; i++ {
		lastEntry = createRandomEntry(t)
	}

	arg := ListEntriesParams{
		SizeLimit:  5,
		SizeOffset: 0,
		AccountID:  lastEntry.AccountID,
	}
	results, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, results)

	for _, entry := range results {
		require.NotEmpty(t, entry)
		require.Equal(t, lastEntry.ID, entry.ID)
	}
}
