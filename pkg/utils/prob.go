package utils

import (
	"crypto/rand"
	"math/big"
)

func Prob(p float64) bool {
	// Generate a random number in [0, 1)
	n, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		panic(err)
	}
	r := float64(n.Int64()) / 100.0

	return r < p
}
