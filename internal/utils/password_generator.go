package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/CSR-LC/csr-be/pkg/domain"
)

const (
	MinPasswordLen                    = 8
	MaxPasswordLen                    = 32
	AllowedRandomResetPasswordSymbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type passwordGenerator struct {
	length int
}

func NewPasswordGenerator(length int) (domain.PasswordGenerator, error) {
	if length < MinPasswordLen {
		return nil, fmt.Errorf("password length must be at least %d", MinPasswordLen)
	}
	if length > MaxPasswordLen {
		return nil, fmt.Errorf("password length must be at most %d", MaxPasswordLen)
	}
	return &passwordGenerator{length: length}, nil
}

func (p passwordGenerator) NewPassword() (string, error) {
	return generateRandomString(p.length, AllowedRandomResetPasswordSymbols)
}

func generateRandomString(n int, symbols string) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(symbols))))
		if err != nil {
			return "", err
		}
		ret[i] = symbols[num.Int64()]
	}

	return string(ret), nil
}
