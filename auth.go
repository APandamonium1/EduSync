package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var firebaseAuth *auth.Client
var store = sessions.NewCookieStore([]byte("super-secret-key"))

func AuthHandler(router *mux.Router) {
	var err error
	firebaseAuth, err = firebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	router.HandleFunc("/login", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/callback", func(res http.ResponseWriter, req *http.Request) {
		idToken := req.URL.Query().Get("token")
		token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Here you can retrieve additional user information from Firebase
		userRecord, err := firebaseAuth.GetUser(context.Background(), token.UID)
		if err != nil {
			http.Error(res, "Unable to get user info", http.StatusInternalServerError)
			return
		}

		user := User{
			GoogleID:      userRecord.UID,
			Name:          userRecord.DisplayName,
			Email:         userRecord.Email,
			ContactNumber: userRecord.PhoneNumber,
			Role:          "User", // or whatever logic you have to determine the role
		}

		session, err := store.Get(req, "auth-session")
		if err != nil {
			http.Error(res, "Unable to get session", http.StatusInternalServerError)
			return
		}

		userData, err := json.Marshal(user)
		if err != nil {
			http.Error(res, "Unable to marshal user data", http.StatusInternalServerError)
			return
		}

		session.Values["user"] = userData
		err = session.Save(req, res)
		if err != nil {
			http.Error(res, "Unable to save session", http.StatusInternalServerError)
			return
		}

		http.Redirect(res, req, "/protected", http.StatusFound)
	}).Methods("GET")

	router.HandleFunc("/protected", func(res http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, "auth-session")
		if err != nil {
			http.Error(res, "Unable to get session", http.StatusInternalServerError)
			return
		}

		userID, ok := session.Values["userID"].(string)
		if !ok || userID == "" {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Use the userID to fetch user details if needed
		fmt.Fprintf(res, "Protected content for user: %s", userID)
	}).Methods("GET")
}
