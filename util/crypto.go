package util

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

const sigma = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const (
	letterIdxBits = 6
	letterIdxMask = 63
	letterIdxMax  = 10
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenToken(n int) string {
	b := make([]byte, n)
	for i, rnd, left := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if left == 0 {
			rnd, left = rand.Int63(), letterIdxMax
		}
		if idx := int(rnd & letterIdxMask); idx < len(sigma) {
			b[i] = sigma[idx]
			i--
		}
		rnd >>= letterIdxBits
		left--
	}
	return string(b)
}

func SHA256(msg []byte) string {
	h := sha256.Sum256(msg)
	return hex.EncodeToString(h[:])
}
