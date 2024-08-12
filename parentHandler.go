package main

import (
	"encoding/json"
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

	router.HandleFunc("/parent/profile", func(res http.ResponseWriter, req *http.Request) {
		parent, err := GetCurrentParent(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Render the profile page
		t, err := template.ParseFiles("templates/parent/profile.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, parent)
	}).Methods("GET")

	router.HandleFunc("/parent/get-folder-id", func(res http.ResponseWriter, req *http.Request) {
		parent, err := GetCurrentParent(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{"folder_id": parent.FolderID}
		res.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("GET")
}
