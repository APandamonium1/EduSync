package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ahandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello, World"))

}

func TestRouter(t *testing.T) {

	// req := httptest.NewRequest("GET", "http://192.168.1.129:8080/3", nil)
	req := httptest.NewRequest("GET", "http://127.0.0.1:8080/1", nil)
	w := httptest.NewRecorder()

	ahandler(w, req)

	resp := w.Result()

	fmt.Println(resp.StatusCode)

}
