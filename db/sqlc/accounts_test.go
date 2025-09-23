package db

import (
	"context"
	"testing"

	"example.com/db/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {

	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	_, err = account.Balance.Float64Value()
	require.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {
	var id int64 = 1

	err := testQueries.DeleteAccount(context.Background(), id)
	require.NoError(t, err)

}

func TestGetAccount(t *testing.T) {
	var id int64 = 2

	account, err := testQueries.GetAccount(context.Background(), id)

	require.NoError(t, err)

	require.NotEmpty(t, account.Balance)
	require.NotEmpty(t, account.CreatedAt)
	require.NotEmpty(t, account.Currency)
	require.NotEmpty(t, account.Owner)
	require.NotEmpty(t, account.ID)

}

// func (q *Queries) UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) error {
// 	_, err := q.db.Exec(ctx, updateAccountBalance, arg.ID, arg.Balance)
// 	return err
// }

func TestUpdateBalance(t *testing.T) {
	var id int64 = 2
	var newBalance pgtype.Numeric
	err := newBalance.Scan("22.2")
	require.NoError(t, err)

	arg := UpdateAccountBalanceParams{
		ID:      id,
		Balance: newBalance,
	}

	err = testQueries.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), id)
	require.NoError(t, err)

	expectedFloat, err := newBalance.Float64Value()
	require.NoError(t, err)
	actualFloat, err := account.Balance.Float64Value()
	require.NoError(t, err)
	require.Equal(t, expectedFloat.Float64, actualFloat.Float64)

}

func CreateRandomAccount(ctx context.Context) (Account, error) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(ctx, arg)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}
