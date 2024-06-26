package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func RoleHandler(router *mux.Router) {
	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, false)
	}).Methods("GET")

	router.HandleFunc("/admin", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/index.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/student", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/student/index.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/instructor", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/instructor/index.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")
}
