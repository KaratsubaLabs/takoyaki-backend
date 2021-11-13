package main

import (
	"time"
	"math/rand"
)

func ContainsString(slice []string, elem string) bool {
	for _, a := range slice {
		if a == elem {
			return true
		}
	}
	return false
}

const randomStringLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const randomStringLen = 10
func RandomString() string {

	rand.Seed(time.Now().UnixNano())

	out := make([]byte, randomStringLen)
	for i := range out {
		out[i] = randomStringLetters[rand.Intn(len(randomStringLetters))]
	}
	return string(out)

}

