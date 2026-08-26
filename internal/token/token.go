package token

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func Generate() (string, error) {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read secure random bytes: %w", err)
	}
	return hex.EncodeToString(buf), nil
}
