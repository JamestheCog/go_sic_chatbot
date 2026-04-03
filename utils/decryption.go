// A file I came up with to store decryption-related helper functions and / or constants.

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// Given a piece of cipher text as a plain string and an AES-GCM key
func DecryptString(cipherText string, key string) (string, error) {
	cipherBytes, err := fetchB64Encoding(cipherText)
	if err != nil {
		return "", err
	}
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
	nonceSize := gcm.NonceSize()
	if len(cipherBytes) < nonceSize {
		return "", fmt.Errorf("Ciphertext too short.")
	}

	nonce, actualCiphertext := cipherBytes[:nonceSize], cipherBytes[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

func fetchB64Encoding(input string) ([]byte, error) {
	if output, err := base64.StdEncoding.DecodeString(input); err != nil {
		if output, err = base64.URLEncoding.DecodeString(input); err != nil {
			if output, err = base64.RawStdEncoding.DecodeString(input); err != nil {
				return nil, fmt.Errorf("All three possible decryption modes for b64 strings have failed.")
			} else {
				return output, nil
			}
		} else {
			return output, nil
		}
	} else {
		return output, nil
	}
}
