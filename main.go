package main

import (
	"log"
	"net/http"

	"github.com/gorilla/pat"
)

func init() {
	// validateJSON("edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json")
	// jsonFilePath := "$HOME/secrets/edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json"

	// err := validateJSON(jsonFilePath)
	// if err != nil {
	// 	log.Fatalf("JSON validation failed: %v", err)
	// }

	// log.Println("JSON is valid.")
	database()
}

func main() {
	// http.HandleFunc("/1", serverhome)
	// http.HandleFunc("/2", setCookieHandler)
	// http.ListenAndServe(":8080", handler())
	// http.ListenAndServeTLS("192.168.1.129:8080", "server.crt", "server.key", handler())

	// err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", handler())

	// Load configuration
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Create a new router
	router := pat.New()

	// Set up authentication routes
	AuthHandler(router, config)

	// Start the HTTPS server with the provided certificate and key files
	log.Println("listening on localhost:8080")
	err = http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", router)
	if err != nil {
		log.Fatal(err)
	}
}
