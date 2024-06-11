package main

import (
	"context"
	"fmt"
	"testing"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// func TestConnect(t *testing.T) {
// 	// Create a new instance of FireDB.
// 	db := FireDB{}

// 	// Call the Connect function.
// 	err := db.Connect()

// 	// Check if there was an error during connection.
// 	if err != nil {
// 		t.Errorf("Failed to connect to Firebase: %v", err)
// 	}

//		// Check if the Firebase Realtime Database client is initialized.
//		if db.Client == nil {
//			t.Error("Firebase Realtime Database client is not initialized")
//		}
//	}
func TestDatabase(t *testing.T) {
	ctx := context.Background()

	// configure database URL
	conf := &firebase.Config{
		DatabaseURL: "https://edusync-test-default-rtdb.firebaseio.com/",
	}

	// fetch service account key
	opt := option.WithCredentialsFile("edusync-test-firebase-adminsdk-hk5kl-9af0162b09.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		t.Errorf("error in initializing firebase app: %v", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		t.Errorf("error in creating firebase DB client: %v", err)
	}

	// create ref at path students/:userId
	ref := client.NewRef("students/" + fmt.Sprint(1))

	// Test case 1: Successful set operation
	data := map[string]interface{}{
		"name":       "Jane Doe",
		"age":        "7",
		"class":      "Tech Explorer",
		"instructor": "Scott Smith",
	}
	if err := ref.Set(ctx, data); err != nil {
		t.Errorf("error in setting data: %v", err)
	}

	// Test case 2: Get the set data
	var getData map[string]interface{}
	if err := ref.Get(ctx, &getData); err != nil {
		t.Errorf("error in getting data: %v", err)
	}
	if getData["name"] != "Jane Doe" {
		t.Errorf("expected name to be 'Jane Doe', got %v", getData["name"])
	}

	// Test case 3: Update the data
	updateData := map[string]interface{}{
		"name": "John Doe",
	}
	if err := ref.Update(ctx, updateData); err != nil {
		t.Errorf("error in updating data: %v", err)
	}

	// Test case 4: Get the updated data
	var updatedData map[string]interface{}
	if err := ref.Get(ctx, &updatedData); err != nil {
		t.Errorf("error in getting updated data: %v", err)
	}
	if updatedData["name"] != "John Doe" {
		t.Errorf("expected name to be 'John Doe', got %v", updatedData["name"])
	}

	// Test case 5: Delete the data
	if err := ref.Delete(ctx); err != nil {
		t.Errorf("error in deleting data: %v", err)
	}

	// Test case 6: Get the deleted data (should return an error)
	var deletedData map[string]interface{}
	// if err := ref.Get(ctx, &deletedData); err == nil {
	// 	t.Errorf("expected error in getting deleted data, but got nil")
	// }
	if err := ref.Get(ctx, &deletedData); err == nil {
		// If no error, check if the data is actually deleted
		if deletedData != nil {
			t.Errorf("Expected data to be deleted, but got %v", deletedData)
		}
	} else {
		// Expecting an error, which indicates the data was not found
		t.Logf("Received expected error after deletion: %v", err)
	}
}
