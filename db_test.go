package main

import (
	"testing"
)

func TestConnect(t *testing.T) {
	// Create a new instance of FireDB.
	db := FireDB{}

	// Call the Connect function.
	err := db.Connect()

	// Check if there was an error during connection.
	if err != nil {
		t.Errorf("Failed to connect to Firebase: %v", err)
	}

	// Check if the Firebase Realtime Database client is initialized.
	if db.Client == nil {
		t.Error("Firebase Realtime Database client is not initialized")
	}
}
