package main

import (
	// "context"
	// "fmt"

	// "os"
	"reflect"

	// "strings"
	"testing"
	// firebase "firebase.google.com/go"
	// "google.golang.org/api/option"
)

func TestInitializeFirebase(t *testing.T) {
	// Test case 1: FirebaseClient is set correctly
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}
	if firebaseClient == nil {
		t.Fatal("FirebaseClient is not set")
	}

	// Run tests
	// os.Exit(t.Run())
}

// Testing for student CRUD operations
func TestCreateStudent(t *testing.T) {
	// Initialize Firebase client
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	// Create a new student
	student := Student{
		GoogleID:   "test-student",
		Name:       "John Doe",
		Email:      "johndoe@example.com",
		Age:        12,
		Class:      "TE",
		Instructor: "Awesomeness",
		ParentName: "Jane Doe",
		Role:       "Student",
	}

	err = createStudent(student.GoogleID, student)
	if err != nil {
		t.Fatalf("Error creating student: %v", err)
	}

	// Read the created student
	readStudent, err := readStudent(student.GoogleID)
	if err != nil {
		t.Fatalf("Error reading student: %v", err)
	}

	// Assert that the created and read students are equal
	if !reflect.DeepEqual(student, readStudent) {
		t.Error("Created and read students are not equal")
	}
}

func TestReadStudent(t *testing.T) {
	googleID := "test-student"

	student, err := readStudent(googleID)
	if err != nil {
		t.Fatalf("Failed to read student: %v", err)
	}

	if student.GoogleID != googleID {
		t.Fatalf("Expected GoogleID %v, got %v", googleID, student.GoogleID)
	}
}

func TestUpdateStudent(t *testing.T) {
	// Initialize Firebase client
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	// Update the student's email
	updates := map[string]interface{}{
		"email": "johndoe@nk.com",
	}

	err = updateStudent("test-student", updates)
	if err != nil {
		t.Fatalf("Error updating student: %v", err)
	}

	// Read the updated student
	readStudent, err := readStudent("test-student")
	if err != nil {
		t.Fatalf("Error reading student: %v", err)
	}

	// Assert that the updated student's email is correct
	if readStudent.Email != updates["email"] {
		t.Errorf("Updated student's email is incorrect. Expected: %v, Got: %v", updates["email"], readStudent.Email)
	}
}

func TestDeleteStudent(t *testing.T) {
	googleID := "test-student"

	// Initialize Firebase client
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	// Delete the student
	err = deleteStudent(googleID)
	if err != nil {
		t.Fatalf("Error deleting student: %v", err)
	}

	// Try to read the deleted student
	// _, err = readStudent(googleID)
	// if err == nil {
	// 	t.Error("Deleted student still exists")
	// }
}

// Testing for instructor CRUD operations
func TestCreateInstructor(t *testing.T) {
	// Initialize Firebase client
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	// Create a new instructor
	instructor := Instructor{
		GoogleID:         "test-instructor",
		Name:             "Awesomeness",
		ContactNumber:    "99999999",
		Email:            "awesome_instructor@nk.com",
		BasePay:          15,
		NumberOfStudents: 24,
		Role:             "Instructor",
	}

	err = createInstructor(instructor.GoogleID, instructor)
	if err != nil {
		t.Fatalf("Error creating instructor: %v", err)
	}

	// Read the created instructor
	readInstructor, err := readInstructor(instructor.GoogleID)
	if err != nil {
		t.Fatalf("Error reading instructor: %v", err)
	}

	// Assert that the created and read instructor are equal
	if !reflect.DeepEqual(instructor, readInstructor) {
		t.Error("Created and read instructors are not equal")
	}
}

func TestReadInstructor(t *testing.T) {
	googleID := "test-instructor"

	instructor, err := readInstructor(googleID)
	if err != nil {
		t.Fatalf("Failed to read instructor: %v", err)
	}

	if instructor.GoogleID != googleID {
		t.Fatalf("Expected GoogleID %v, got %v", googleID, instructor.GoogleID)
	}
}

