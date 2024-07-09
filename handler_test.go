package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestRoleHandler(t *testing.T) {
	router := mux.NewRouter()
	MainHandler(router)
	InstructorHandler(router)
	StudentHandler(router)
	ParentHandler(router)
	AdminHandler(router)

	tests := []struct {
		route          string
		expectedStatus int
	}{
		{"/", http.StatusOK},
		{"/admin", http.StatusOK},
		{"/student", http.StatusOK},
		{"/instructor", http.StatusOK},
		{"/parent", http.StatusOK},
	}

	for _, tt := range tests {
		req, err := http.NewRequest("GET", tt.route, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		if res.Code != tt.expectedStatus {
			t.Errorf("expected status %v; got %v", tt.expectedStatus, res.Code)
		}
	}
}
