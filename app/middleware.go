package app

import (
	"context"
	"net/http"
)

// A custom type we'll be using to define handler functions:
type handlerFunc func(w http.ResponseWriter, r *http.Request) (any, int, error)

// -- Structs and constants for what we'll be sending over to the front-end --
const successMsg = "All went well, my brotha / sister!"

type APIResponse struct {
	Status Status `json:"status"`
	Data   any    `json:"data,omitempty"`
}
type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// -- End --

// Middleware:

// Wraps all our internal routes - so that it sends JSON in our desired format
// regardless of success or failure.
func (a *App) handleJsonPayload(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, code, err := h(w, r)
		response := APIResponse{Status: Status{Code: code, Message: successMsg}}

		if err != nil && payload == nil {
			response.Status.Message = err.Error()
			if code == 0 {
				response.Status.Code = http.StatusInternalServerError
			}
		} else {
			response.Data = payload
		}
		a.sendJSON(w, code, response)
	}
}

// Passes contexts from one request to another:
func (a *App) sessionPropagation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := a.Cookie.Get(r, cookieName)
		ctx := context.WithValue(r.Context(), cookieName, session)
		next.ServeHTTP(w, r.WithContext(ctx))
		session.Save(r, w)
	})
}
