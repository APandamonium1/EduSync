package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func StudentHandler(router *mux.Router) {
	router.HandleFunc("/student", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/student/index.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/student/get-folder-id", func(res http.ResponseWriter, req *http.Request) {
		folderID, err := GetStudentFolder(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{"folder_id": folderID}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("GET")
}
