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

	router.HandleFunc("/admin/profile", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/profile.html")
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

		currentUser, err := GetCurrentUser(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		student, err := readStudent(currentUser, googleID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("templates/admin/edit_student.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, student)
	}).Methods("GET")

	router.HandleFunc("/admin/student/{googleID}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]

		currentUser, err := GetCurrentUser(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		switch req.Method {
		case http.MethodGet:
			student, err := readStudent(currentUser, googleID)
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
			if err := updateStudent(currentUser, googleID, updates); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT")

	router.HandleFunc("/admin/search_parent", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/search_parent.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/admin/api/search_parent", func(res http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		parents, err := searchParents(name)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(parents)
	}).Methods("GET")

	router.HandleFunc("/admin/parent/{googleID}/edit", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]

		currentUser, err := GetCurrentUser(req)
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parent, err := readParent(currentUser, googleID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("templates/admin/edit_parent.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, parent)
	}).Methods("GET")

	router.HandleFunc("/admin/parent/{googleID}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]

		currentUser, err := GetCurrentUser(req)
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		switch req.Method {
		case http.MethodGet:
			parent, err := readParent(currentUser, googleID)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(parent)
		case http.MethodPut:
			var updates map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			if err := updateParent(currentUser, googleID, updates); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT")

	router.HandleFunc("/admin/search_instructor", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/search_instructor.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/admin/api/search_instructor", func(res http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		instructors, err := searchInstructors(name)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(instructors)
	}).Methods("GET")

	router.HandleFunc("/admin/instructor/{googleID}/edit", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]

		currentUser, err := GetCurrentUser(req)
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		instructor, err := readInstructor(currentUser, googleID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("templates/admin/edit_instructor.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, instructor)
	}).Methods("GET")

	router.HandleFunc("/admin/instructor/{googleID}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]

		currentUser, err := GetCurrentUser(req)
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		switch req.Method {
		case http.MethodGet:
			instructor, err := readInstructor(currentUser, googleID)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(instructor)
		case http.MethodPut:
			var updates map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			if err := updateInstructor(currentUser, googleID, updates); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT")

	router.HandleFunc("/admin/api/profile", func(res http.ResponseWriter, req *http.Request) {
		currentUser, err := GetCurrentUser(req)
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}
		switch req.Method {
		case http.MethodGet:
			admin, err := readAdmin(currentUser, currentUser.GoogleID)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(admin)
		case http.MethodPut:
			var updates map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			if err := updateAdmin(currentUser, currentUser.GoogleID, updates); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT")
}
