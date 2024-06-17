package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Test using test database
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

// Test using actual database
func TestDatabaseCRUD(t *testing.T) {
	ctx := context.Background()
	databaseURL := goDotEnvVariable("DATABASE_URL")
	conf := &firebase.Config{DatabaseURL: databaseURL}
	opt := option.WithCredentialsFile("edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("error in initializing firebase app: ", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase DB client: ", err)
	}

	// Student operations
	student := NewStudent("Jane Doe", 7, 119.5, "jane_doe@nk.com", "91234567", "Tech Explorer", "Scott Smith", "Jackie Doe")
	err = createStudent(client, student.ID.String(), student)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Student added/updated successfully!")

	readStudent, err := readStudent(client, student.ID.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Student read successfully:", readStudent)

	studentUpdates := map[string]interface{}{
		"class":      "Tech Explorer 2",
		"updated_at": time.Now(),
	}
	err = updateStudent(client, student.ID.String(), studentUpdates)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Student updated successfully!")

	err = deleteStudent(client, student.ID.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Student deleted successfully!")

	// Instructor operations
	instructor := NewInstructor("Scott Smith", "123-456-7890", "scott@example.com", 50000.00, 10)
	err = createInstructor(client, instructor.ID.String(), instructor)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Instructor added/updated successfully!")

	readInstructor, err := readInstructor(client, instructor.ID.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Instructor read successfully:", readInstructor)

	instructorUpdates := map[string]interface{}{
		"base_pay":   55000.00,
		"updated_at": time.Now(),
	}
	err = updateInstructor(client, instructor.ID.String(), instructorUpdates)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Instructor updated successfully!")

	err = deleteInstructor(client, instructor.ID.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Instructor deleted successfully!")

	// Parent operations
	parent := NewParent("Jackie Doe", "jackjack@example.com", "98765432")
	err = createParent(client, parent.ID.String(), parent)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parent added/updated successfully!")

	readParent, err := readParent(client, parent.ID.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parent read successfully:", readParent)

	parentUpdates := map[string]interface{}{
		"email":      "jackiejack@nk.com",
		"updated_at": time.Now(),
	}
	err = updateParent(client, parent.ID.String(), parentUpdates)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parent updated successfully!")

	err = deleteParent(client, parent.ID.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parent deleted successfully!")
}
