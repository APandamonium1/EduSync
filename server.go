package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("hello.html"))

func init() {
	fmt.Println("Hello init!")
}

func handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/1", serverhome)
	r.HandleFunc("/2", setCookieHandler)
	return r
}

func serverhome(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Hello World!"))
	fmt.Println(r)
	templates.ExecuteTemplate(w, "hello.html", nil)
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize a new cookie containing the string "Hello world!" and some
	// non-default attributes.
	cookie := http.Cookie{
		Name:  "exampleCookie",
		Value: "Hello world!",
		// Path:     "/",
		// MaxAge:   3600,
		// HttpOnly: true,
		// Secure:   true,
		// SameSite: http.SameSiteLaxMode,
	}

	// Use the http.SetCookie() function to send the cookie to the client.
	// Behind the scenes this adds a `Set-Cookie` header to the response
	// containing the necessary cookie data.
	http.SetCookie(w, &cookie)

	// Write a HTTP response as normal.
	w.Write([]byte("cookie set!"))
}
