package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	initializeFirebase()
	err := initializeFirebase()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
}

func main() {
	// Load configuration
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	router := mux.NewRouter()

	// Serving static files
	fs := http.FileServer(http.Dir("assets"))
	router.PathPrefix("/assets").Handler(http.StripPrefix("/assets", fs))

	// Set up authentication routes
	AuthHandler(router, config)
	MainHandler(router)

	log.Println("listening on localhost:8080")
	err = http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", router)
	if err != nil {
		log.Fatal(err)
	}
}
