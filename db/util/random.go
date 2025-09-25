package util

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() pgtype.Numeric {
	value := RandomInt(0, 1000)

	var numeric pgtype.Numeric
	// Format as string and let Scan handle conversion
	if err := numeric.Scan(fmt.Sprintf("%d", value)); err != nil {
		panic(err)
	}

	return numeric
}

func RandomCurrency() string {
	currencies := []string{"EUR", "BHD", "USD", "AED", "SAR", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
