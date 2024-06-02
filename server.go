package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var homePages = map[string]string{
	"":        "index.html",
	"about":   "about.html",
	"contact": "contact.html",
	"login":   "login.html",
}

var adminPages = map[string]string{
	"":      "index.html",
	"about": "about.html",
}

var tmpl *template.Template
var adminTmpl *template.Template
var templateDir = "templates"

func init() {
	var err error
	tmpl, err = template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		fmt.Println("Error parsing home templates: ", err)
		os.Exit(1)
	}

	adminTmpl, err = template.ParseGlob(filepath.Join(templateDir, "admin", "*.html"))
	if err != nil {
		fmt.Println("Error parsing admin templates: ", err)
		os.Exit(1)
	}
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/admin/", adminHandler)
	mux.HandleFunc("/", mainHandler)
	fmt.Println("Server started")
	// http.ListenAndServe(":8080", mux)
	// r.HandleFunc("/1", serverhome)
	// r.HandleFunc("/2", setCookieHandler)
	return mux
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/admin/"):]
	templateFile, ok := adminPages[path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	adminTmpl.ExecuteTemplate(w, templateFile, nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	templateFile, ok := homePages[path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	tmpl.ExecuteTemplate(w, templateFile, nil)
}

// func serverhome(w http.ResponseWriter, r *http.Request) {
// 	// w.Write([]byte("Hello World!"))
// 	fmt.Println(r)
// 	templates.ExecuteTemplate(w, "hello.html", nil)
// }

// func setCookieHandler(w http.ResponseWriter, r *http.Request) {
// 	// Initialize a new cookie containing the string "Hello world!" and some
// 	// non-default attributes.
// 	cookie := http.Cookie{
// 		Name:  "exampleCookie",
// 		Value: "Hello world!",
// 		// Path:     "/",
// 		// MaxAge:   3600,
// 		// HttpOnly: true,
// 		// Secure:   true,
// 		// SameSite: http.SameSiteLaxMode,
// 	}

// 	// Use the http.SetCookie() function to send the cookie to the client.
// 	// Behind the scenes this adds a `Set-Cookie` header to the response
// 	// containing the necessary cookie data.
// 	http.SetCookie(w, &cookie)

// 	// Write a HTTP response as normal.
// 	w.Write([]byte("cookie set!"))
// }
