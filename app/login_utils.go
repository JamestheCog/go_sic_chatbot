// A collection of helper functions, structs, and constants (maybe) for the logging in
// functionality of our application.

package app

import (
	"github.com/jamesthecog/go_chatbot/utils"
)

// --- Structs and constants ---

// The struct of interest to be sent over to the front-end if
type LoginInfo struct {
	IsSuccessful bool
	Username     string
	SessionID    string
}

// Given a username and a password, see if a unique ID has been generated for the said
// user.  If not, then return a false LoginInfo item:
func (a *App) fetchLoginInfo(username, password string) (LoginInfo, error) {
	const selectStatement = "SELECT session_id, username " +
		"FROM user_info " +
		"wHERE username = ? AND password = ?;"
	toReturn := LoginInfo{}
	result, err := a.DbClient.SelectArray(selectStatement, []any{username, password})
	if err != nil {
		return toReturn, err
	}

	sessionID, _ := result.GetStringValue(0, 0)
	fetchedUsername, _ := result.GetStringValue(0, 1)
	if sessionID != "" && fetchedUsername != "" {
		toReturn.IsSuccessful = true
		toReturn.SessionID = sessionID
		toReturn.Username = fetchedUsername
	}
	return toReturn, nil
}

// Given a username and password, create a new, unique ID for this user and their
// (probably lame-as-anything password):
func (a *App) createUser(username, password string) error {
	const insertStatement = "INSERT INTO user_info VALUES (?, ?, ?);"
	sessionID, err := utils.GenerateID(idLength)
	if err != nil {
		return err
	}

	if err := a.DbClient.ExecuteArray(insertStatement, []any{username, password, sessionID}); err != nil {
		return err
	}
	return nil
}
