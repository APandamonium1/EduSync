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

// Utility function to get current user
func getCurrentUser(req *http.Request) (User, error) {
	// Implement session or context based user retrieval
	return User{}, nil
}

// Utility functions to check roles
func isAdmin(user User) bool {
	return user.Role == "Admin"
}

func isInstructor(user User) bool {
	return user.Role == "Instructor"
}

func isParent(user User) bool {
	return user.Role == "Parent"
}

func isStudent(user User) bool {
	return user.Role == "Student"
}

// Function to check if CurrentUser is checking their own details
func isSelf(user User, googleID string) bool {
	return user.GoogleID == googleID
}

// Check if instructor can access student (student's class' instructor = instructor name)
func canInstructorAccessStudent(currentUser User, student Student, classes []Class) bool {
	for _, class := range classes {
		if class.Instructor == currentUser.GoogleID && class.ClassID == student.ClassID {
			return true
		}
	}
	return false
}

// Check if parent can access child (student's parent's id = parent id)
func canParentAccessChild(currentUser User, student Student) bool {
	// Implement logic to check if parent can access the child
	return currentUser.GoogleID == student.ParentID
}

// Check if student can access parent (student's parent's name = parent name)
func canChildAccessParent(currentUser User, parent Parent) bool {
	// Implement logic to check if parent can access the child
	return currentUser.GoogleID == parent.GoogleID
}

// Check if student is in the class
func isStudentInClass(currentUser User, students []Student, class Class) bool {
	for _, student := range students {
		if student.GoogleID == currentUser.GoogleID && student.ClassID == class.ClassID {
			return true
		}
	}
	return false
}

// Check if the parent's child is in the class
func isParentChildInClass(currentUser User, students []Student, class Class) bool {
	for _, student := range students {
		if student.ParentID == currentUser.GoogleID && student.ClassID == class.ClassID {
			return true
		}
	}
	return false
}

// CRUD operations with role checks

// Student CRUD
func createStudent(currentUser User, student Student) error {
	// If user is not an admin, return error when attempting to create student
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can create students")
	}
	ref := firebaseClient.NewRef("students/" + student.GoogleID)
	if err := ref.Set(context.TODO(), student); err != nil {
		return fmt.Errorf("error creating student: %v", err)
	}
	return ref.Set(context.TODO(), student)
}

func readStudent(currentUser User, student Student, classes []Class) (Student, error) {
	// If user is not an admin, instructor, or parent, return error when attempting to read student
	if !isAdmin(currentUser) && //not admin
		!(isSelf(currentUser, student.GoogleID) && isStudent(currentUser)) && //not student and reading self
		!(currentUser.Role == "Instructor" && canInstructorAccessStudent(currentUser, student, classes)) && //instructor can access only their students' info
		!(currentUser.Role == "Parent" && canParentAccessChild(currentUser, student)) { // parent can access only their child's info
		return Student{}, fmt.Errorf("unauthorized access: you can only read your own details or the details of students you are authorized to access")
	}

	ref := firebaseClient.NewRef("students/" + student.GoogleID)
	// var student Student
	if err := ref.Get(context.TODO(), &student); err != nil {
		return Student{}, fmt.Errorf("error reading student: %v", err)
	}
	return student, nil
}

func updateStudent(currentUser User, student Student, classes []Class, updates map[string]interface{}) error {
	// If user is not admin, instructor, or parent, return error when attempting to update student
	if !isAdmin(currentUser) && //not admin
		!(isSelf(currentUser, student.GoogleID) && isStudent(currentUser)) && //not student and reading self
		!(currentUser.Role == "Instructor" && canInstructorAccessStudent(currentUser, student, classes)) && //instructor can access only their students' info
		!(currentUser.Role == "Parent" && canParentAccessChild(currentUser, student)) { // parent can access only their child's info {
		return fmt.Errorf("unauthorized access: you can only update your own details")
	}
	ref := firebaseClient.NewRef("students/" + student.GoogleID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating student: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteStudent(currentUser User, student Student) error {
	// If user is not admin, return error when attempting to delete student
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can delete students")
	}
	ref := firebaseClient.NewRef("students/" + student.GoogleID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting student: %v", err)
	}
	return ref.Delete(context.TODO())
}

// Instructor CRUD
func createInstructor(currentUser User, instructor Instructor) error {
	// If user is not admin, return error when attempting to create instructor
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can create instructors")
	}
	ref := firebaseClient.NewRef("instructors/" + instructor.GoogleID)
	if err := ref.Set(context.TODO(), instructor); err != nil {
		return fmt.Errorf("error creating instructor: %v", err)
	}
	return ref.Set(context.TODO(), instructor)
}

