package main

import (
	"context"
	"fmt"

	"log"
	"os"

	// "github.com/joho/godotenv"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

// Use godot package to load/read the .env file and
// return the value of the key (for local env)
// func goDotEnvVariable(key string) string {

// 	// load .env file
// 	err := godotenv.Load(".env")

// 	if err != nil {
// 		log.Fatalf("Error loading .env file")
// 	}

// 	return os.Getenv(key)
// }

// Initialize Firebase client
var firebaseClient *db.Client

// InitializeFirebase initializes the Firebase app and sets the global firebaseClient variable
func initializeFirebase() error {
	ctx := context.Background()

	databaseURL, found := os.LookupEnv("DATABASE_URL")
	if !found {
		log.Fatalf("DATABASE_URL is not set in the environment variables")
	}
	// databaseURL := goDotEnvVariable("DATABASE_URL")
	// if databaseURL == "" {
	// 	return fmt.Errorf("DATABASE_URL is not set in the environment variables")
	// }

	conf := &firebase.Config{DatabaseURL: databaseURL}

	opt := option.WithCredentialsFile("edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json")
	// opt := option.WithCredentialsFile("edusync-test-firebase-adminsdk-hk5kl-9af0162b09.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		return fmt.Errorf("error creating firebase DB client: %v", err)
	}

	firebaseClient = client
	return nil
}

// Student CRUD
func createStudent(googleID string, student Student) error {
	ref := firebaseClient.NewRef("students/" + googleID)
	if err := ref.Set(context.TODO(), student); err != nil {
		return fmt.Errorf("error creating student: %v", err)
	}
	return ref.Set(context.TODO(), student)
}

func readStudent(googleID string) (Student, error) {
	ref := firebaseClient.NewRef("students/" + googleID)
	var student Student
	if err := ref.Get(context.TODO(), &student); err != nil {
		return Student{}, fmt.Errorf("error reading student: %v", err)
	}
	return student, nil
}

