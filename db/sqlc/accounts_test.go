package db

import (
	"context"

	"example.com/db/util"
)


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
