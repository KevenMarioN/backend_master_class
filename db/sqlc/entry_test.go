package db_test

import (
	"context"
	"testing"
	"time"

	db "github.com/kevenmarion/backend_master_class/db/sqlc"
	mocks "github.com/kevenmarion/backend_master_class/mocks/db/sqlc"
	"github.com/kevenmarion/backend_master_class/util"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var entryIDCount int64

func createRandomEntry(t *testing.T) db.Entry {
	entryIDCount++
	var (
		mockAccount = createRandomAccount(t)
		arg         = db.CreateEntryParams{
			AccountID: mockAccount.ID,
			Amount:    util.RandomMoney(),
		}
		mockQueries = mocks.NewQuerier(t)
		mockEntry   = db.Entry{
			ID:        entryIDCount,
			AccountID: arg.AccountID,
			Amount:    arg.Amount,
			CreatedAt: time.Now(),
		}
	)
	mockQueries.On("CreateEntry", mock.Anything, mock.Anything).
		Return(mockEntry, nil)
	entry, err := mockQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.Amount, arg.Amount)
	require.Equal(t, entry.AccountID, arg.AccountID)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	var (
		mockEntry   = createRandomEntry(t)
		mockQueries = mocks.NewQuerier(t)
	)

	mockQueries.On("GetEntry", mock.Anything, mock.Anything).
		Return(mockEntry, nil)

	result, err := mockQueries.GetEntry(context.Background(), mockEntry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, mockEntry.ID, result.ID)
	require.Equal(t, mockEntry.Amount, result.Amount)
	require.Equal(t, mockEntry.AccountID, result.AccountID)
	require.Equal(t, mockEntry.CreatedAt, result.CreatedAt)
}

func TestListEntries(t *testing.T) {
	var (
		lastEntry   db.Entry
		mockQueries = mocks.NewQuerier(t)
		mockEntries = []db.Entry{}
	)
	for i := 0; i < 5; i++ {
		lastEntry = createRandomEntry(t)
		mockEntries = append(mockEntries, lastEntry)
	}

	mockQueries.On("ListEntries", mock.Anything, mock.Anything).
		Return(mockEntries, nil)

	arg := db.ListEntriesParams{
		SizeLimit:  5,
		SizeOffset: 0,
		AccountID:  lastEntry.AccountID,
	}
	results, err := mockQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, results)

	require.NotEmpty(t, results[len(results)-1])
	require.Equal(t, lastEntry.ID, results[len(results)-1].ID)
}
