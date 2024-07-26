package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/gorilla/sessions"
)

var firebaseClient *db.Client

func SessionCookie() (string, error) {
	sessionCookieStore := goDotEnvVariable("SESSION_COOKIE_STORE")
	if sessionCookieStore == "" {
		return sessionCookieStore, fmt.Errorf("SESSION_COOKIE_STORE is not set in the environment variables")
	}
	// sessionCookieStore, found := os.LookupEnv("SESSION_COOKIE_STORE")
	// if !found {
	// 	log.Fatalf("SESSION_COOKIE_STORE is not set in the environment variables")
	// }

	return sessionCookieStore, nil
}

var sessionCookieStore, err = SessionCookie()
var store = sessions.NewCookieStore([]byte(sessionCookieStore))

func initDB(app *firebase.App) error {
	// Initialize Firebase client
	client, err := app.Database(context.Background())
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
func canInstructorAccessStudent(user User, student Student, classes []Class) bool {
	for _, class := range classes {
		if class.Instructor == user.GoogleID && class.ClassID == student.ClassID {
			return true
		}
	}
	return false
}

// Check if parent can access child (student's parent's id = parent id)
func canParentAccessChild(user User, student Student) bool {
	// Implement logic to check if parent can access the child
	return user.GoogleID == student.ParentID
}

// Check if student can access parent (student's parent's name = parent name)
func canChildAccessParent(user User, parent Parent) bool {
	// Implement logic to check if parent can access the child
	return user.GoogleID == parent.GoogleID
}

// Check if student is in the class
func isStudentInClass(user User, students []Student, class Class) bool {
	for _, student := range students {
		if student.GoogleID == user.GoogleID && student.ClassID == class.ClassID {
			return true
		}
	}
	return false
}

// Check if the parent's child is in the class
func isParentChildInClass(user User, students []Student, class Class) bool {
	for _, student := range students {
		if student.ParentID == user.GoogleID && student.ClassID == class.ClassID {
			return true
		}
	}
	return false
}

func getUserRole(email string) (User, string, error) {
	ctx := context.Background()
	var user User
	var userRole string

	// Check if firebaseClient is initialized
	if firebaseClient == nil {
		log.Println("Firebase client is not initialized")
		return user, userRole, fmt.Errorf("firebase client is not initialized")
	}

	// Map categories to Firebase references
	categoryRefs := map[string]string{
		"Student":    "/students",
		"Parent":     "/parents",
		"Instructor": "/instructors",
		"Admin":      "/admins",
	}

	// Iterate through each category and check if the email exists
	for category, ref := range categoryRefs {
		categoryRef := firebaseClient.NewRef(ref)
		dataSnapshot, err := categoryRef.OrderByChild("email").EqualTo(email).GetOrdered(ctx)
		if err != nil {
			log.Printf("Error fetching data from %s: %v", category, err)
			continue
		}

		// Check if dataSnapshot has any children
		if len(dataSnapshot) > 0 {
			userRole = category
			// Assuming dataSnapshot[0] is the first match and it contains the user data
			if err := dataSnapshot[0].Unmarshal(&user); err != nil {
				log.Printf("Error unmarshalling data for %s: %v", category, err)
				continue
			}
			break
		}
	}
	return user, userRole, nil
}

// CRUD operations with role checks

// Student CRUD
func createStudent(student Student, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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

func readStudent(studentGoogleID string, req *http.Request) (Student, error) {
	ref := firebaseClient.NewRef("students/" + studentGoogleID)
	var student Student
	if err := ref.Get(context.TODO(), &student); err != nil {
		return Student{}, fmt.Errorf("error reading student: %v", err)
	}

	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return Student{}, fmt.Errorf("error getting current user")
	}

	classes, err := readAllClasses(req)
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

func updateStudent(studentGoogleID string, updates map[string]interface{}, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
	// Fetch the student information
	student, err := readStudent(studentGoogleID, req)
	if err != nil {
		return fmt.Errorf("error fetching student: %v", err)
	}

	// Fetch all classes
	classes, err := readAllClasses(req)
	if err != nil {
		return fmt.Errorf("error reading classes: %v", err)
	}

	// If user is not admin, instructor, or parent, return error when attempting to update student
	if !isAdmin(currentUser) && //not admin
		!(isSelf(currentUser, studentGoogleID) && isStudent(currentUser)) && //not student and reading self
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

func deleteStudent(student Student, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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
func createInstructor(instructor Instructor, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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

func readInstructor(instructorGoogleID string, req *http.Request) (Instructor, error) {
	ref := firebaseClient.NewRef("instructors/" + instructorGoogleID)
	var instructor Instructor
	if err := ref.Get(context.TODO(), &instructor); err != nil {
		return Instructor{}, fmt.Errorf("error reading instructor: %v", err)
	}

	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return Instructor{}, fmt.Errorf("error getting current user: %v", err)
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

func updateInstructor(instructorGoogleID string, updates map[string]interface{}, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
	// Fetch the instructor information using the provided GoogleID
	instructor, err := readInstructor(instructorGoogleID, req)
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

func deleteInstructor(instructor Instructor, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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
func createAdmin(admin Admin, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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

func readAdmin(adminGoogleID string, req *http.Request) (Admin, error) {
	ref := firebaseClient.NewRef("admins/" + adminGoogleID)
	var admin Admin
	if err := ref.Get(context.TODO(), &admin); err != nil {
		return Admin{}, fmt.Errorf("error reading admin: %v", err)
	}

	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return Admin{}, fmt.Errorf("error getting current user: %v", err)
	}

	// If user is not admin, return error when attempting to read admin
	if !isAdmin(currentUser) {
		return Admin{}, fmt.Errorf("unauthorized access: you can only read your own details")
	}
	return admin, nil
}

func updateAdmin(adminGoogleID string, updates map[string]interface{}, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
	// Fetch the admin information using the provided GoogleID
	admin, err := readAdmin(adminGoogleID, req)
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

func deleteAdmin(admin Admin, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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
func createParent(parent Parent, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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

func readParent(parentGoogleID string, req *http.Request) (Parent, error) {
	ref := firebaseClient.NewRef("parents/" + parentGoogleID)
	var parent Parent
	if err := ref.Get(context.TODO(), &parent); err != nil {
		return Parent{}, fmt.Errorf("error reading parent: %v", err)
	}

	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return Parent{}, fmt.Errorf("error getting current user: %v", err)
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

func updateParent(parentGoogleID string, updates map[string]interface{}, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
	// Fetch the parent information using the provided GoogleID
	parent, err := readParent(parentGoogleID, req)
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

func deleteParent(parent Parent, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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
func createClass(class Class, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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

func readClass(students []Student, classID string, req *http.Request) (Class, error) {
	ref := firebaseClient.NewRef("classes/" + classID)
	var class Class
	if err := ref.Get(context.TODO(), &class); err != nil {
		return Class{}, fmt.Errorf("error reading class: %v", err)
	}

	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return Class{}, fmt.Errorf("error getting current user: %v", err)
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

func readAllClasses(req *http.Request) ([]Class, error) {
	ref := firebaseClient.NewRef("classes")

	var classesMap map[string]Class
	if err := ref.Get(context.TODO(), &classesMap); err != nil {
		return nil, fmt.Errorf("error reading classes: %v", err)
	}

	var classes []Class
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return []Class{}, fmt.Errorf("error getting current user: %v", err)
	}
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

func updateClass(class Class, updates map[string]interface{}, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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

func deleteClass(class Class, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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
func createAnnouncement(announcement Announcement, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
	// If user is not admin, return error when attempting to create announcement
	if !isAdmin(currentUser) {
		return fmt.Errorf("unauthorized access: only admins can create announcements")
	}
	ref := firebaseClient.NewRef("announcements/" + announcement.AnnouncementID)
	if err := ref.Set(context.TODO(), announcement); err != nil {
		return fmt.Errorf("error creating admin: %v", err)
	}
	return ref.Set(context.TODO(), announcement)
}

func readAnnouncement(announcementID string, req *http.Request) (Announcement, error) {
	ref := firebaseClient.NewRef("announcements/" + announcementID)
	var announcement Announcement
	if err := ref.Get(context.TODO(), &announcement); err != nil {
		return Announcement{}, fmt.Errorf("error reading admin: %v", err)
	}

	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return Announcement{}, fmt.Errorf("error getting current user: %v", err)
	}

	// If user is not NK user, return error when attempting to read admin
	if !isAdmin(currentUser) &&
		!isInstructor(currentUser) &&
		!isParent(currentUser) &&
		!isStudent(currentUser) {
		return Announcement{}, fmt.Errorf("unauthorized access: you are not allowed to read this announcement")
	}
	return announcement, nil
}

func updateAnnouncement(announcement Announcement, updates map[string]interface{}, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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

func deleteAnnouncement(announcement Announcement, req *http.Request) error {
	currentUser, err := GetCurrentUser(req)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
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
