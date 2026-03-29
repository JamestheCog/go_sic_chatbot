// A file to store all miscellaneous routes.  Mostly routes that've to do with
// keeping the database alive and the application from going cold:

package app

import (
	"fmt"
	"net/http"
)

// --- Constants and structs for dying ---

// What we're going to be sending to the middleware to return to
// the front-end (i.e., JavaScript):
type SimpleResponse struct {
	Message string `json:"message"`
}

// --- End ---

// The route for keeping the application awake - to prevent cold starts on
// Render's free tiers.  It'll work regardless of what REST verb is
// sent over to this route:
func (a *App) wakeApplication(w http.ResponseWriter, r *http.Request) (any, int, error) {
	awakener := SimpleResponse{Message: "Chatty botty now awaky, mommy!"}
	return awakener, http.StatusOK, nil
}

// The route that wakes the database up - I'm not sure if a SELECT statement will
// be sufficient, but can't hurt to try:
func (a *App) wakeDB(w http.ResponseWriter, r *http.Request) (any, int, error) {
	const wakeupStatement = "SELECT message, role FROM conversations LIMIT 10;"
	awakener := SimpleResponse{}

	if err := a.DbClient.Execute(wakeupStatement); err != nil {
		awakener.Message = "DB too tired cannot wake up :("
		return awakener, http.StatusInternalServerError, fmt.Errorf("DB no waky waky.")
	}
	awakener.Message = "DB now awake - happy?"
	return awakener, http.StatusOK, nil
}
