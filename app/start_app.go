package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	sqlitecloud "github.com/sqlitecloud/sqlitecloud-go"
	"google.golang.org/genai"
)

// The App object that contains important objects needed
// for the bot's proper functioning.  Also what InitializeApp() will
// return.
type App struct {
	GeminiClient  *genai.Client
	DbClient      *sqlitecloud.SQCloud
	AesKey        string
	ContentConfig *genai.GenerateContentConfig
	Cookie        *sessions.CookieStore
}

// Constants necessary for the proper functioning of the bot:
const (
	promptPath = "./resources/base_bot_prompt.txt"
)

// The function of interest that'll allow us to set up our application's struct -
// it mainly ensures that our Gemini client and dbConnection strings are in place.
func InitializeApp(ctx context.Context) (*App, error) {
	app := &App{}
	geminiKey := os.Getenv("GEMINI_KEY")
	dbConnString := os.Getenv("SQLITECLOUD_CONNECTION_KEY")
	aesKey := os.Getenv("AES_KEY")
	appSecret := os.Getenv("APP_SECRET")
	if geminiKey == "" {
		return nil, fmt.Errorf("null value for `geminiKey`.")
	}
	if dbConnString == "" {
		return nil, fmt.Errorf("null value for `dbConnString`.")
	}
	if aesKey == "" {
		return nil, fmt.Errorf("null value for `aesKey`.")
	}
	if appSecret == "" {
		return nil, fmt.Errorf("null value for `appSecret`.")
	}

	// Create what needs to be created here:
	geminiClient, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: geminiKey})
	if err != nil {
		return nil, err
	}
	dbClient, err := sqlitecloud.Connect(dbConnString)
	if err != nil {
		dbClient.Close()
		return nil, err
	}
	app.AesKey = aesKey
	contentConfig, err := app.createContentConfig(promptPath)
	if err != nil {
		return nil, err
	}
	cookieStore := sessions.NewCookieStore([]byte(appSecret))

	// Assign values to the app object here:
	app.GeminiClient = geminiClient
	app.DbClient = dbClient
	app.ContentConfig = contentConfig
	app.Cookie = cookieStore

	return app, nil
}

// --- Struct-related methods ---

// A struct-based method to register the App struct's routes -
// both internal and otherwise:
func (a *App) RegisterRoutes() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/manifest.json", http.HandlerFunc(a.serveManifest))
	http.Handle("/sw.js", http.HandlerFunc(a.serveServiceWorker))

	// Webpage-serving routes:
	http.Handle("/", http.HandlerFunc(a.indexPageHandler))
	http.Handle("/signup", http.HandlerFunc(a.creationPageHandler))
	http.Handle("/login", http.HandlerFunc(a.loginPageHandler))
	http.Handle("/chat", a.sessionPropagation(http.HandlerFunc(a.chatPageHandler)))

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
