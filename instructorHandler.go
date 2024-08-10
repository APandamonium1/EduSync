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

	router.HandleFunc("/instructor/learning-materials/classes", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/instructor/lm-classes.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/instructor/classes/get-classes", func(res http.ResponseWriter, req *http.Request) {
		GetInstructorClasses(res, req)
	}).Methods("GET")

	router.HandleFunc("/instructor/learning-materials/upload", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/instructor/lm-upload.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/instructor/profile", func(res http.ResponseWriter, req *http.Request) {
		instructor, err := GetCurrentInstructor(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Render the profile page
		t, err := template.ParseFiles("templates/instructor/profile.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, instructor)
	}).Methods("GET")
}
