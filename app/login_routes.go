// A file made to contain route handlers that are responsible for session authentication - or
// in plain language, logging in and out:

package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

// -- Constants --

const formMemory = 10 << 20 // Allocate 10 MB for form processing - overkill, but should be sufficient.

// --- End ---

// For creating new users:
func (a *App) newUserHandler(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("Only POSTs are allowed to this route.")
	}

	if err := r.ParseMultipartForm(formMemory); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	curLoginInfo, err := a.fetchLoginInfo(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if curLoginInfo.IsSuccessful && curLoginInfo.SessionID != "" {
		return nil, http.StatusConflict, fmt.Errorf("The username already exists!")
	}

	if err := a.createUser(r.FormValue("username"), r.FormValue("password")); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return SimpleResponse{Message: "Your profile has been created - redirecting you to the login page now!"},
		http.StatusOK, nil
}

// The login handler of interest:
func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("Only POSTs are allowed in this route.")
	}

	if err := r.ParseMultipartForm(formMemory); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	loginAttempt, err := a.fetchLoginInfo(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if !loginAttempt.IsSuccessful && loginAttempt.SessionID == "" {
		return nil, http.StatusNotFound, fmt.Errorf("The entered username and / or password is incorrect.")
	}

	val := r.Context().Value(cookieName)
	if val == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Why is there no session?  WHY!?")
	}
	session := val.(*sessions.Session)
	session.Values[sessIdKey] = loginAttempt.SessionID
	session.Save(r, w)

	return SimpleResponse{Message: "You've been successfully logged in!  Redirecting you now..."}, http.StatusOK, nil
}

// The route of interest for logging out:
func (a *App) logoutHandler(w http.ResponseWriter, r *http.Request) (any, int, error) {
	if r.Method != http.MethodPost {
		return nil, http.StatusMethodNotAllowed, fmt.Errorf("Only POSTs are allowed to logout.")
	}

	val := r.Context().Value(cookieName)
	if val == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("There's no session object; something is seriously wrong.")
	}
	session := val.(*sessions.Session)

	if _, ok := session.Values[sessIdKey]; !ok {
		return nil, http.StatusInternalServerError, fmt.Errorf("There's no active session to log out from.")
	}
	session.Values[sessIdKey] = ""
	session.Save(r, w)

	return SimpleResponse{Message: "Logged out went well!"}, http.StatusOK, nil
}