func TestUpdateInstructor(t *testing.T) {
	// Initialize Firebase client
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	// Update the instructor's email
	updates := map[string]interface{}{
		"email": "amazing_instructor@nk.com",
	}

	err = updateInstructor("test-instructor", updates)
	if err != nil {
		t.Fatalf("Error updating instructor: %v", err)
	}

	// Read the updated instructor
	readInstructor, err := readInstructor("test-instructor")
	if err != nil {
		t.Fatalf("Error reading instructor: %v", err)
	}

	// Assert that the updated instructor's email is correct
	if readInstructor.Email != updates["email"] {
		t.Errorf("Updated instructor's email is incorrect. Expected: %v, Got: %v", updates["email"], readInstructor.Email)
	}
}

func TestDeleteInstructor(t *testing.T) {
	googleID := "test-instructor"

	// Initialize Firebase client
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	// Delete the instructor
	err = deleteInstructor(googleID)
	if err != nil {
		t.Fatalf("Error deleting instructor: %v", err)
	}

	// Try to read the deleted instructor
	// _, err = readInstructor(googleID)
	// if err == nil {
	// 	t.Error("Deleted instructor still exists")
	// }
}

// Testing for admin CRUD operations
func TestCreateAdmin(t *testing.T) {
	// Initialize Firebase client
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	// Create a new admin
	admin := Admin{
		GoogleID:      "test-admin",
		Name:          "Awesomeness",
		ContactNumber: "99999999",
		Email:         "awesome_admin@nk.com",
		BasePay:       15,
		Incentive:     24,
		Role:          "Admin",
	}

	err = createAdmin(admin.GoogleID, admin)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}

	// Read the created admin
	readAdmin, err := readAdmin(admin.GoogleID)
	if err != nil {
		t.Fatalf("Error reading admin: %v", err)
	}

	// Assert that the created and read admin are equal
	if !reflect.DeepEqual(admin, readAdmin) {
		t.Error("Created and read admins are not equal")
	}
}

func TestReadAdmin(t *testing.T) {
	googleID := "test-admin"

	admin, err := readAdmin(googleID)
	if err != nil {
		t.Fatalf("Failed to read instructor: %v", err)
	}

	if admin.GoogleID != googleID {
		t.Fatalf("Expected GoogleID %v, got %v", googleID, admin.GoogleID)
	}
}

func TestUpdateAdmin(t *testing.T) {
	// Initialize Firebase client
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	// Update the admin's email
	updates := map[string]interface{}{
		"email": "amazing_admin@nk.com",
	}

	err = updateAdmin("test-admin", updates)
	if err != nil {
		t.Fatalf("Error updating admin: %v", err)
	}

	// Read the updated admin
	readAdmin, err := readAdmin("test-admin")
	if err != nil {
		t.Fatalf("Error reading admin: %v", err)
	}

	// Assert that the updated admin's email is correct
	if readAdmin.Email != updates["email"] {
		t.Errorf("Updated admin's email is incorrect. Expected: %v, Got: %v", updates["email"], readAdmin.Email)
	}
}

func TestDeleteAdmin(t *testing.T) {
	googleID := "test-admin"

	// Initialize Firebase client
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	// Delete the admin
	err = deleteAdmin(googleID)
	if err != nil {
		t.Fatalf("Error deleting admin: %v", err)
	}

	// Try to read the deleted admin
	// _, err = readAdmin(googleID)
	// if err == nil {
	// 	t.Error("Deleted admin still exists")
	// }
}

// Testing for parent CRUD operations
func TestCreateParent(t *testing.T) {
	// Initialize Firebase client
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	// Create a new parent
	parent := Parent{
		GoogleID:      "test-parent",
		Name:          "Awesomeness",
		ContactNumber: "99999999",
		Email:         "janedoe@nk.com",
		Role:          "Parent",
	}

	err = createParent(parent.GoogleID, parent)
	if err != nil {
		t.Fatalf("Error creating parent: %v", err)
	}

	// Read the created parent
	readParent, err := readParent(parent.GoogleID)
	if err != nil {
		t.Fatalf("Error reading parent: %v", err)
	}

	// Assert that the created and read parent are equal
	if !reflect.DeepEqual(parent, readParent) {
		t.Error("Created and read parents are not equal")
	}
}

func TestReadParent(t *testing.T) {
	googleID := "test-parent"

	parent, err := readParent(googleID)
	if err != nil {
		t.Fatalf("Failed to read parent: %v", err)
	}

	if parent.GoogleID != googleID {
		t.Fatalf("Expected GoogleID %v, got %v", googleID, parent.GoogleID)
	}
}

