package main

import (
	"fmt"
	"testing"

	"net/http"

	"net/http/httptest"
)

func ahandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello, Worldn"))

}

func TestRouter(t *testing.T) {

	req := httptest.NewRequest("GET", "http://192.168.1.129:8080/3", nil)

	w := httptest.NewRecorder()

	ahandler(w, req)

	resp := w.Result()

	fmt.Println(resp.StatusCode)

}
