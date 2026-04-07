// A package I came up with to store helper functions that've to do with the fetching and storing of data
// from the SQLitecloud database that this project relies on.
//
// Formatting utilities are also present here - though they're private.

package app

import (
	"encoding/base64"
	"fmt"
	"time"

	"google.golang.org/genai"
)

// === Structs and Constants ==

const (
	limitAmount = 50      // Just for good measure - each conversation probably won't hit this amount of messages though.
	userRole    = "user"  // The user's role
	modelRole   = "model" // The model's role
)

// A struct to store results fetched from the remote SQLitecloud database:
type DbMessage struct {
	Message   string
	Role      string
	ImgString string
	ImgType   string
}

// --- Public helper functions here ---

// Given a randomly generated conversation ID and user ID, fetch its associated
// messages from the SQLitecloud database.
func (a *App) fetchConversationHistory(connID, sessionID string) ([]*genai.Content, error) {
	const selectStatement = "SELECT message, role, img_b64, img_mimetype " +
		"FROM conversations " +
		"WHERE conversation_id = ? AND session_id = ? " +
		"ORDER BY time_sent LIMIT ?;"
	selectVals := []any{connID, sessionID, limitAmount}
	contents, err := a.DbClient.SelectArray(selectStatement, selectVals)
	if err != nil {
		return nil, err
	}

	formattedMsgs := make([]DbMessage, contents.GetNumberOfRows())
	for i := uint64(0); i < contents.GetNumberOfRows(); i++ {
		msg, _ := contents.GetStringValue(i, 0)
		role, _ := contents.GetStringValue(i, 1)
		rawImg, _ := contents.GetStringValue(i, 2)
		imgMimeType, _ := contents.GetStringValue(i, 3)

		formattedMsgs[i] = DbMessage{Message: msg, Role: role, ImgString: rawImg, ImgType: imgMimeType}
	}

	toReturn, err := formatForGemini(formattedMsgs)
	if err != nil {
		return nil, err
	}
	return toReturn, nil
}

// A helper function that - given a MsgPayload struct or a string, formats
// the data to be uploaded onto the remote SQLitecloud database.
func (a *App) uploadMessage(messageItems any, connID, sessionID string) error {
	var valArray []any
	const uploadStatement = "INSERT INTO conversations VALUES (?, ?, ?, ?, ?, ?, ?);"
	timeToday := time.Now().Format("2006-01-02 15:03:04")

	switch v := messageItems.(type) {
	case MsgPayload:
		imgByte := []byte{}
		if v.Message.ImgB64 != "" {
			imgByte = []byte(v.Message.ImgB64)
		}
		valArray = []any{v.Message.Msg, v.Message.Role, timeToday, v.Message.ImgType, imgByte, connID, sessionID}
	case string:
		valArray = []any{messageItems, modelRole, timeToday, "", []byte{}, connID, sessionID}
	default:
		return fmt.Errorf("uploadMessage only knows how to work with strings and MsgPayload!")
	}

	if err := a.DbClient.ExecuteArray(uploadStatement, valArray); err != nil {
		return err
	}
	return nil
}

// Non-App
func formatForGemini(dbMsgs []DbMessage) ([]*genai.Content, error) {
	toReturn := make([]*genai.Content, len(dbMsgs))

	for i, v := range dbMsgs {
		if v.Role != userRole && v.Role != modelRole {
			return nil, fmt.Errorf("`v.Role` needs to either be of value `userRole` or `modelRole`.")
		}

		contentPart := []*genai.Part{{Text: v.Message}}
		if v.ImgString != "" && v.ImgType != "" {
			imgPart, err := decodeImage(v)
			if err != nil {
				return nil, err
			}
			contentPart = append(contentPart, imgPart)
		}
		content := &genai.Content{
			Role:  v.Role,
			Parts: contentPart,
		}
		toReturn[i] = content
	}
	return toReturn, nil
}

// Given a bsae64 string that's fetched from the database, decode it into a byte slice.
func decodeImage(msg DbMessage) (*genai.Part, error) {
	unencoded, err := base64.StdEncoding.DecodeString(msg.ImgString)
	if err != nil {
		return nil, err
	}

	imagePart := &genai.Part{InlineData: &genai.Blob{MIMEType: msg.ImgType, Data: unencoded}}
	return imagePart, nil
}
