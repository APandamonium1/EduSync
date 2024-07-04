package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func AdminHandler(router *mux.Router) {
	router.HandleFunc("/admin", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/index.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/admin/search_student", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/search_student.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/admin/api/search_student", func(res http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		class := req.URL.Query().Get("class")
		students, err := searchStudents(name, class)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(students)
	}).Methods("GET")

	router.HandleFunc("/admin/student/{googleID}/edit", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]
		t, err := template.ParseFiles("templates/admin/edit_student.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		student, err := readStudent(googleID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, student)
	}).Methods("GET")

	router.HandleFunc("/admin/student/{googleID}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]

		switch req.Method {
		case http.MethodGet:
			student, err := readStudent(googleID)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(student)
		case http.MethodPut:
			var updates map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			if err := updateStudent(googleID, updates); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT")
}
