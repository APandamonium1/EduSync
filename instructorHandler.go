package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func InstructorHandler(router *mux.Router) {
	router.HandleFunc("/instructor", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/instructor/index.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/instructor/classes", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/instructor/classes.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/instructor/classes/get-classes", func(res http.ResponseWriter, req *http.Request) {
		GetInstructorClasses(res, req)
	}).Methods("GET")
}
