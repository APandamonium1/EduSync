package main

import (
	"context"
	"fmt"
	"net/http"

	"log"
	"os"

	"github.com/joho/godotenv"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

// Use godot package to load/read the .env file and
// return the value of the key (for local env)
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// Initialize Firebase client
var firebaseClient *db.Client

// InitializeFirebase initializes the Firebase app and sets the global firebaseClient variable
func initializeFirebase() error {
	ctx := context.Background()

	// databaseURL, found := os.LookupEnv("DATABASE_URL")
	// if !found {
	// 	log.Fatalf("DATABASE_URL is not set in the environment variables")
	// }
	databaseURL := goDotEnvVariable("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL is not set in the environment variables")
	}

	conf := &firebase.Config{DatabaseURL: databaseURL}

	// opt := option.WithCredentialsFile("edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json")
	opt := option.WithCredentialsFile("edusync-test-firebase-adminsdk-hk5kl-9af0162b09.json")

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

// Utility function to get current user
func getCurrentUser(req *http.Request) (User, error) {
	// Implement session or context based user retrieval
	return User{}, nil
}

// Utility functions to check roles
func isAdmin(user User) bool {
	return user.Role == "Admin"
}

func isSelf(user User, googleID string) bool {
	return user.GoogleID == googleID
}

// CRUD operations with role checks

// Student CRUD
func createStudent(currentUser User, googleID string, student Student) error {
	if currentUser.Role != "Admin" {
		return fmt.Errorf("unauthorized access: only admins can create students")
	}
	ref := firebaseClient.NewRef("students/" + googleID)
	if err := ref.Set(context.TODO(), student); err != nil {
		return fmt.Errorf("error creating student: %v", err)
	}
	return ref.Set(context.TODO(), student)
}

func readStudent(currentUser User, googleID string) (Student, error) {
	if !isAdmin(currentUser) && !isSelf(currentUser, googleID) {
		return Student{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}
	ref := firebaseClient.NewRef("students/" + googleID)
	var student Student
	if err := ref.Get(context.TODO(), &student); err != nil {
		return Student{}, fmt.Errorf("error reading student: %v", err)
	}
	return student, nil
}

func updateStudent(currentUser User, googleID string, updates map[string]interface{}) error {
	if !isAdmin(currentUser) && !isSelf(currentUser, googleID) {
		return fmt.Errorf("unauthorized access: you can only update your own details")
	}
	ref := firebaseClient.NewRef("students/" + googleID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating student: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteStudent(currentUser User, googleID string) error {
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can delete students")
	}
	ref := firebaseClient.NewRef("students/" + googleID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting student: %v", err)
	}
	return ref.Delete(context.TODO())
}

// Instructor CRUD
func createInstructor(currentUser User, googleID string, instructor Instructor) error {
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can create instructors")
	}
	ref := firebaseClient.NewRef("instructors/" + googleID)
	if err := ref.Set(context.TODO(), instructor); err != nil {
		return fmt.Errorf("error creating instructor: %v", err)
	}
	return ref.Set(context.TODO(), instructor)
}

func readInstructor(currentUser User, googleID string) (Instructor, error) {
	if !isAdmin(currentUser) && !isSelf(currentUser, googleID) {
		return Instructor{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}
	ref := firebaseClient.NewRef("instructors/" + googleID)
	var instructor Instructor
	if err := ref.Get(context.TODO(), &instructor); err != nil {
		return Instructor{}, fmt.Errorf("error reading instructor: %v", err)
	}
	return instructor, nil
}

func updateInstructor(currentUser User, googleID string, updates map[string]interface{}) error {
	if !isAdmin(currentUser) && !isSelf(currentUser, googleID) {
		return fmt.Errorf("unauthorized access: you can only update your own details")
	}
	ref := firebaseClient.NewRef("instructors/" + googleID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating instructor: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteInstructor(currentUser User, googleID string) error {
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can delete instructors")
	}
	ref := firebaseClient.NewRef("instructors/" + googleID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting instructor: %v", err)
	}
	return ref.Delete(context.TODO())
}

// Admin CRUD
func createAdmin(currentUser User, googleID string, admin Admin) error {
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can create admins")
	}
	ref := firebaseClient.NewRef("admins/" + googleID)
	if err := ref.Set(context.TODO(), admin); err != nil {
		return fmt.Errorf("error creating admin: %v", err)
	}
	return ref.Set(context.TODO(), admin)
}

func readAdmin(currentUser User, googleID string) (Admin, error) {
	if !isAdmin(currentUser) && !isSelf(currentUser, googleID) {
		return Admin{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}
	ref := firebaseClient.NewRef("admins/" + googleID)
	var admin Admin
	if err := ref.Get(context.TODO(), &admin); err != nil {
		return Admin{}, fmt.Errorf("error reading admin: %v", err)
	}
	return admin, nil
}

func updateAdmin(currentUser User, googleID string, updates map[string]interface{}) error {
	if !isAdmin(currentUser) && !isSelf(currentUser, googleID) {
		return fmt.Errorf("unauthorized access: you can only update your own details")
	}
	ref := firebaseClient.NewRef("admins/" + googleID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating admin: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteAdmin(currentUser User, googleID string) error {
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can delete admins")
	}
	ref := firebaseClient.NewRef("admins/" + googleID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting admin: %v", err)
	}
	return ref.Delete(context.TODO())
}

// Parent CRUD
func createParent(currentUser User, googleID string, parent Parent) error {
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can create parents")
	}
	ref := firebaseClient.NewRef("parents/" + googleID)
	if err := ref.Set(context.TODO(), parent); err != nil {
		return fmt.Errorf("error creating parent: %v", err)
	}
	return ref.Set(context.TODO(), parent)
}

func readParent(currentUser User, googleID string) (Parent, error) {
	if !isAdmin(currentUser) && !isSelf(currentUser, googleID) {
		return Parent{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}
	ref := firebaseClient.NewRef("parents/" + googleID)
	var parent Parent
	if err := ref.Get(context.TODO(), &parent); err != nil {
		return Parent{}, fmt.Errorf("error reading parent: %v", err)
	}
	return parent, nil
}

func updateParent(currentUser User, googleID string, updates map[string]interface{}) error {
	if !isAdmin(currentUser) && !isSelf(currentUser, googleID) {
		return fmt.Errorf("unauthorized access: you can only update your own details")
	}
	ref := firebaseClient.NewRef("parents/" + googleID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating parent: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteParent(currentUser User, googleID string) error {
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can delete parents")
	}
	ref := firebaseClient.NewRef("parents/" + googleID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting parent: %v", err)
	}
	return ref.Delete(context.TODO())
}
