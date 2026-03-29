// A file to store all routes that've to do with the serving of webpages - which in our case,
// only happens to be one for now.

package app

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
)

const (
	sessUsernameKey = "username"
	conIdKey        = "conversation_id"
	sessIdKey       = "session_id"
	cookieName      = "conversation"
	indexRoute      = "./views/index.html"
	loginRoute      = "./views/login.html"
	newUserRoute    = "./views/signup.html"
	chatRoute       = "./views/chat.html"
	idLength        = 10
)

// The index route:
func (a *App) indexPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, indexRoute)
}

// The login page:
func (a *App) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, loginRoute)
}

// The user creation page:
func (a *App) creationPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, newUserRoute)
}

// The main route of interest that's our application's chatting page.  We be baking cookies
// here that Famous Amos would yearn for.
func (a *App) chatPageHandler(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value(cookieName)
	if val == nil {
		http.Error(w, "No session available to fetch from.", http.StatusInternalServerError)
		return
	}
	session := val.(*sessions.Session)

	sessionID, ok := session.Values[sessIdKey]
	if ok {
		if sessionID == nil || sessionID == "" {
			msg := url.QueryEscape("You need to be logged in to chat with the assistant.")
			http.Redirect(w, r, fmt.Sprintf("/login?msg=%s&type=error", msg), http.StatusFound)
			return
		}
		http.ServeFile(w, r, chatRoute)
	} else {
		http.Error(w, "No session ID present; something is very wrong!", http.StatusInternalServerError)
	}
}
