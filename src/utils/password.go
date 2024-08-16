package utils

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

const LETTER_BYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func HashPassword(password string) ([]byte, error) {
	p := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func VerifyPassword(password []byte, confirmPassword []byte) error {
	if err := bcrypt.CompareHashAndPassword(password, confirmPassword); err != nil {
		return err
	}
	return nil
}

func GenerateRandomPassword(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LETTER_BYTES[rand.Intn(len(LETTER_BYTES))]
	}
	return string(b)
}
