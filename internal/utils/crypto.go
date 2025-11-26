package utils

import (
	"crypto/rand"
	"math/big"

	"github.com/HeRaNO/xcpc-team-reg/internal/berrors"
	"golang.org/x/crypto/bcrypt"
)

const sigma = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenToken generates a token whose length is `n`.
func GenToken(n int) (string, berrors.Berror) {
	b := make([]byte, n)
	rng := new(big.Int).SetInt64(int64(len(sigma)))
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, rng)
		if err != nil {
			return "", berrors.New(berrors.ErrInternal, err.Error())
		}
		b[i] = sigma[idx.Int64()]
	}
	return string(b), nil
}

// GenSecret generates a bytes-secret whose length is `n`.
func GenSecret(n int) ([]byte, error) {
	b := make([]byte, n)
	rng := new(big.Int).SetInt64(256)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, rng)
		if err != nil {
			return b, err
		}
		b[i] = byte(idx.Uint64())
	}
	return b, nil
}

func HashPassword(pwd *string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(*pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashStr := string(hash)
	return &hashStr, nil
}

func ValidatePassword(hashedPwd, plainPwd *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hashedPwd), []byte(*plainPwd))
	return err == nil
}
