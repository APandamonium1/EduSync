// This file contains the main handler for the web application.
package main

//The import section declares packages that this source file depends on.
import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// MainHandler handles the main page of the web application.
func MainHandler(router *mux.Router) {
	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, false)
	}).Methods("GET")

	router.HandleFunc("/unregistered", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/unregistered.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, false)
	}).Methods("GET")
}
