package main

import (
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

func main() {
	amount := int64(10)

	amountNumeric := pgtype.Numeric{
		Int:   big.NewInt(amount),
		Exp:   0,
		Valid: true,
	}

	fmt.Println(amountNumeric.Int.Int64())

}
