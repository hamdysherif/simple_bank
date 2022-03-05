package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabit = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generate randome number between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generate randome string of n length
func RandomString(n int) string {
	var sb strings.Builder
	k := int64(len(alphabit))

	for i := 0; i < n; i++ {
		c := alphabit[rand.Int63n(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(8)
}

func RandomeBalance() int64 {
	return RandomInt(0, 10000)
}

func RandomEntryAmount() int64 {
	return RandomInt(-1000, 1000)
}

func RandomCurrency() string {
	return AllowedCurrencies()[rand.Intn(len(AllowedCurrencies()))]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@bank.com", RandomString(7))
}
