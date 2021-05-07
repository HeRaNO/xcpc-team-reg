package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

const sigma = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Generate a token whose length is `n`
func GenToken(n int) (string, error) {
	b := make([]byte, n)
	rng := new(big.Int).SetInt64(int64(len(sigma)))
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, rng)
		if err != nil {
			return "", err
		}
		b[i] = sigma[idx.Int64()]
	}
	return string(b), nil
}

// Get the SHA256 of `msg`
func SHA256(msg []byte) string {
	h := sha256.Sum256(msg)
	return hex.EncodeToString(h[:])
}
