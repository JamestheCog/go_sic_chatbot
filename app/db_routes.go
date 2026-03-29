// A file I came up with to store routes that've to do with the database directly.
// There's not much in here for now (i.e., Sunday, 22nd March, 2026), but things
// might explode in the future depending on feature requests, so...
//
// Just keep this file in handy!

package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jamesthecog/go_chatbot/utils"
)

// --- Constants and structs here... ---
//
// Note: I'm going to use a one-off for incoming JSON payloads - conversation IDs
// 		 sent by the front-end to the respective routes - by using a map[string]string
//		 to marshal the data into.

// The struct of interest for passing items to the application's front-end -
// so that JavaScript can then format those messages accordingly in the
// chat container:
type ChatMessage struct {
	Message      string `json:"message"`
	Role         string `json:"role"`
	B64ImgString string `json:"img_b64"`
	ImgType      string `json:"img_mimetype"`
}

// The struct of interest to send over for the preview.
type ChatPreview struct {
	Date           string `json:"date_sent"`
	ConversationID string `json:"conversation_id"`
}

// What we'll be sending over if we want to delete a
// conversation:
type DeleteConvo struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
}

// --- End ---

// The route of interest for fetching messages given a conversation and a
// session ID:
func (a *App) fetchMessages(w http.ResponseWriter, r *http.Request) (any, int, error) {
	val := r.Context().Value(cookieName)
	if val == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("No session available.")
	}
	session := val.(*sessions.Session)

	sessionID, ok := session.Values[sessIdKey].(string)
	if !ok {
		return nil, http.StatusInternalServerError, fmt.Errorf("Cookie named '%s' not set.", sessIdKey)
	}
	var jsonPayload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&jsonPayload); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Wack-ass JSON encountered.  Aborting.")
	}
	conversationID, ok := jsonPayload["conversation_id"]
	if !ok {
		return nil, http.StatusBadRequest, fmt.Errorf("No `conversationID` field lah, bodoh!")
	}

	const selectStatement = "SELECT message, role, img_b64, img_mimetype " +
		"FROM conversations " +
		"WHERE conversation_id = ? AND session_id = ? " +
		"ORDER BY time_sent;"
	vals := []any{conversationID, sessionID}
	messages, err := a.DbClient.SelectArray(selectStatement, vals)
	if err != nil {
		fmt.Println(messages)
		return nil, http.StatusInternalServerError, err
	}

	if messages.GetNumberOfRows() == 0 {
		return nil, http.StatusNotFound, fmt.Errorf("Could not find conversation for conversation ID %s", conversationID)
	}

	toReturn := make([]ChatMessage, messages.GetNumberOfRows())
	for i := uint64(0); i < messages.GetNumberOfRows(); i++ {
		message, _ := messages.GetStringValue(i, 0)
		role, _ := messages.GetStringValue(i, 1)
		b64String, _ := messages.GetStringValue(i, 2)
		imgType, _ := messages.GetStringValue(i, 3)
		toReturn[i] = ChatMessage{Message: message, Role: role, B64ImgString: b64String, ImgType: imgType}
	}
	return toReturn, http.StatusOK, nil
}

// The route of interest for the app. to fetch all snippets of conversations:
func (a *App) fetchSnippets(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("No other method but POSTs are allowed here.")
	}

	val := r.Context().Value(cookieName)
	if val == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Something is seriously wrong.  Why is there no session?")
	}
	session := val.(*sessions.Session)
	sessionID, ok := session.Values[sessIdKey]
	if !ok {
		return nil, http.StatusInternalServerError, fmt.Errorf("There's no session ID.  Bruh.")
	}

	const selectStatement = "SELECT id, time_sent " +
		"FROM" +
		" (SELECT session_id, conversation_id AS id, time_sent," +
		" ROW_NUMBER() OVER (PARTITION BY conversation_id ORDER BY time_sent ASC) AS rn" +
		" FROM conversations" +
		" WHERE role = 'model' AND session_id = ?) " +
		"where rn = 1;"
	previews, err := a.DbClient.SelectArray(selectStatement, []any{sessionID})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	previewinfo := make([]ChatPreview, previews.GetNumberOfRows())
	for i := uint64(0); i < previews.GetNumberOfRows(); i++ {
		connID, _ := previews.GetStringValue(i, 0)
		timeSent, _ := previews.GetStringValue(i, 1)
		previewinfo[i] = ChatPreview{ConversationID: connID, Date: timeSent}
	}
	return previewinfo, http.StatusOK, nil
}

// A route to erase all conversations associated with the conversation ID
// of interest:
func (a *App) eraseConversation(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("POST to erase.  Capiche?")
	}

	val := r.Context().Value(cookieName)
	if val == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Something is seriously wrong.  Why is there no session?")
	}
	session := val.(*sessions.Session)
	sessionID, ok := session.Values[sessIdKey]
	if !ok {
		return nil, http.StatusInternalServerError, fmt.Errorf("There's no session ID.  Bruh.")
	}

	var jsonPayload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&jsonPayload); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("What kind of wack-ass JSON did you send over, man?")
	}
	conversationID, ok := jsonPayload["conversation_id"]
	if !ok {
		return nil, http.StatusBadRequest, fmt.Errorf("No `conversationID` field lah, bodoh!")
	}

	const deleteStatement = "DELETE FROM conversations WHERE conversation_id = ? AND session_id = ?;"
	if err := a.DbClient.ExecuteArray(deleteStatement, []any{conversationID, sessionID}); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Can't delete this weird-ass ID: %s", conversationID)
	}
	return SimpleResponse{Message: "Done erasing the conversation!"}, http.StatusOK, nil
}

// Given a new, incoming POST request to start a new conversation, send over a new ChatPreview struct in the
// data field of the JSON response:
func (a *App) createNewConversation(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("Post to make a new conversation.")
	}

	val := r.Context().Value(cookieName)
	if val == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Why got no session one?")
	}
	session := val.(*sessions.Session)
	if _, ok := session.Values[sessIdKey]; !ok {
		return nil, http.StatusForbidden, fmt.Errorf("Gotta log in first to use this, man.")
	}

	newID, err := utils.GenerateID(idLength)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	toReturn := ChatPreview{
		Date:           time.Now().Format("2006-01-02 15:03:04"),
		ConversationID: newID,
	}
	return toReturn, http.StatusOK, nil
}
