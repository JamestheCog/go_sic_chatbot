package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

// Basicallt the inverse of the DecryptFile function - it encodes stuff.
// I don't know why we'd need it, but it's better to have and not need than to
// need and not have if you ask me:
func EncryptString(plainText, key string) (string, error) {
	plainBytes := []byte(plainText)
	byteKey, err := hex.DecodeString(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(byteKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, plainBytes, nil)
	toReturn := base64.StdEncoding.EncodeToString(cipherText)
	return toReturn, nil
}
