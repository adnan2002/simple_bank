package db

import (
	"context"
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(context.Background()); rbErr != nil {
			return fmt.Errorf("tx err: %v", rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		amountNumeric := pgtype.Numeric{
			Int:   big.NewInt(arg.Amount),
			Exp:   0,
			Valid: true,
		}

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        amountNumeric,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount: pgtype.Numeric{
				Int:   new(big.Int).Neg(big.NewInt(arg.Amount)),
				Exp:   0,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    amountNumeric,
		})
		if err != nil {
			return err
		}

		err = q.SubtractAccountBalance(ctx, SubtractAccountBalanceParams{
			ID:     arg.FromAccountId,
			Amount: amountNumeric,
		})

		if err != nil {
			return err
		}

		err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountId,
			Amount: amountNumeric,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
func SubtractNumericInt64(n pgtype.Numeric, amount int64) (pgtype.Numeric, error) {
	// If the numeric is NULL
	if !n.Valid {
		return pgtype.Numeric{Valid: false}, nil
	}

	// Convert amount into a scaled big.Int
	// n.Int holds the unscaled integer, and n.Exp is the base-10 exponent (negative for decimals)
	scale := big.NewInt(1)
	if n.Exp < 0 {
		scale.Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil)
	}

	// amount * scale
	amt := big.NewInt(amount)
	amt.Mul(amt, scale)

	// Copy the balance to avoid mutating original
	newInt := new(big.Int).Sub(n.Int, amt)

	// Build result
	return pgtype.Numeric{
		Int:   newInt,
		Exp:   n.Exp,
		Valid: true,
	}, nil
}

func AddNumericInt64(n pgtype.Numeric, amount int64) (pgtype.Numeric, error) {
	// If the numeric is NULL
	if !n.Valid {
		return pgtype.Numeric{Valid: false}, nil
	}

	// Convert amount into a scaled big.Int
	scale := big.NewInt(1)
	if n.Exp < 0 {
		scale.Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil)
	}

	amt := big.NewInt(amount)
	amt.Mul(amt, scale)

	// Copy to avoid mutating original
	newInt := new(big.Int).Add(n.Int, amt)

	// Build result
	return pgtype.Numeric{
		Int:   newInt,
		Exp:   n.Exp,
		Valid: true,
	}, nil
}
