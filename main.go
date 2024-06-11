package main

import (
	"log"
	"net/http"
)

func init() {
	database()
}

func main() {
	// http.HandleFunc("/1", serverhome)
	// http.HandleFunc("/2", setCookieHandler)
	// http.ListenAndServe(":8080", handler())
	// http.ListenAndServeTLS("192.168.1.129:8080", "server.crt", "server.key", handler())

	// err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", handler())
	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", handler())
	if err != nil {
		log.Fatal(err)
	}
}
