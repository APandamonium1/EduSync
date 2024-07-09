package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func ParentHandler(router *mux.Router) {
	router.HandleFunc("/parent", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/parent/index.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")
}
