package main

import (
	"encoding/json"
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

	store = sessions.NewCookieStore(
		[]byte(config.AuthKey),
		[]byte(config.EncryptKey),
	)
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   true, // This should be true if your application is served over HTTPS
	}

	gothic.Store = store
	goth.UseProviders(google.New(config.GoogleClientID, config.GoogleClientSecret, "https://localhost:8080/auth/google/callback", "email", "profile"))

	router.HandleFunc("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}

		userObj, userRole, err := getUserRole(user.Email)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		log.Println("User role:", userRole)

		// Only store the user object into the session if userRole is not an empty string
		if userRole != "" {
			// // Create a User object with the user role
			// currentUser := User{
			// 	GoogleID:      user.UserID,
			// 	Name:          user.Name,
			// 	Email:         user.Email,
			// 	ContactNumber: userObj.ContactNumber, // Use contact number from the retrieved user object
			// 	Role:          userObj.Role,
			// 	CreatedAt:     userObj.CreatedAt,
			// 	UpdatedAt:     userObj.UpdatedAt,
			// }

			// // Serialize the user object to JSON
			// userData, err := json.Marshal(currentUser)
			// if err != nil {
			// 	http.Error(res, err.Error(), http.StatusInternalServerError)
			// 	return
			// }

			// // Get the session and store the user data
			// session, err := store.Get(req, "auth-session")
			// if err != nil {
			// 	http.Error(res, err.Error(), http.StatusInternalServerError)
			// 	return
			// }
			// session.Values["user"] = userData
			// err = session.Save(req, res)
			// if err != nil {
			// 	http.Error(res, err.Error(), http.StatusInternalServerError)
			// 	return
			// }

			SetCurrentUser(res, req, userObj)

			// Redirect based on user role
			if userRole == "Admin" {
				AdminHandler(router)
				http.Redirect(res, req, "/admin", http.StatusFound)
			} else if userRole == "Instructor" {
				InstructorHandler(router)
				http.Redirect(res, req, "/instructor", http.StatusFound)
			} else if userRole == "Student" {
				StudentHandler(router)
				http.Redirect(res, req, "/student", http.StatusFound)
			} else if userRole == "Parent" {
				ParentHandler(router)
				http.Redirect(res, req, "/parent", http.StatusFound)
			}
		} else {
			http.Redirect(res, req, "/unregistered", http.StatusFound)
		}
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

func SetCurrentUser(res http.ResponseWriter, req *http.Request, user User) error {
	session, err := store.Get(req, "auth-session")
	if err != nil {
		return fmt.Errorf("error retrieving session: %v", err)
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error marshalling user data: %v", err)
	}

	session.Values["user"] = userData
	err = session.Save(req, res)
	if err != nil {
		return fmt.Errorf("error saving session: %v", err)
	}

	return nil
}