func TestUpdateParent(t *testing.T) {
	// Initialize Firebase client
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	// Update the parent's email
	updates := map[string]interface{}{
		"email": "janedoe_parent@nk.com",
	}

	err = updateParent("test-parent", updates)
	if err != nil {
		t.Fatalf("Error updating parent: %v", err)
	}

	// Read the updated parent
	readParent, err := readParent("test-parent")
	if err != nil {
		t.Fatalf("Error reading parent: %v", err)
	}

	// Assert that the updated parent's email is correct
	if readParent.Email != updates["email"] {
		t.Errorf("Updated parent's email is incorrect. Expected: %v, Got: %v", updates["email"], readParent.Email)
	}
}

func TestDeleteParent(t *testing.T) {
	googleID := "test-parent"

	// Initialize Firebase client
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	// Delete the parent
	err = deleteParent(googleID)
	if err != nil {
		t.Fatalf("Error deleting parent: %v", err)
	}

	// Try to read the deleted parent
	// _, err = readParent(googleID)
	// if err == nil {
	// 	t.Error("Deleted parent still exists")
	// }
}

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"testing"
// 	"time"

// 	firebase "firebase.google.com/go"
// 	"google.golang.org/api/option"
// )

// // Test using test database
// // func TestDatabase(t *testing.T) {
// // 	ctx := context.Background()

// // 	// configure database URL
// // 	conf := &firebase.Config{
// // 		DatabaseURL: "https://edusync-test-default-rtdb.firebaseio.com/",
// // 	}

// // 	// fetch service account key
// // 	opt := option.WithCredentialsFile("edusync-test-firebase-adminsdk-hk5kl-9af0162b09.json")

// // 	app, err := firebase.NewApp(ctx, conf, opt)
// // 	if err != nil {
// // 		t.Errorf("error in initializing firebase app: %v", err)
// // 	}

// // 	client, err := app.Database(ctx)
// // 	if err != nil {
// // 		t.Errorf("error in creating firebase DB client: %v", err)
// // 	}

// // 	// create ref at path students/:userId
// // 	ref := client.NewRef("students/" + fmt.Sprint(1))

// // 	// Test case 1: Successful set operation
// // 	data := map[string]interface{}{
// // 		"name":       "Jane Doe",
// // 		"age":        "7",
// // 		"class":      "Tech Explorer",
// // 		"instructor": "Scott Smith",
// // 	}
// // 	if err := ref.Set(ctx, data); err != nil {
// // 		t.Errorf("error in setting data: %v", err)
// // 	}

// // 	// Test case 2: Get the set data
// // 	var getData map[string]interface{}
// // 	if err := ref.Get(ctx, &getData); err != nil {
// // 		t.Errorf("error in getting data: %v", err)
// // 	}
// // 	if getData["name"] != "Jane Doe" {
// // 		t.Errorf("expected name to be 'Jane Doe', got %v", getData["name"])
// // 	}

// // 	// Test case 3: Update the data
// // 	updateData := map[string]interface{}{
// // 		"name": "John Doe",
// // 	}
// // 	if err := ref.Update(ctx, updateData); err != nil {
// // 		t.Errorf("error in updating data: %v", err)
// // 	}

// // 	// Test case 4: Get the updated data
// // 	var updatedData map[string]interface{}
// // 	if err := ref.Get(ctx, &updatedData); err != nil {
// // 		t.Errorf("error in getting updated data: %v", err)
// // 	}
// // 	if updatedData["name"] != "John Doe" {
// // 		t.Errorf("expected name to be 'John Doe', got %v", updatedData["name"])
// // 	}

// // 	// Test case 5: Delete the data
// // 	if err := ref.Delete(ctx); err != nil {
// // 		t.Errorf("error in deleting data: %v", err)
// // 	}

// // 	// Test case 6: Get the deleted data (should return an error)
// // 	var deletedData map[string]interface{}
// // 	// if err := ref.Get(ctx, &deletedData); err == nil {
// // 	// 	t.Errorf("expected error in getting deleted data, but got nil")
// // 	// }
// // 	if err := ref.Get(ctx, &deletedData); err == nil {
// // 		// If no error, check if the data is actually deleted
// // 		if deletedData != nil {
// // 			t.Errorf("Expected data to be deleted, but got %v", deletedData)
// // 		}
// // 	} else {
// // 		// Expecting an error, which indicates the data was not found
// // 		t.Logf("Received expected error after deletion: %v", err)
// // 	}
// // }

