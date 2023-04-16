package util

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const numbers = "123456789"

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

func RandomMoney() int64 {
	return rand.Int63n(10000)
}

func RandomCurrency() string {
	currencies := []string{
		"EUR",
		"USD",
		"CAD",
	}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
