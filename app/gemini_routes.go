// A file I came up with to store handler functions that concern itself with Gemini's API.

package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

// --- Structs and Constants ---

const numTries = 5

// The expected of the incoming payload from JavaScript that
// we're expecting to see:
type MsgPayload struct {
	Session ConversationDetails `json:"session"`
	Message ChatItem            `json:"message"`
}

type ConversationDetails struct {
	ConversationID string `json:"conversation_id"`
}

type ChatItem struct {
	Msg     string `json:"msg"`
	ImgB64  string `json:"img_b64"`
	ImgType string `json:"img_mime"`
	Role    string `json:"role"`
}

// --- End ---

// --- The handlers --

// Coordinates the execution of the various app-based helper functions for uploading and
// formatting messages and image strings to be sent to Gemini's API.
//
// Note that the session's variables and then some are going to be initialized in the
// template route - if any of the variables are missing, then this function's gonna
// throw.
func (a *App) messageHandler(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("Only POSTs be allowed here.")
	}

	val := r.Context().Value(cookieName)
	if val == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("No session available.")
	}
	session := val.(*sessions.Session)

	// Fetch our session and conversation IDs:
	sessionID, ok := session.Values[sessIdKey].(string)
	if !ok {
		return nil, http.StatusInternalServerError, fmt.Errorf("Cookie named '%s' not set.", sessIdKey)
	}

	var payload MsgPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Handle database stuff here:
	if err := a.uploadMessage(payload, payload.Session.ConversationID, sessionID); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	msgHistory, err := a.fetchConversationHistory(payload.Session.ConversationID, sessionID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	geminiMsg, err := a.fetchResponse(msgHistory, numTries, r.Context(), a.GeminiKeys, a.ConversationPrompt)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if err := a.uploadMessage(geminiMsg, payload.Session.ConversationID, sessionID); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return SimpleResponse{Message: geminiMsg}, http.StatusOK, nil
}
