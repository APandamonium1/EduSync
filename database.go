package main

import (
	"context"
	"fmt"
	"time"

	"log"
	"os"

	// "github.com/joho/godotenv"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

// Use godot package to load/read the .env file and
// return the value of the key
// func goDotEnvVariable(key string) string {

// 	// load .env file
// 	err := godotenv.Load(".env")

// 	if err != nil {
// 		log.Fatalf("Error loading .env file")
// 	}

// 	return os.Getenv(key)
// }

func database() {
	// Find home directory.
	// home, err := os.Getwd()
	// if err != nil {
	// 	return err
	// }

	ctx := context.Background()

	// configure database URL
	// databaseURL := goDotEnvVariable("DATABASE_URL")
	// if databaseURL == "" {
	// 	return fmt.Errorf("DATABASE_URL is not set in the .env file")
	// }
	databaseURL, found := os.LookupEnv("DATABASE_URL")
	if !found {
		log.Fatalf("DATABASE_URL is not set in the environment variables")
	}
	conf := &firebase.Config{DatabaseURL: databaseURL}

	// conf := &firebase.Config{
	// 	DatabaseURL: "https://edusync-test-default-rtdb.firebaseio.com/",
	// }

	// Set up the Firebase app with the provided JSON file containing the service account key.
	// opt := option.WithCredentialsFile(home + "edusync-test-firebase-adminsdk-hk5kl-9af0162b09.json")

	// fetch service account key
	// opt := option.WithCredentialsFile("edusync-test-firebase-adminsdk-hk5kl-9af0162b09.json")
	opt := option.WithCredentialsFile("edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json")
	// opt := option.WithCredentialsFile("$HOME/secrets/edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json")

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

// Student CRUD
func createStudent(client *db.Client, userId string, student Student) error {
	ref := client.NewRef("students/" + userId)
	return ref.Set(context.TODO(), student)
}

func readStudent(client *db.Client, userId string) (Student, error) {
	ref := client.NewRef("students/" + userId)
	var student Student
	if err := ref.Get(context.TODO(), &student); err != nil {
		return Student{}, err
	}
	return student, nil
}

func updateStudent(client *db.Client, userId string, updates map[string]interface{}) error {
	ref := client.NewRef("students/" + userId)
	return ref.Update(context.TODO(), updates)
}

func deleteStudent(client *db.Client, userId string) error {
	ref := client.NewRef("students/" + userId)
	return ref.Delete(context.TODO())
}

// Instructor CRUD
func createInstructor(client *db.Client, userId string, instructor Instructor) error {
	ref := client.NewRef("instructors/" + userId)
	return ref.Set(context.TODO(), instructor)
}

func readInstructor(client *db.Client, userId string) (Instructor, error) {
	ref := client.NewRef("instructors/" + userId)
	var instructor Instructor
	if err := ref.Get(context.TODO(), &instructor); err != nil {
		return Instructor{}, err
	}
	return instructor, nil
}

func updateInstructor(client *db.Client, userId string, updates map[string]interface{}) error {
	ref := client.NewRef("instructors/" + userId)
	return ref.Update(context.TODO(), updates)
}

func deleteInstructor(client *db.Client, userId string) error {
	ref := client.NewRef("instructors/" + userId)
	return ref.Delete(context.TODO())
}

// Parent CRUD
func createParent(client *db.Client, userId string, parent Parent) error {
	ref := client.NewRef("parents/" + userId)
	return ref.Set(context.TODO(), parent)
}

func readParent(client *db.Client, userId string) (Parent, error) {
	ref := client.NewRef("parents/" + userId)
	var parent Parent
	if err := ref.Get(context.TODO(), &parent); err != nil {
		return Parent{}, err
	}
	return parent, nil
}

func updateParent(client *db.Client, userId string, updates map[string]interface{}) error {
	ref := client.NewRef("parents/" + userId)
	return ref.Update(context.TODO(), updates)
}

func deleteParent(client *db.Client, userId string) error {
	ref := client.NewRef("parents/" + userId)
	return ref.Delete(context.TODO())
}
