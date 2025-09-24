package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HashKey(plainKey string) string {
	sum := sha256.Sum256([]byte(plainKey))
	return hex.EncodeToString(sum[:])
}

func GenerateAPIKey() (plainKey, hashedKey string, err error) {
	bytes := make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	plainKey = hex.EncodeToString(bytes)
	sum := sha256.Sum256([]byte(plainKey))
	hashedKey = hex.EncodeToString(sum[:])
	return plainKey, hashedKey, nil
}