func updateStudent(googleID string, updates map[string]interface{}) error {
	ref := firebaseClient.NewRef("students/" + googleID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating student: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteStudent(googleID string) error {
	ref := firebaseClient.NewRef("students/" + googleID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting student: %v", err)
	}
	return ref.Delete(context.TODO())
}

// Instructor CRUD
func createInstructor(googleID string, instructor Instructor) error {
	ref := firebaseClient.NewRef("instructors/" + googleID)
	if err := ref.Set(context.TODO(), instructor); err != nil {
		return fmt.Errorf("error creating instructor: %v", err)
	}
	return ref.Set(context.TODO(), instructor)
}

func readInstructor(googleID string) (Instructor, error) {
	ref := firebaseClient.NewRef("instructors/" + googleID)
	var instructor Instructor
	if err := ref.Get(context.TODO(), &instructor); err != nil {
		return Instructor{}, fmt.Errorf("error reading instructor: %v", err)
	}
	return instructor, nil
}

func updateInstructor(googleID string, updates map[string]interface{}) error {
	ref := firebaseClient.NewRef("instructors/" + googleID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating instructor: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteInstructor(googleID string) error {
	ref := firebaseClient.NewRef("instructors/" + googleID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting instructor: %v", err)
	}
	return ref.Delete(context.TODO())
}

// Admin CRUD
func createAdmin(googleID string, admin Admin) error {
	ref := firebaseClient.NewRef("admins/" + googleID)
	if err := ref.Set(context.TODO(), admin); err != nil {
		return fmt.Errorf("error creating admin: %v", err)
	}
	return ref.Set(context.TODO(), admin)
}

func readAdmin(googleID string) (Admin, error) {
	ref := firebaseClient.NewRef("admins/" + googleID)
	var admin Admin
	if err := ref.Get(context.TODO(), &admin); err != nil {
		return Admin{}, fmt.Errorf("error reading admin: %v", err)
	}
	return admin, nil
}

func updateAdmin(googleID string, updates map[string]interface{}) error {
	ref := firebaseClient.NewRef("admins/" + googleID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating admin: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteAdmin(googleID string) error {
	ref := firebaseClient.NewRef("admins/" + googleID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting admin: %v", err)
	}
	return ref.Delete(context.TODO())
}

// Parent CRUD
func createParent(googleID string, parent Parent) error {
	ref := firebaseClient.NewRef("parents/" + googleID)
	if err := ref.Set(context.TODO(), parent); err != nil {
		return fmt.Errorf("error creating parent: %v", err)
	}
	return ref.Set(context.TODO(), parent)
}

func readParent(googleID string) (Parent, error) {
	ref := firebaseClient.NewRef("parents/" + googleID)
	var parent Parent
	if err := ref.Get(context.TODO(), &parent); err != nil {
		return Parent{}, fmt.Errorf("error reading parent: %v", err)
	}
	return parent, nil
}

func updateParent(googleID string, updates map[string]interface{}) error {
	ref := firebaseClient.NewRef("parents/" + googleID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating parent: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteParent(googleID string) error {
	ref := firebaseClient.NewRef("parents/" + googleID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting parent: %v", err)
	}
	return ref.Delete(context.TODO())
}

// func database() {
// 	// Find home directory.
// 	// home, err := os.Getwd()
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	ctx := context.Background()

// 	// configure database URL
// 	// databaseURL := goDotEnvVariable("DATABASE_URL")
// 	// if databaseURL == "" {
// 	// 	return fmt.Errorf("DATABASE_URL is not set in the .env file")
// 	// }
// 	// databaseURL, found := os.LookupEnv("DATABASE_URL")
// 	// if !found {
// 	// 	log.Fatalf("DATABASE_URL is not set in the environment variables")
// 	// }
// 	// conf := &firebase.Config{DatabaseURL: databaseURL}

// 	conf := &firebase.Config{
// 		DatabaseURL: "https://edusync-test-default-rtdb.firebaseio.com/",
// 	}

// 	// Set up the Firebase app with the provided JSON file containing the service account key.
// 	// opt := option.WithCredentialsFile(home + "edusync-test-firebase-adminsdk-hk5kl-9af0162b09.json")
// 	opt := option.WithCredentialsFile("edusync-test-firebase-adminsdk-hk5kl-9af0162b09.json")
// 	// opt := option.WithCredentialsFile("edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json")
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

// Student CRUD
// func createStudent(client *db.Client, userId string, student Student) error {
// 	ref := client.NewRef("students/" + userId)
// 	return ref.Set(context.TODO(), student)
// }

// func readStudent(client *db.Client, userId string) (Student, error) {
// 	ref := client.NewRef("students/" + userId)
// 	var student Student
// 	if err := ref.Get(context.TODO(), &student); err != nil {
// 		return Student{}, err
// 	}
// 	return student, nil
// }

// func updateStudent(client *db.Client, userId string, updates map[string]interface{}) error {
// 	ref := client.NewRef("students/" + userId)
// 	return ref.Update(context.TODO(), updates)
// }

//	func deleteStudent(client *db.Client, userId string) error {
//		ref := client.NewRef("students/" + userId)
//		return ref.Delete(context.TODO())
//	}

// Instructor CRUD
// func createInstructor(client *db.Client, userId string, instructor Instructor) error {
// 	ref := client.NewRef("instructors/" + userId)
// 	return ref.Set(context.TODO(), instructor)
// }

// func readInstructor(client *db.Client, userId string) (Instructor, error) {
// 	ref := client.NewRef("instructors/" + userId)
// 	var instructor Instructor
// 	if err := ref.Get(context.TODO(), &instructor); err != nil {
// 		return Instructor{}, err
// 	}
// 	return instructor, nil
// }

// func updateInstructor(client *db.Client, userId string, updates map[string]interface{}) error {
// 	ref := client.NewRef("instructors/" + userId)
// 	return ref.Update(context.TODO(), updates)
// }

//	func deleteInstructor(client *db.Client, userId string) error {
//		ref := client.NewRef("instructors/" + userId)
//		return ref.Delete(context.TODO())
//	}

// Parent CRUD
// func createParent(client *db.Client, userId string, parent Parent) error {
// 	ref := client.NewRef("parents/" + userId)
// 	return ref.Set(context.TODO(), parent)
// }

// func readParent(client *db.Client, userId string) (Parent, error) {
// 	ref := client.NewRef("parents/" + userId)
// 	var parent Parent
// 	if err := ref.Get(context.TODO(), &parent); err != nil {
// 		return Parent{}, err
// 	}
// 	return parent, nil
// }

// func updateParent(client *db.Client, userId string, updates map[string]interface{}) error {
// 	ref := client.NewRef("parents/" + userId)
// 	return ref.Update(context.TODO(), updates)
// }

//	func deleteParent(client *db.Client, userId string) error {
//		ref := client.NewRef("parents/" + userId)
//		return ref.Delete(context.TODO())
//	}

// // Student CRUD Ver 2
// func createStudent(client *db.Client, googleID string, student Student) error {
// 	ref := client.NewRef("students/" + googleID)
// 	return ref.Set(context.TODO(), student)
// }

// func readStudent(client *db.Client, googleID string) (Student, error) {
// 	ref := client.NewRef("students/" + googleID)
// 	var student Student
// 	if err := ref.Get(context.TODO(), &student); err != nil {
// 		return Student{}, err
// 	}
// 	return student, nil
// }

// func updateStudent(client *db.Client, googleID string, updates map[string]interface{}) error {
// 	ref := client.NewRef("students/" + googleID)
// 	return ref.Update(context.TODO(), updates)
// }

// func deleteStudent(client *db.Client, googleID string) error {
// 	ref := client.NewRef("students/" + googleID)
// 	return ref.Delete(context.TODO())
// }

// // Instructor CRUD Ver 2
// func createInstructor(client *db.Client, googleID string, instructor Instructor) error {
// 	ref := client.NewRef("instructors/" + googleID)
// 	return ref.Set(context.TODO(), instructor)
// }

// func readInstructor(client *db.Client, googleID string) (Instructor, error) {
// 	ref := client.NewRef("instructors/" + googleID)
// 	var instructor Instructor
// 	if err := ref.Get(context.TODO(), &instructor); err != nil {
// 		return Instructor{}, err
// 	}
// 	return instructor, nil
// }

// func updateInstructor(client *db.Client, googleID string, updates map[string]interface{}) error {
// 	ref := client.NewRef("instructors/" + googleID)
// 	return ref.Update(context.TODO(), updates)
// }

// func deleteInstructor(client *db.Client, googleID string) error {
// 	ref := client.NewRef("instructors/" + googleID)
// 	return ref.Delete(context.TODO())
// }

// // Parent CRUD Ver 2
// func createParent(client *db.Client, googleID string, parent Parent) error {
// 	ref := client.NewRef("parents/" + googleID)
// 	return ref.Set(context.TODO(), parent)
// }

// func readParent(client *db.Client, googleID string) (Parent, error) {
// 	ref := client.NewRef("parents/" + googleID)
// 	var parent Parent
// 	if err := ref.Get(context.TODO(), &parent); err != nil {
// 		return Parent{}, err
// 	}
// 	return parent, nil
// }

// func updateParent(client *db.Client, googleID string, updates map[string]interface{}) error {
// 	ref := client.NewRef("parents/" + googleID)
// 	return ref.Update(context.TODO(), updates)
// }

// func deleteParent(client *db.Client, googleID string) error {
// 	ref := client.NewRef("parents/" + googleID)
// 	return ref.Delete(context.TODO())
// }
