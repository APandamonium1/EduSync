package main

import (
	"context"
	"fmt"
	"time"

	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

func database() {
	ctx := context.Background()

	// configure database URL
	conf := &firebase.Config{
		DatabaseURL: "https://edusync-test-default-rtdb.firebaseio.com/",
	}

	// fetch service account key
	opt := option.WithCredentialsFile("edusync-test-firebase-adminsdk-hk5kl-9af0162b09.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("error in initializing firebase app: ", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase DB client: ", err)
	}

	// Student operations
	student := NewStudent("Jane Doe", 7, "Tech Explorer", "91234567", "Scott Smith")
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
}

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
