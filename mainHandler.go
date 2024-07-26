/*
Gofmt formats Go programs.
It uses tabs for indentation and blanks for alignment.
Alignment assumes that an editor is using a fixed-width font.
package main
*/
package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func MainHandler(router *mux.Router) {
	/*
		Gofmt formats Go programs.
		It uses tabs for indentation and blanks for alignment.
		Alignment assumes that an editor is using a fixed-width font.
		package main
		MAIN
	*/
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
