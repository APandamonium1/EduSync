package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	// database()
	initializeFirebase()
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

	// Start the HTTPS server with the provided certificate and key files
	log.Println("listening on localhost:8080")
	err = http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", router)
	if err != nil {
		log.Fatal(err)
	}
}
