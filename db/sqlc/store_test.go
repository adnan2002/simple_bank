package db

import (
	"context"
	"testing"
	"math/big"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)
func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	ctx := context.Background()

	account1, err := CreateRandomAccount(ctx)
	require.NoError(t, err)
	account2, err := CreateRandomAccount(ctx)
	require.NoError(t, err)


	n := 5
	amount := int64(10)

type TransferTxResultAndErr struct {
	Result TransferTxResult
	Err    error
}
resErrChan := make(chan TransferTxResultAndErr, n)

	for i := 0; i < n; i++ {
	go func() {
		res, err := store.TransferTx(ctx, TransferTxParams{
			FromAccountId: account1.ID,
			ToAccountId:   account2.ID,
			Amount:        amount,
		})
		resErrChan <- TransferTxResultAndErr{Result: res, Err: err}
	}()
}


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

	}

}


func convertToInt64(num pgtype.Numeric) (int64, error) {
	if !num.Valid {
		return 0, errors.New("input pgtype.Numeric is not valid")
	}

	// Calculate the value by considering the exponent.
	// We need to handle potential floating point inaccuracies by using a more robust approach.
	val := new(big.Int).Set(num.Int)
	if num.Exp > 0 {
		val.Mul(val, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(num.Exp)), nil))
	} else if num.Exp < 0 {
		val.Quo(val, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-num.Exp)), nil))
	}

	// Check for overflow before converting to int64.
	if val.IsInt64() {
		return val.Int64(), nil
	}

	return 0, errors.New("conversion to int64 would cause an overflow")
}
