package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
var (
	firebaseClient *db.Client
	firebaseApp    *firebase.App
)

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

	//opt := option.WithCredentialsFile("edusync-7bd5e-firebase-adminsdk-x49uh-af084a6314.json")
	opt := option.WithCredentialsFile("edusync-test-firebase-adminsdk-hk5kl-9af0162b09.json")

	firebaseApp, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return fmt.Errorf("error initializing firebase app: %v", err)
	}

	client, err := firebaseApp.Database(ctx)
	if err != nil {
		return fmt.Errorf("error creating firebase DB client: %v", err)
	}

	firebaseClient = client
	return nil
}

// Utility function to get current user
func GetCurrentUser(req *http.Request) (User, error) {
	session, err := store.Get(req, "auth-session")
	if err != nil {
		return User{}, err
	}

	userData, ok := session.Values["user"].([]byte)
	if !ok || userData == nil {
		return User{}, fmt.Errorf("user not found in session")
	}

	var user User
	err = json.Unmarshal(userData, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
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

func readStudent(currentUser User, studentGoogleID string) (Student, error) {
	ref := firebaseClient.NewRef("students/" + studentGoogleID)
	var student Student
	if err := ref.Get(context.TODO(), &student); err != nil {
		return Student{}, fmt.Errorf("error reading student: %v", err)
	}

	classes, err := readAllClasses(currentUser)
	if err != nil {
		return Student{}, fmt.Errorf("error reading classes: %v", err)
	}

	// If user is not an admin, instructor, or parent, return error when attempting to read student
	if !isAdmin(currentUser) && //not admin
		!(isSelf(currentUser, student.GoogleID) && isStudent(currentUser)) && //not student and reading self
		!(currentUser.Role == "Instructor" && canInstructorAccessStudent(currentUser, student, classes)) && //instructor can access only their students' info
		!(currentUser.Role == "Parent" && canParentAccessChild(currentUser, student)) { // parent can access only their child's info
		return Student{}, fmt.Errorf("unauthorized access: you can only read your own details or the details of students you are authorized to access")
	}

	return student, nil
}

func readStudents() ([]Student, error) {
	var studentsMap map[string]Student
	ref := firebaseClient.NewRef("students")
	if err := ref.Get(context.TODO(), &studentsMap); err != nil {
		return nil, fmt.Errorf("error reading students: %v", err)
	}
	// Convert map to slice
	students := make([]Student, 0, len(studentsMap))
	for _, student := range studentsMap {
		students = append(students, student)
	}
	return students, nil
}

func searchStudents(name, class string) ([]Student, error) {
	if name == "" && class == "" {
		return readStudents()
	}
	students, err := readStudents()
	if err != nil {
		return nil, err
	}
	var filteredStudents []Student
	for _, student := range students {
		if name == "" || strings.Contains(strings.ToLower(student.Name), strings.ToLower(name)) {
			filteredStudents = append(filteredStudents, student)
		}
	}
	return filteredStudents, nil
}

func updateStudent(currentUser User, studentGoogleID string, updates map[string]interface{}) error {
	// Fetch the student information using the provided GoogleID
	student, err := readStudent(currentUser, studentGoogleID)
	if err != nil {
		return fmt.Errorf("error fetching student: %v", err)
	}

	// Fetch all classes
	classes, err := readAllClasses(currentUser)
	if err != nil {
		return fmt.Errorf("error reading classes: %v", err)
	}

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
	return nil
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

func readInstructor(currentUser User, instructorGoogleID string) (Instructor, error) {
	ref := firebaseClient.NewRef("instructors/" + instructorGoogleID)
	var instructor Instructor
	if err := ref.Get(context.TODO(), &instructor); err != nil {
		return Instructor{}, fmt.Errorf("error reading instructor: %v", err)
	}

	// If user is not admin or (self & instructor), return error when attempting to read instructor
	if !isAdmin(currentUser) && //not admin
		!(isSelf(currentUser, instructor.GoogleID) && isInstructor(currentUser)) {
		return Instructor{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}

	return instructor, nil
}

func readInstructors() ([]Instructor, error) {
	var instructorsMap map[string]Instructor
	ref := firebaseClient.NewRef("instructors")
	if err := ref.Get(context.TODO(), &instructorsMap); err != nil {
		return nil, fmt.Errorf("error reading students: %v", err)
	}
	// Convert map to slice
	instructors := make([]Instructor, 0, len(instructorsMap))
	for _, instructor := range instructorsMap {
		instructors = append(instructors, instructor)
	}
	return instructors, nil
}

func searchInstructors(name string) ([]Instructor, error) {
	if name == "" {
		return readInstructors()
	}
	instructors, err := readInstructors()
	if err != nil {
		return nil, err
	}
	var filteredInstructors []Instructor
	for _, instructor := range instructors {
		if name == "" || strings.Contains(strings.ToLower(instructor.Name), strings.ToLower(name)) {
			filteredInstructors = append(filteredInstructors, instructor)
		}
	}
	return filteredInstructors, nil
}

func updateInstructor(currentUser User, instructorGoogleID string, updates map[string]interface{}) error {
	// Fetch the instructor information using the provided GoogleID
	instructor, err := readInstructor(currentUser, instructorGoogleID)
	if err != nil {
		return fmt.Errorf("error fetching instructor: %v", err)
	}

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

func readAdmin(currentUser User, adminGoogleID string) (Admin, error) {
	ref := firebaseClient.NewRef("admins/" + adminGoogleID)
	var admin Admin
	if err := ref.Get(context.TODO(), &admin); err != nil {
		return Admin{}, fmt.Errorf("error reading admin: %v", err)
	}

	// If user is not admin, return error when attempting to read admin
	if !isAdmin(currentUser) {
		return Admin{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}
	return admin, nil
}

func updateAdmin(currentUser User, adminGoogleID string, updates map[string]interface{}) error {
	// Fetch the admin information using the provided GoogleID
	admin, err := readAdmin(currentUser, adminGoogleID)
	if err != nil {
		return fmt.Errorf("error fetching admin: %v", err)
	}

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

func readParent(currentUser User, parentGoogleID string) (Parent, error) {
	ref := firebaseClient.NewRef("parents/" + parentGoogleID)
	var parent Parent
	if err := ref.Get(context.TODO(), &parent); err != nil {
		return Parent{}, fmt.Errorf("error reading parent: %v", err)
	}

	// If user is not admin or (self and parent), return error when attempting to update parent
	if !isAdmin(currentUser) && //not admin
		!(isSelf(currentUser, parent.GoogleID) && isParent(currentUser)) && //not parent and reading self
		!(currentUser.Role == "Student" && canChildAccessParent(currentUser, parent)) {
		return Parent{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}
	return parent, nil
}

func readParents() ([]Parent, error) {
	var parentsMap map[string]Parent
	ref := firebaseClient.NewRef("parents")
	if err := ref.Get(context.TODO(), &parentsMap); err != nil {
		return nil, fmt.Errorf("error reading students: %v", err)
	}
	// Convert map to slice
	parents := make([]Parent, 0, len(parentsMap))
	for _, parent := range parentsMap {
		parents = append(parents, parent)
	}
	return parents, nil
}

func searchParents(name string) ([]Parent, error) {
	if name == "" {
		return readParents()
	}
	parents, err := readParents()
	if err != nil {
		return nil, err
	}
	var filteredParents []Parent
	for _, parent := range parents {
		if name == "" || strings.Contains(strings.ToLower(parent.Name), strings.ToLower(name)) {
			filteredParents = append(filteredParents, parent)
		}
	}
	return filteredParents, nil
}

func updateParent(currentUser User, parentGoogleID string, updates map[string]interface{}) error {
	// Fetch the parent information using the provided GoogleID
	parent, err := readParent(currentUser, parentGoogleID)
	if err != nil {
		return fmt.Errorf("error fetching parent: %v", err)
	}

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

func readClass(currentUser User, students []Student, classID string) (Class, error) {
	ref := firebaseClient.NewRef("classes/" + classID)
	var class Class
	if err := ref.Get(context.TODO(), &class); err != nil {
		return Class{}, fmt.Errorf("error reading class: %v", err)
	}

	// If user is not admin or (self and class), return error when attempting to read class
	if !isAdmin(currentUser) && //not admin
		!isInstructor(currentUser) && //not instructor
		!(isStudent(currentUser) && isStudentInClass(currentUser, students, class)) &&
		!(isParent(currentUser) && isParentChildInClass(currentUser, students, class)) {
		return Class{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}
	return class, nil
}

func readAllClasses(currentUser User) ([]Class, error) {
	ref := firebaseClient.NewRef("classes")

	var classesMap map[string]Class
	if err := ref.Get(context.TODO(), &classesMap); err != nil {
		return nil, fmt.Errorf("error reading classes: %v", err)
	}

	var classes []Class
	for _, class := range classesMap {
		// If user is not authorized to read the class, skip it
		if !isAdmin(currentUser) && // not admin
			!isInstructor(currentUser) { // not instructor
			continue
		}
		classes = append(classes, class)
	}

	return classes, nil
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
