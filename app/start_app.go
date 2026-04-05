package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/jamesthecog/go_chatbot/utils"
	sqlitecloud "github.com/sqlitecloud/sqlitecloud-go"
)

// The App object that contains important objects needed
// for the bot's proper functioning.  Also what InitializeApp() will
// return.
type App struct {
	DbClient           *sqlitecloud.SQCloud
	AesKey             string
	GeminiKeys         []string
	Cookie             *sessions.CookieStore
	ConversationPrompt string
}

// Constants necessary for the proper functioning of the bot:
const (
	promptPath    = "./resources/base_bot_prompt.txt"
	geminiKeyPath = "./resources/gemini_keys.txt"
)

// The function of interest that'll allow us to set up our application's struct -
// it mainly ensures that our Gemini client and dbConnection strings are in place.
func InitializeApp(ctx context.Context) (*App, error) {
	app := &App{}
	dbConnString := os.Getenv("SQLITECLOUD_CONNECTION_STRING")
	aesKey := os.Getenv("AES_KEY")
	appSecret := os.Getenv("APP_SECRET")
	if dbConnString == "" {
		return nil, fmt.Errorf("null value for `SQLITECLOUD_CONNECTION_STRING`.")
	}
	if aesKey == "" {
		return nil, fmt.Errorf("null value for `AES_KEY`.")
	}
	if appSecret == "" {
		return nil, fmt.Errorf("null value for `APP_SECRET`.")
	}

	app.AesKey = aesKey
	dbClient, err := sqlitecloud.Connect(dbConnString)
	if err != nil {
		dbClient.Close()
		return nil, err
	}
	conversationPrompt, err := utils.LoadFile(promptPath, aesKey)
	if err != nil {
		return nil, err
	}
	geminiKeys, err := fetchGeminiKeys(geminiKeyPath, aesKey)
	if err != nil {
		return nil, err
	}
	cookieStore := sessions.NewCookieStore([]byte(appSecret))

	// Assign values to the app object here:
	app.DbClient = dbClient
	app.ConversationPrompt = conversationPrompt
	app.Cookie = cookieStore
	app.GeminiKeys = geminiKeys
	return app, nil
}

// Given a file path and a hex key, return the file's decrypted contents
// as a string:

// Given a the file path to the Gemini keys, return a string slice containing its
// contents:
func fetchGeminiKeys(filePath, hexKey string) ([]string, error) {
	decryptedContents, err := utils.LoadFile(filePath, hexKey)
	if err != nil {
		return nil, err
	}

	tokens := strings.Split(strings.TrimSpace(decryptedContents), "\n")
	toReturn := make([]string, len(tokens))
	for i, v := range tokens {
		toReturn[i] = strings.TrimSpace(v)
	}
	return toReturn, nil
}

// --- Struct-related methods ---

// A struct-based method to register the App struct's routes -
// both internal and otherwise:
func (a *App) RegisterRoutes() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/favicon.ico", http.HandlerFunc(a.faviconHandler))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/manifest.json", http.HandlerFunc(a.serveManifest))
	http.Handle("/sw.js", http.HandlerFunc(a.serveServiceWorker))

	// Webpage-serving routes:
	http.Handle("/", http.HandlerFunc(a.indexPageHandler))
	http.Handle("/signup", http.HandlerFunc(a.creationPageHandler))
	http.Handle("/login", http.HandlerFunc(a.loginPageHandler))
	http.Handle("/chat", a.sessionPropagation(http.HandlerFunc(a.chatPageHandler)))
	http.Handle("/offline", http.HandlerFunc(a.offlinePageHandler))

	// Internal routes:
	//
	// -- For the application's core functionalities --
	http.Handle("/internal/chat", a.sessionPropagation(a.handleJsonPayload(a.messageHandler)))
	http.Handle("/internal/new_chat", a.sessionPropagation(a.handleJsonPayload(a.createNewConversation)))
	http.Handle("/internal/fetch_messages", a.sessionPropagation(a.handleJsonPayload(a.fetchMessages)))
	http.Handle("/internal/fetch_snippets", a.sessionPropagation(a.handleJsonPayload(a.fetchSnippets)))
	http.Handle("/internal/delete_conversation", a.sessionPropagation(a.handleJsonPayload((a.eraseConversation))))

	// -- For new user / login / logout handlers --
	http.Handle("/internal/new_user", a.handleJsonPayload(a.newUserHandler))
	http.Handle("/internal/login", a.sessionPropagation(a.handleJsonPayload(a.loginHandler)))
	http.Handle("/internal/logout", a.sessionPropagation(a.handleJsonPayload(a.logoutHandler)))

	// -- Miscellaneous routes - mainly to prevent cold starts --
	http.HandleFunc("/internal/wake_app", a.handleJsonPayload(a.wakeApplication))
	http.HandleFunc("/internal/wake_db", a.handleJsonPayload(a.wakeDB))
}
