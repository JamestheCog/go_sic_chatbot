package app

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
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

func (a *App) fetchResponse(history []*genai.Content, numRetries int, ctx context.Context, tokens []string, prompt string) (string, error) {
	contentConfig := &genai.GenerateContentConfig{
		MaxOutputTokens:   tokenLimit,
		SystemInstruction: genai.NewContentFromText(prompt, genai.RoleUser),
	}

	for i, token := range tokens {
		clientObj, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: token})
		if err != nil {
			return "", err
		}
		for j := 0; j < numRetries; j++ {
			response, err := clientObj.Models.GenerateContent(
				ctx, modelName, history, contentConfig,
			)
			if err != nil {
				if strings.Contains(err.Error(), "503") {
					if j == (numRetries - 1) {
						return "", fmt.Errorf("Gemini is currently swamped; try again later.")
					}
					log.Printf("Error 503 encountered; backing off before trying again (attempt #%d)...", j+1)
					delayDuration := math.Pow(float64(baseDelay), float64(j+1)) + (rand.Float64() * baseDelay)
					time.Sleep(time.Duration(delayDuration) * time.Second)
					continue
				} else if strings.Contains(err.Error(), "429") {
					if i < len(tokens) {
						log.Printf("Token %d is exhausted; switching to token %d now...\n", i+1, i+2)
					}
					break
				} else {
					return "", fmt.Errorf("Could not fetch Gemini's response for the following reason: %v", err)
				}
			} else {
				return response.Text(), nil
			}
		}
	}
	return "", fmt.Errorf("All tokens have been exhausted.")
}
