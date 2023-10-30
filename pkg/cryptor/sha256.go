package cryptor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// Sha256 кодируем пароль
func Sha256(src string) (string, error) {
	hash := sha256.New()

	if src == "" {
		return "", ErrEmptyPassword
	}

	_, err := io.WriteString(hash, src)
	if err != nil {
		return "", fmt.Errorf("write hash: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
