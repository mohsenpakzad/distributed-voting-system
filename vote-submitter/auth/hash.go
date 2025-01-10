package auth

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return fmt.Sprintf("%x%x", salt, hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	salt, hashedPassword, err := splitHash(hash)
	if err != nil {
		return false
	}
	decodedSalt, err := hex.DecodeString(salt)
	if err != nil {
		return false
	}
	decodedHashedPassword, err := hex.DecodeString(hashedPassword)
	if err != nil {
		return false
	}
	calculatedHash := argon2.IDKey([]byte(password), decodedSalt, 1, 64*1024, 4, 32)
	return equal(calculatedHash, decodedHashedPassword)
}

func splitHash(hash string) (string, string, error) {
	if len(hash)%2 != 0 {
		return "", "", fmt.Errorf("invalid hash length")
	}
	salt := hash[:32]
	hashedPassword := hash[32:]
	return salt, hashedPassword, nil
}

func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
