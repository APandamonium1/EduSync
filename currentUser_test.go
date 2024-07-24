package main

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentUser_SessionFound(t *testing.T) {
	// Arrange
	req := httptest.NewRequest("GET", "/test", nil)
	session, _ := store.New(req, "auth-session")
	session.Values["user"] = []byte(`{"googleID": "testUser", "name": "Test User", "email": "test@example.com", "role": "admin"}`)
	session.Save(req, nil)

	// Act
	user, err := GetCurrentUser(req)

	// Assert
	assert.Nil(t, err, "Expected no error when session is found")
	assert.Equal(t, User{GoogleID: "testUser", Name: "Test User", Email: "test@example.com", Role: "admin"}, user, "Expected correct user data")
}

func TestGetCurrentUser_SessionNotFound(t *testing.T) {
	// Arrange
	req := httptest.NewRequest("GET", "/test", nil)

	// Act
	user, err := GetCurrentUser(req)

	// Assert
	assert.Error(t, err, "Expected error when session is not found")
	assert.Equal(t, User{}, user, "Expected empty user when session is not found")
}

func TestGetCurrentUser_InvalidUserData(t *testing.T) {
	// Arrange
	req := httptest.NewRequest("GET", "/test", nil)
	session, _ := store.New(req, "auth-session")
	session.Values["user"] = []byte(`{"googleID": "testUser", "name": "Test User", "email": "test@example.com", "role": "invalid"}`)
	session.Save(req, nil)

	// Act
	_, err := GetCurrentUser(req)

	// Assert
	assert.Error(t, err, "Expected error when invalid user data is found")
}

func TestGetCurrentUser_EmptyUserData(t *testing.T) {
	// Arrange
	req := httptest.NewRequest("GET", "/test", nil)
	session, _ := store.New(req, "auth-session")
	session.Values["user"] = nil
	session.Save(req, nil)

	// Act
	user, err := GetCurrentUser(req)

	// Assert
	assert.Error(t, err, "Expected error when user data is empty")
	assert.Equal(t, User{}, user, "Expected empty user when user data is empty")
}

func TestGetCurrentUser_NoSession(t *testing.T) {
	// Arrange
	req := httptest.NewRequest("GET", "/test", nil)

	// Act
	user, err := GetCurrentUser(req)

	// Assert
	assert.Error(t, err, "Expected error when no session is found")
	assert.Equal(t, User{}, user, "Expected empty user when no session is found")
}