func readInstructor(currentUser User, instructor Instructor) (Instructor, error) {
	// If user is not admin or (self & instructor), return error when attempting to read instructor
	if !isAdmin(currentUser) && //not admin
		!(isSelf(currentUser, instructor.GoogleID) && isInstructor(currentUser)) {
		return Instructor{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}
	ref := firebaseClient.NewRef("instructors/" + instructor.GoogleID)
	if err := ref.Get(context.TODO(), &instructor); err != nil {
		return Instructor{}, fmt.Errorf("error reading instructor: %v", err)
	}
	return instructor, nil
}

func updateInstructor(currentUser User, instructor Instructor, updates map[string]interface{}) error {
	// If user is not admin or (self & instructor), return error when attempting to update instructor
	if !isAdmin(currentUser) && //not admin
		!(isSelf(currentUser, instructor.GoogleID) && isInstructor(currentUser)) {
		return fmt.Errorf("unauthorized access: you can only update your own details")
	}
	ref := firebaseClient.NewRef("instructors/" + instructor.GoogleID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating instructor: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteInstructor(currentUser User, instructor Instructor) error {
	// If user is not admin, return error when attempting to delete instructor
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can delete instructors")
	}
	ref := firebaseClient.NewRef("instructors/" + instructor.GoogleID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting instructor: %v", err)
	}
	return ref.Delete(context.TODO())
}

// Admin CRUD
func createAdmin(currentUser User, admin Admin) error {
	// If user is not admin, return error when attempting to create admin
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can create admins")
	}
	ref := firebaseClient.NewRef("admins/" + admin.GoogleID)
	if err := ref.Set(context.TODO(), admin); err != nil {
		return fmt.Errorf("error creating admin: %v", err)
	}
	return ref.Set(context.TODO(), admin)
}

func readAdmin(currentUser User, admin Admin) (Admin, error) {
	// If user is not admin, return error when attempting to read admin
	if !isAdmin(currentUser) {
		return Admin{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}
	ref := firebaseClient.NewRef("admins/" + admin.GoogleID)
	if err := ref.Get(context.TODO(), &admin); err != nil {
		return Admin{}, fmt.Errorf("error reading admin: %v", err)
	}
	return admin, nil
}

func updateAdmin(currentUser User, admin Admin, updates map[string]interface{}) error {
	// If user is not admin, return error when attempting to update admin
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: you can only update your own details")
	}
	ref := firebaseClient.NewRef("admins/" + admin.GoogleID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating admin: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteAdmin(currentUser User, admin Admin) error {
	// If user is not admin, return error when attempting to delete admin
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can delete admins")
	}
	ref := firebaseClient.NewRef("admins/" + admin.GoogleID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting admin: %v", err)
	}
	return ref.Delete(context.TODO())
}

// Parent CRUD
func createParent(currentUser User, parent Parent) error {
	// If user is not admin, return error when attempting to create parent
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can create parents")
	}
	ref := firebaseClient.NewRef("parents/" + parent.GoogleID)
	if err := ref.Set(context.TODO(), parent); err != nil {
		return fmt.Errorf("error creating parent: %v", err)
	}
	return ref.Set(context.TODO(), parent)
}

func readParent(currentUser User, parent Parent) (Parent, error) {
	// If user is not admin or (self and parent), return error when attempting to update parent
	if !isAdmin(currentUser) && //not admin
		!(isSelf(currentUser, parent.GoogleID) && isParent(currentUser)) && //not parent and reading self
		!(currentUser.Role == "Student" && canChildAccessParent(currentUser, parent)) {
		return Parent{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}
	ref := firebaseClient.NewRef("parents/" + parent.GoogleID)
	if err := ref.Get(context.TODO(), &parent); err != nil {
		return Parent{}, fmt.Errorf("error reading parent: %v", err)
	}
	return parent, nil
}

func updateParent(currentUser User, parent Parent, updates map[string]interface{}) error {
	if !isAdmin(currentUser) && //not admin
		!(isSelf(currentUser, parent.GoogleID) && isParent(currentUser)) {
		return fmt.Errorf("unauthorized access: you can only update your own details")
	}
	ref := firebaseClient.NewRef("parents/" + parent.GoogleID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating parent: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteParent(currentUser User, parent Parent) error {
	// If user is not admin, return error when attempting to delete parent
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can delete parents")
	}
	ref := firebaseClient.NewRef("parents/" + parent.GoogleID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting parent: %v", err)
	}
	return ref.Delete(context.TODO())
}

// class CRUD

func createClass(currentUser User, class Class) error {
	// If user is not admin, return error when attempting to create class
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can create classes")
	}
	ref := firebaseClient.NewRef("classes/" + class.ClassID)
	if err := ref.Set(context.TODO(), class); err != nil {
		return fmt.Errorf("error creating class: %v", err)
	}
	return ref.Set(context.TODO(), class)
}

func readClass(currentUser User, students []Student, class Class) (Class, error) {
	// If user is not admin or (self and class), return error when attempting to read class
	if !isAdmin(currentUser) && //not admin
		!isInstructor(currentUser) && //not instructor
		!(isStudent(currentUser) && isStudentInClass(currentUser, students, class)) &&
		!(isParent(currentUser) && isParentChildInClass(currentUser, students, class)) {
		return Class{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}
	ref := firebaseClient.NewRef("classes/" + class.ClassID)
	if err := ref.Get(context.TODO(), &class); err != nil {
		return Class{}, fmt.Errorf("error reading class: %v", err)
	}
	return class, nil
}

func updateClass(currentUser User, class Class, updates map[string]interface{}) error {
	// If user is not admin or (self and class), return error when attempting to update class
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: you can only update your own details")
	}
	ref := firebaseClient.NewRef("classes/" + class.ClassID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating class: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteClass(currentUser User, class Class) error {
	// If user is not admin, return error when attempting to delete class
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can delete classes")
	}
	ref := firebaseClient.NewRef("classes/" + class.ClassID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting class: %v", err)
	}
	return ref.Delete(context.TODO())
}

// Announcements CRUD
func createAnnouncement(currentUser User, announcement Announcement) error {
	// If user is not admin, return error when attempting to create announcement
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can create announcements")
	}
	ref := firebaseClient.NewRef("admins/" + announcement.AnnouncementID)
	if err := ref.Set(context.TODO(), announcement); err != nil {
		return fmt.Errorf("error creating admin: %v", err)
	}
	return ref.Set(context.TODO(), announcement)
}

func readAnnouncement(currentUser User, announcement Announcement) (Announcement, error) {
	// If user is not admin, return error when attempting to read admin
	if !isAdmin(currentUser) &&
		!isInstructor(currentUser) &&
		!isParent(currentUser) &&
		!isStudent(currentUser) {
		return Announcement{}, fmt.Errorf("unauthorized access: you are not allowed to read this announcement")
	}
	ref := firebaseClient.NewRef("announcements/" + announcement.AnnouncementID)
	if err := ref.Get(context.TODO(), &announcement); err != nil {
		return Announcement{}, fmt.Errorf("error reading admin: %v", err)
	}
	return announcement, nil
}

func updateAnnouncement(currentUser User, announcement Announcement, updates map[string]interface{}) error {
	// If user is not admin, return error when attempting to update announcement
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can update this announcement")
	}
	ref := firebaseClient.NewRef("announcements/" + announcement.AnnouncementID)
	if err := ref.Update(context.TODO(), updates); err != nil {
		return fmt.Errorf("error updating announcement: %v", err)
	}
	return ref.Update(context.TODO(), updates)
}

func deleteAnnouncement(currentUser User, announcement Announcement) error {
	// If user is not admin, return error when attempting to delete announcement
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can delete announcements")
	}
	ref := firebaseClient.NewRef("announcements/" + announcement.AnnouncementID)
	if err := ref.Delete(context.TODO()); err != nil {
		return fmt.Errorf("error deleting announcement: %v", err)
	}
	return ref.Delete(context.TODO())
}
