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
	router := mux.NewRouter()

	// Serving static files
	fs := http.FileServer(http.Dir("assets"))
	router.PathPrefix("/assets").Handler(http.StripPrefix("/assets", fs))

	// Set up authentication routes
	AuthHandler(router)
	MainHandler(router)
	AdminHandler(router)

	log.Println("listening on localhost:8080")
	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", router)
	if err != nil {
		log.Fatal(err)
	}
}
