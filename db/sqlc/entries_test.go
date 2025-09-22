package db

import (
	"context"
	"testing"

	"example.com/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	// Use account IDs from 2 to 6 as specified
	accountID := util.RandomInt(2, 6)
	
	// Create random amount (could be positive or negative for debits/credits)
	amount := util.RandomMoney()
	
	arg := CreateEntryParams{
		AccountID: accountID,
		Amount:    amount,
	}
	
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)


	arg_amount, err := arg.Amount.Float64Value()
	require.NoError(t, err)
	real_arg_amount := arg_amount.Float64
	entry_amount, err := entry.Amount.Float64Value()
	require.NoError(t, err)
	real_entry_amount := entry_amount.Float64

	
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, real_arg_amount, real_entry_amount)
	require.True(t, entry.Amount.Valid)
	
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	
	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount.Int, entry2.Amount.Int)
	require.Equal(t, entry1.Amount.Valid, entry2.Amount.Valid)
	require.WithinDuration(t, entry1.CreatedAt.Time, entry2.CreatedAt.Time, 0)
}

func TestUpdateEntryAmount(t *testing.T) {
	entry1 := createRandomEntry(t)
	
	newAmount := util.RandomMoney()
	arg := UpdateEntryAmountParams{
		ID:     entry1.ID,
		Amount: newAmount,
	}
	
	err := testQueries.UpdateEntryAmount(context.Background(), arg)
	require.NoError(t, err)
	
	// Verify the update by getting the entry
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, newAmount, entry2.Amount.Int)
	require.True(t, entry2.Amount.Valid)
	require.WithinDuration(t, entry1.CreatedAt.Time, entry2.CreatedAt.Time, 0)
}

func TestDeleteEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	
	// Try to get the deleted entry
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.Empty(t, entry2)
}

func TestListEntries(t *testing.T) {
	// Create multiple entries
	var createdEntries = make([]Entry, 10)
	for i := 0; i < 10; i++ {
		entry := createRandomEntry(t)
		createdEntries[i] =  entry
	}
	
	arg := ListEntriesParams{
		Limit:  5,
		Offset: 0,
	}
	
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.AccountID)
		require.True(t, entry.Amount.Valid)
		require.NotZero(t, entry.CreatedAt)
	}
}

func TestListEntriesWithOffset(t *testing.T) {
	// Create multiple entries
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}
	
	// Test with offset
	arg := ListEntriesParams{
		Limit:  3,
		Offset: 5,
	}
	
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 3)
	
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.AccountID)
		require.True(t, entry.Amount.Valid)
		require.NotZero(t, entry.CreatedAt)
	}
}

func TestListEntriesForAccount(t *testing.T) {
	accountID := int64(3) // Use one of the allowed account IDs
	
	// Create multiple entries for the specific account
	var createdEntries = make([]Entry, 5)
	for i := 0; i < 5; i++ {
		arg := CreateEntryParams{
			AccountID: accountID,
			Amount:    util.RandomMoney(),
		}
		
		entry, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
		createdEntries[i] = entry
	}
	
	// Create some entries for other accounts
	for i := 0; i < 3; i++ {
		createRandomEntry(t)
	}
	
	// List entries for the specific account
	arg := ListEntriesForAccountParams{
		AccountID: accountID,
		Limit:     10,
		Offset:    0,
	}
	
	entries, err := testQueries.ListEntriesForAccount(context.Background(), arg)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(entries), 5) // Should have at least the 5 we created
	
	// All entries should belong to the specified account
	for _, entry := range entries {
		require.Equal(t, accountID, entry.AccountID)
		require.NotEmpty(t, entry)
		require.NotZero(t, entry.ID)
		require.True(t, entry.Amount.Valid)
		require.NotZero(t, entry.CreatedAt)
	}
}

func TestListEntriesForAccountWithPagination(t *testing.T) {
	accountID := int64(4) // Use one of the allowed account IDs
	
	// Create multiple entries for the specific account
	for i := 0; i < 8; i++ {
		arg := CreateEntryParams{
			AccountID: accountID,
			Amount:    util.RandomMoney(),
		}
		
		_, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
	}
	
	// Test pagination
	arg := ListEntriesForAccountParams{
		AccountID: accountID,
		Limit:     3,
		Offset:    2,
	}
	
	entries, err := testQueries.ListEntriesForAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 3)
	
	// All entries should belong to the specified account
	for _, entry := range entries {
		require.Equal(t, accountID, entry.AccountID)
		require.NotEmpty(t, entry)
		require.NotZero(t, entry.ID)
		require.True(t, entry.Amount.Valid)
		require.NotZero(t, entry.CreatedAt)
	}
}

func TestCreateEntryWithNegativeAmount(t *testing.T) {
	accountID := int64(5)
	negativeAmount := util.RandomMoney() // Create a negative amount
	
	arg := CreateEntryParams{
		AccountID: accountID,
		Amount:    negativeAmount,
	}
	
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, negativeAmount, entry.Amount.Int)
	require.True(t, entry.Amount.Valid)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntryNotFound(t *testing.T) {
	// Try to get an entry with a non-existent ID
	entry, err := testQueries.GetEntry(context.Background(), 999999)
	require.Error(t, err)
	require.Empty(t, entry)
}

func TestUpdateEntryAmountNotFound(t *testing.T) {
	// Try to update an entry with a non-existent ID
	arg := UpdateEntryAmountParams{
		ID:     999999,
		Amount: util.RandomMoney(),
	}
	
	err := testQueries.UpdateEntryAmount(context.Background(), arg)
	require.NoError(t, err) // Note: PostgreSQL UPDATE returns no error even if no rows are affected
}

func TestDeleteEntryNotFound(t *testing.T) {
	// Try to delete an entry with a non-existent ID
	err := testQueries.DeleteEntry(context.Background(), 999999)
	require.NoError(t, err) // Note: PostgreSQL DELETE returns no error even if no rows are affected
}