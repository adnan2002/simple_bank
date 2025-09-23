package db

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	ctx := context.Background()

	// create two accounts
	account1, err := CreateRandomAccount(ctx)
	require.NoError(t, err)
	account2, err := CreateRandomAccount(ctx)
	require.NoError(t, err)

	// get initial balances
	initialBalance1 := account1.Balance
	initialBalance2 := account2.Balance

	n := 2
	amount := int64(10)

	type TransferTxResultAndErr struct {
		Result TransferTxResult
		Err    error
	}
	resErrChan := make(chan TransferTxResultAndErr, n)

	// run n concurrent transfers
	for i := 0; i < n; i++ {
		go func() {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})
			resErrChan <- TransferTxResultAndErr{Result: res, Err: err}
		}()
	}

	// verify all transfers
	for i := 0; i < n; i++ {
		resErr := <-resErrChan
		require.NoError(t, resErr.Err)

		resp := resErr.Result
		require.NotEmpty(t, resp)

		transfer := resp.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)

		actual, err := convertToInt64(transfer.Amount)
		require.NoError(t, err)
		require.Equal(t, amount, actual)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// entries check
		fromEntry := resp.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)

		toEntry := resp.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)

	}

	// --- Final Balance Check ---
	updatedAccount1, err := store.GetAccount(ctx, account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(ctx, account2.ID)
	require.NoError(t, err)

	// expected balances
	totalTransferred := amount * int64(n)

	expectedBalance1 := mustSub(initialBalance1, totalTransferred)
	expectedBalance2 := mustAdd(initialBalance2, totalTransferred)

	require.Equal(t, expectedBalance1, updatedAccount1.Balance)
	require.Equal(t, expectedBalance2, updatedAccount2.Balance)
}

// --- Helpers ---

func convertToInt64(num pgtype.Numeric) (int64, error) {
	if !num.Valid {
		return 0, errors.New("input pgtype.Numeric is not valid")
	}
	val := new(big.Int).Set(num.Int)
	if num.Exp > 0 {
		val.Mul(val, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(num.Exp)), nil))
	} else if num.Exp < 0 {
		val.Quo(val, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-num.Exp)), nil))
	}
	if val.IsInt64() {
		return val.Int64(), nil
	}
	return 0, errors.New("conversion to int64 would cause an overflow")
}

func mustSub(n pgtype.Numeric, amount int64) pgtype.Numeric {
	newN, err := SubtractNumericInt64(n, amount)
	if err != nil {
		panic(err)
	}
	return newN
}

func mustAdd(n pgtype.Numeric, amount int64) pgtype.Numeric {
	newN, err := AddNumericInt64(n, amount)
	if err != nil {
		panic(err)
	}
	return newN
}
