package util

import (
	"crypto/rand"
	"math/big"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandShortlyCode(n int) string {
	b := make([]rune, n)
	for i := range b {
		val, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		b[i] = letterRunes[val.Int64()]
	}
	return string(b)
}
