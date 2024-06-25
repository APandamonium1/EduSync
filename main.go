package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	database()
}

func index(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "home.html", nil)

	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	// // Load configuration
	// config, err := LoadConfig("config.json")
	// if err != nil {
	// 	log.Fatalf("could not load config: %v", err)
	// }

	router := mux.NewRouter()
	fs := http.FileServer(http.Dir("assets"))
	router.PathPrefix("/assets").Handler(http.StripPrefix("/assets", fs))
	//http.Handle("/resources/", http.StripPrefix("/resources", fs))
	router.HandleFunc("/", index)
	log.Println("listening on localhost:8080")
	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", router)
	if err != nil {
		log.Fatal(err)
	}

	// // Set up authentication routes
	// AuthHandler(router, config)
}
