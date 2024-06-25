package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// AuthHandler sets up the authentication routes on the provided router
func AuthHandler(router *pat.Router, config *Config) {
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

	router.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}

		// Extract user details from the goth.User
		googleID := user.UserID
		name := user.Name
		email := user.Email

		// Example: Create or update a student with the extracted details
		student := NewStudent(googleID, name, 18, 10.0, email, "91234567", "TE", "Mr. Smith", "Mrs. Doe")
		err = createStudent(student.GoogleID, student)
		if err != nil {
			log.Println("Error creating student:", err)
		} else {
			log.Println("Student created/updated successfully!")
		}

		t, _ := template.ParseFiles("templates/success.html")
		t.Execute(res, user)
	})

	router.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})
}
