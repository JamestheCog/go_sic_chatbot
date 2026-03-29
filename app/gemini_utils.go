package app

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"google.golang.org/genai"
)

// --- Gemini-related constants here ---
const (
	modelName  = "gemini-3-flash-preview"
	tokenLimit = 2000
	baseDelay  = 3
)

// --- END ---

// Fetches the prompt that we're going to be using to make the bot behave as is (i.e.,
// the very same prompt we passed off to Dylan on the SUTD's side of things).
func (a *App) decryptFileToText(filePath string) (string, error) {
	ciphertextB64, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	cleanB64 := strings.TrimSpace(string(ciphertextB64))
	cipherBytes, err := base64.StdEncoding.DecodeString(cleanB64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}
	byteKey, err := hex.DecodeString(a.AesKey)
	if err != nil {
		return "", fmt.Errorf("invalid hex key: %w", err)
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
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, actualCiphertext := cipherBytes[:nonceSize], cipherBytes[nonceSize:]
	plainBytes, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed (check your key): %w", err)
	}
	return string(plainBytes), nil
}

func (a *App) createContentConfig(filePath string) (*genai.GenerateContentConfig, error) {
	decryptedPrompt, err := a.decryptFileToText(filePath)
	if err != nil {
		return nil, err
	}

	return &genai.GenerateContentConfig{
		MaxOutputTokens:   tokenLimit,
		SystemInstruction: genai.NewContentFromText(decryptedPrompt, genai.RoleUser),
	}, nil
}

// The helper function of interest that tries to - given a history - fetch Gemini's response.
func (a *App) fetchResponse(history []*genai.Content, numRetries int, ctx context.Context) (string, error) {
	for i := 0; i < numRetries; i++ {
		response, err := a.GeminiClient.Models.GenerateContent(
			ctx, modelName,
			history, a.ContentConfig,
		)

		if err != nil {
			if strings.Contains(err.Error(), "429") {
				return "", fmt.Errorf("Kevin, we've a problem.  The token has called for MSW's services!")
			} else if strings.Contains(err.Error(), "503") {
				delayDuration := math.Pow(float64(baseDelay), float64(i+1))
				log.Printf("Gemini's too busy right now - waiting %.2f seconds before trying again (attempt #%d out of %d)...",
					delayDuration, i+1, numRetries)
				time.Sleep(time.Duration(delayDuration+(rand.Float64()*delayDuration)) * time.Second)
				continue
			}
			return "", fmt.Errorf("Gemini is too busy right now; try again later perhaps?")
		}
		return response.Text(), nil
	}
	return "", fmt.Errorf("Could not fetch any responses from Gemini for some reason.")
}
