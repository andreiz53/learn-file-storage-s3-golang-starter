package main

import (
	"math/rand"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz1234567890"

func RandomString(n int) string {
	var str string

	for i := 0; i < n; i++ {
		str += string(alphabet[rand.Intn(len(alphabet))])
	}
	return str
}
