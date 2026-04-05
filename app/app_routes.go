// A file to store app-related routes - for the application's key functionalities that don't
// really concern itself with the database or the Gemini server.

package app

import "net/http"

// --- Constants ---

const (
	manifestPath = "./manifest.json"
	swPath       = "./sw.js"
	faviconPath  = "./static/img/logo.png"
)

// --- End ---

// The handler function of interest for serving our manifest.json file - so that we can turn
// our application into a PWA (I still hate this idea, but whatever - Kevin ain't the one
// in charge of this project, plus he's still getting paid anyways)
func (a *App) serveManifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/manifest+json")
	http.ServeFile(w, r, manifestPath)
}

// The handler of interest that'll handle our service worker.  It's currently just sitting in our
// project's root directory:
func (a *App) serveServiceWorker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeFile(w, r, swPath)
}

// The handler of interest that's going to be responsible for serving
// our favicon.ico file.  I'm just gonna re-use Dr. Natalie's icons
// for this - the one from her presentation slides:
func (a *App) faviconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	http.ServeFile(w, r, faviconPath)
}
