package main

import (
	"fmt"
	"html/template"
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
