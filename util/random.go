package util

import (
	"fmt"
	"github.com/CrunchyBlue/Golang-Bank/constants"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomString(10))
}

func RandomOwner() string {
	return RandomString(10)
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomCurrency() string {
	currencies := []string{
		constants.USD,
		constants.EUR,
		constants.CAD,
	}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
