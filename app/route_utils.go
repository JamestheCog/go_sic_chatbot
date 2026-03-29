package app

import (
	"encoding/json"
	"log"
	"net/http"
)

// === END ===

func (a *App) sendJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Something bad happened: %v", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("Something went horribly wrong - check the console!"))
	}
}
