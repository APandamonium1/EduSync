package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func AuthHandler(router *mux.Router, config *Config) {
	maxAge := 86400 * 30 // 30 days
	isProd := true       // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(config.SessionSecret))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(config.GoogleClientID, config.GoogleClientSecret, "https://localhost:8080/auth/google/callback", "email", "profile"),
	)

	router.HandleFunc("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}

		// Example role assignment logic
		var role string
		switch user.Email {
		case "admin@example.com":
			role = "Admin"
		case "instructor@example.com":
			role = "Instructor"
		case "parent@example.com":
			role = "Parent"
		default:
			role = "Student"
		}

		// Create or update the user in Firebase with the assigned role
		// student := NewStudent(user.UserID, user.Name, user.Email, "91234567", "TE", "Mr. Smith", "Jane Doe", role, 10, 10.0)
		student := NewStudent(user.UserID, user.Name, user.Email, "91234567", "TE", "Jane Doe", role, 10, 10.0)
		err = createStudent(student.User, student)
		if err != nil {
			log.Println("Error creating student:", err)
		} else {
			log.Println("Student created/updated successfully!")
		}

		t, err := template.ParseFiles("templates/success.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, user)
	}).Methods("GET")

	router.HandleFunc("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	}).Methods("GET")

	router.HandleFunc("/login", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, false)
	}).Methods("GET")
}
