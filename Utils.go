package main

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		b[i] = charset[num.Int64()]
	}
	return string(b)
}

func RemoveSpacesAndColons(input string) string {
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, ":", "")
	return input

}

func RemovePeriods(input string) string {
	input = strings.ReplaceAll(input, ".", "")
	input = strings.ReplaceAll(input, " ", "")

	return input
}
