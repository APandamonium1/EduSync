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
	// req := httptest.NewRequest("GET", "http://127.0.0.1:8080/1", nil)
	// w := httptest.NewRecorder()

	// ahandler(w, req)

	// resp := w.Result()

	// fmt.Println(resp.StatusCode)
	req, err := http.NewRequest("GET", "/3", nil)
	fmt.Println(req.Body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandler)

	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	// expected := "Hi"
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }

}