// // Test using actual database
// func TestDatabaseCRUD(t *testing.T) {
// 	ctx := context.Background()
// 	databaseURL := goDotEnvVariable("DATABASE_URL")
// 	// databaseURL, found := os.LookupEnv("DATABASE_URL")
// 	// if !found {
// 	// 	log.Fatalf("DATABASE_URL is not set in the environment variables")
// 	// }
// 	conf := &firebase.Config{DatabaseURL: databaseURL}
// 	opt := option.WithCredentialsFile("edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json")
// 	// opt := option.WithCredentialsFile("$HOME/secrets/edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json")

// 	app, err := firebase.NewApp(ctx, conf, opt)
// 	if err != nil {
// 		log.Fatalln("error in initializing firebase app: ", err)
// 	}

// 	client, err := app.Database(ctx)
// 	if err != nil {
// 		log.Fatalln("error in creating firebase DB client: ", err)
// 	}

// 	// Student operations
// 	// student := NewStudent("Jane Doe", 7, 119.5, "jane_doe@nk.com", "91234567", "Tech Explorer", "Scott Smith", "Jackie Doe")
// 	// err = createStudent(client, student.ID.String(), student)
// 	googleIDStudent := "google-id-student"
// 	student := NewStudent(googleIDStudent, "Jane Doe", 7, 119.5, "jane_doe@nk.com", "91234567", "Tech Explorer", "Scott Smith", "Jackie Doe")
// 	err = createStudent(client, student.GoogleID, student)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Student added/updated successfully!")

// 	// readStudent, err := readStudent(client, student.ID.String())
// 	readStudent, err := readStudent(client, student.GoogleID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Student read successfully:", readStudent)

// 	studentUpdates := map[string]interface{}{
// 		"class":      "Tech Explorer 2",
// 		"updated_at": time.Now(),
// 	}
// 	// err = updateStudent(client, student.ID.String(), studentUpdates)
// 	err = updateStudent(client, student.GoogleID, studentUpdates)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Student updated successfully!")

// 	// err = deleteStudent(client, student.ID.String())
// 	err = deleteStudent(client, student.GoogleID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Student deleted successfully!")

// 	// Instructor operations
// 	// instructor := NewInstructor("Scott Smith", "123-456-7890", "scott@example.com", 50000.00, 10)
// 	// err = createInstructor(client, instructor.ID.String(), instructor)
// 	googleIDInstructor := "google-id-instructor"
// 	instructor := NewInstructor(googleIDInstructor, "Scott Smith", "123-456-7890", "scott@example.com", 50000.00, 10)
// 	err = createInstructor(client, instructor.GoogleID, instructor)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Instructor added/updated successfully!")

// 	// readInstructor, err := readInstructor(client, instructor.ID.String())
// 	readInstructor, err := readInstructor(client, instructor.GoogleID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Instructor read successfully:", readInstructor)

// 	instructorUpdates := map[string]interface{}{
// 		"base_pay":   55000.00,
// 		"updated_at": time.Now(),
// 	}
// 	// err = updateInstructor(client, instructor.ID.String(), instructorUpdates)
// 	err = updateInstructor(client, instructor.GoogleID, instructorUpdates)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Instructor updated successfully!")

// 	// err = deleteInstructor(client, instructor.ID.String())
// 	err = deleteInstructor(client, instructor.GoogleID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Instructor deleted successfully!")

// 	// Parent operations
// 	// parent := NewParent("Jackie Doe", "jackjack@example.com", "98765432")
// 	// err = createParent(client, parent.ID.String(), parent)
// 	googleIDParent := "google-id-parent"
// 	parent := NewParent(googleIDParent, "Jackie Doe", "jackjack@example.com", "98765432")
// 	err = createParent(client, parent.GoogleID, parent)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Parent added/updated successfully!")

// 	// readParent, err := readParent(client, parent.ID.String())
// 	readParent, err := readParent(client, parent.GoogleID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Parent read successfully:", readParent)

// 	parentUpdates := map[string]interface{}{
// 		"email":      "jackiejack@nk.com",
// 		"updated_at": time.Now(),
// 	}
// 	// err = updateParent(client, parent.ID.String(), parentUpdates)
// 	err = updateParent(client, parent.GoogleID, parentUpdates)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Parent updated successfully!")

// 	// err = deleteParent(client, parent.ID.String())
// 	err = deleteParent(client, parent.GoogleID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Parent deleted successfully!")
// }
