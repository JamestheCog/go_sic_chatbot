package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jamesthecog/go_chatbot/app"
	"github.com/joho/godotenv"
)

const startUpDuration = 10 // Wait for this amount of seconds for the application to start up.

func main() {
	_ = godotenv.Load()

	timeout, cancel := context.WithTimeout(context.Background(), time.Duration(startUpDuration)*time.Second)
	defer cancel()

	app, err := app.InitializeApp(timeout)
	if err != nil {
		log.Fatalf("Unable to start the application after %d seconds for the following reason: %v\n",
			startUpDuration, err)
	}
	app.RegisterRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting the application now at Singaporean local time: %s", time.Now().Format(time.RFC1123))
	log.Printf("Starting the entire thing on http://localhost:%s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Printf("Couldn't start the application for some reason: %v", err)
	}
}
