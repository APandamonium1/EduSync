package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"encoding/gob"
	"encoding/json"
)

func init() {
	// Register the User type with gob
	gob.Register(User{})
	gob.Register(Announcement{})
}

var (
	// testStore = sessions.NewCookieStore([]byte("a-very-secret-key"))

	currentUser = User{
		GoogleID:      "admin-user",
		Name:          "Admin User",
		Email:         "jiaweicodetest@gmail.com",
		ContactNumber: "12345678",
		Role:          "Admin",
	}

	// Sample data for testing
	students = []Student{
		{
			User: User{
				GoogleID:      "test-student",
				Name:          "John Doe",
				Email:         "jeyvianangjieen@gmail.com",
				ContactNumber: "91234567",
				Role:          "Student",
			},
			Age:           12,
			LessonCredits: 10.0,
			ClassID:       "TE-6-10",
			ParentID:      "test-parent",
		},
	}

	instructor = Instructor{
		User: User{
			GoogleID:      "test-student",
			Name:          "John Doe",
			Email:         "jeyvianang112462@gmail.com",
			ContactNumber: "91234567",
			Role:          "Student",
		},
		BasePay:          15,
		NumberOfStudents: 24,
	}

	admin = Admin{
		User: User{
			GoogleID:      "test-admin",
			Name:          "Awesomeness",
			ContactNumber: "99999999",
			Email:         "awesome_admin@nk.com",
			Role:          "Admin",
		},
		BasePay:   15,
		Incentive: 24,
	}

// Create a new admin
var admin = Admin{
	User: User{
		GoogleID:      "test-admin",
		Name:          "Awesomeness",
		ContactNumber: "99999999",
		Email:         "jeyvianangjieen@gmail.com",
		Role:          "Admin",
	},
	BasePay:   15,
	Incentive: 24,
}

	classes = []Class{
		{
			ClassID: "te-6-10",
			Name:    "Test Class",
		},
	}

	announcement = Announcement{
		AnnouncementID: "test-announcement",
		Subject:        "Test Announcement",
		Content:        "This is a test announcement.",
	}
)

// func mockRequest() *http.Request {
// 	req, _ := http.NewRequest("GET", "/test", nil)
// 	session, _ := store.Get(req, "session-name")
// 	session.Values["user"] = currentUser
// 	err := session.Save(req, httptest.NewRecorder())
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return req
// }

func mockRequest() *http.Request {
	req, _ := http.NewRequest("GET", "/test", nil)
	recorder := httptest.NewRecorder()
	session, _ := store.Get(req, "auth-session")

	userData, err := json.Marshal(currentUser)
	if err != nil {
		log.Println(err)
	}

	session.Values["user"] = userData
	err = session.Save(req, recorder)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Cookie", recorder.Header().Get("Set-Cookie"))
	return req
}

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
	err := initializeFirebase()
	if err != nil {
		t.Fatalf("Error initializing Firebase: %v", err)
	}

	err = createStudent(students[0], mockRequest())
	if err != nil {
		t.Fatalf("Error creating student: %v", err)
	}

	readStudent, err := readStudent(students[0].GoogleID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading student: %v", err)
	}

	if !reflect.DeepEqual(students[0], readStudent) {
		t.Errorf("Created and read students do not match")
	}
}

func TestReadStudent(t *testing.T) {
	readStudent, err := readStudent(students[0].GoogleID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading student: %v", err)
	}

	if readStudent.GoogleID != "test-student" {
		t.Errorf("Expected ID %v, got %v", "test-student", readStudent.GoogleID)
	}
}

func TestUpdateStudent(t *testing.T) {
	// Update the student's email
	updates := map[string]interface{}{
		"name":           "John Doe",
		"contact_number": "99999999",
	}

	err := updateStudent(students[0].GoogleID, updates, mockRequest())
	if err != nil {
		t.Fatalf("Error updating student: %v", err)
	}

	// Read the updated student
	readStudent, err := readStudent(students[0].GoogleID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading student after updating: %v", err)
	}

	// Assert that the updated student's email is correct
	if students[0].Name != updates["name"] {
		t.Errorf("Updated student's name is incorrect. Expected: %v, Got: %v", updates["name"], readStudent.Name)
	}
}

func TestDeleteStudent(t *testing.T) {
	// Delete the student
	err := deleteStudent(students[0], mockRequest())
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
	err := createInstructor(instructor, mockRequest())
	if err != nil {
		t.Fatalf("Error creating instructor: %v", err)
	}

	// Read the created instructor
	readInstructor, err := readInstructor(instructor.GoogleID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading instructor: %v", err)
	}

	// Assert that the created and read instructor are equal
	if !reflect.DeepEqual(instructor, readInstructor) {
		t.Error("Created and read instructors are not equal")
	}
}

func TestReadInstructor(t *testing.T) {

	instructor, err := readInstructor(instructor.GoogleID, mockRequest())
	if err != nil {
		t.Fatalf("Failed to read instructor: %v", err)
	}

	if instructor.GoogleID != "test-instructor" {
		t.Fatalf("Expected GoogleID %v, got %v", "test-instructor", instructor.GoogleID)
	}
}

func TestUpdateInstructor(t *testing.T) {
	// Update the instructor's email
	updates := map[string]interface{}{
		"email": "amazing_instructor@nk.com",
	}

	err := updateInstructor(instructor.GoogleID, updates, mockRequest())
	if err != nil {
		t.Fatalf("Error updating instructor: %v", err)
	}

	// Read the updated instructor
	readInstructor, err := readInstructor(instructor.GoogleID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading instructor: %v", err)
	}

	// Assert that the updated instructor's email is correct
	if readInstructor.Email != updates["email"] {
		t.Errorf("Updated instructor's email is incorrect. Expected: %v, Got: %v", updates["email"], readInstructor.Email)
	}
}

func TestDeleteInstructor(t *testing.T) {
	// Delete the instructor
	err := deleteInstructor(instructor, mockRequest())
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
	err := createAdmin(admin, mockRequest())
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}

	// Read the created admin
	readAdmin, err := readAdmin(admin.GoogleID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading admin: %v", err)
	}

	// Assert that the created and read admin are equal
	if !reflect.DeepEqual(admin, readAdmin) {
		t.Error("Created and read admins are not equal")
	}
}

func TestReadAdmin(t *testing.T) {
	admin, err := readAdmin(admin.GoogleID, mockRequest())
	if err != nil {
		t.Fatalf("Failed to read instructor: %v", err)
	}

	if admin.GoogleID != "test-admin" {
		t.Fatalf("Expected GoogleID %v, got %v", "test-admin", admin.GoogleID)
	}
}

func TestUpdateAdmin(t *testing.T) {
	// Update the admin's email
	updates := map[string]interface{}{
		"email": "amazing_admin@nk.com",
	}

	err := updateAdmin(admin.GoogleID, updates, mockRequest())
	if err != nil {
		t.Fatalf("Error updating admin: %v", err)
	}

	// Read the updated admin
	readAdmin, err := readAdmin(admin.GoogleID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading admin: %v", err)
	}

	// Assert that the updated admin's email is correct
	if readAdmin.Email != updates["email"] {
		t.Errorf("Updated admin's email is incorrect. Expected: %v, Got: %v", updates["email"], readAdmin.Email)
	}
}

func TestDeleteAdmin(t *testing.T) {
	// Delete the admin
	err := deleteAdmin(admin, mockRequest())
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
	err := createParent(parent, mockRequest())
	if err != nil {
		t.Fatalf("Error creating parent: %v", err)
	}

	// Read the created parent
	readParent, err := readParent(parent.GoogleID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading parent: %v", err)
	}

	// Assert that the created and read parent are equal
	if !reflect.DeepEqual(parent, readParent) {
		t.Error("Created and read parents are not equal")
	}
}

func TestReadParent(t *testing.T) {
	parent, err := readParent(parent.GoogleID, mockRequest())
	if err != nil {
		t.Fatalf("Failed to read parent: %v", err)
	}

	if parent.GoogleID != "test-parent" {
		t.Fatalf("Expected GoogleID %v, got %v", "test-parent", parent.GoogleID)
	}
}

func TestUpdateParent(t *testing.T) {
	// Update the parent's email
	updates := map[string]interface{}{
		"email": "jane_doe_parent@nk.com",
	}

	err := updateParent(parent.GoogleID, updates, mockRequest())
	if err != nil {
		t.Fatalf("Error updating parent: %v", err)
	}

	// Read the updated parent
	readParent, err := readParent(parent.GoogleID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading parent: %v", err)
	}

	// Assert that the updated parent's email is correct
	if readParent.Email != updates["email"] {
		t.Errorf("Updated parent's email is incorrect. Expected: %v, Got: %v", updates["email"], readParent.Email)
	}
}

func TestDeleteParent(t *testing.T) {
	// Delete the parent
	err := deleteParent(parent, mockRequest())
	if err != nil {
		t.Fatalf("Error deleting parent: %v", err)
	}

	// Try to read the deleted parent
	// _, err = readParent(googleID)
	// if err == nil {
	// 	t.Error("Deleted parent still exists")
	// }
}

// Testing for class CRUD operations

func TestCreateClass(t *testing.T) {
	err := createClass(classes[0], mockRequest())
	if err != nil {
		t.Fatalf("Error creating class: %v", err)
	}

	// Read the created class
	readClass, err := readClass(students, classes[0].ClassID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading class: %v", err)
	}

	// Assert that the created and read class are equal
	if !reflect.DeepEqual(classes[0], readClass) {
		t.Error("Created and read classes are not equal")
	}
}

func TestReadClass(t *testing.T) {
	class, err := readClass(students, classes[0].ClassID, mockRequest())
	if err != nil {
		t.Fatalf("Failed to read class: %v", err)
	}

	if class.ClassID != classes[0].ClassID {
		t.Fatalf("Expected ID %v, got %v", classes[0].ClassID, class.ClassID)
	}
}

func TestUpdateClass(t *testing.T) {
	// Update the class's name
	updates := map[string]interface{}{
		"class_name": "DN",
	}

	err := updateClass(classes[0], updates, mockRequest())
	if err != nil {
		t.Fatalf("Error updating class: %v", err)
	}

	// Read the updated class
	readClass, err := readClass(students, classes[0].ClassID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading class: %v", err)
	}

	// Assert that the updated class's name is correct
	if readClass.Name != updates["class_name"] {
		t.Errorf("Updated class's name is incorrect. Expected: %v, Got: %v", updates["class_name"], readClass.Name)
	}
}

func TestDeleteClass(t *testing.T) {
	// Delete the class
	err := deleteClass(classes[0], mockRequest())
	if err != nil {
		t.Fatalf("Error deleting class: %v", err)
	}

	// Try to read the deleted class
	// _, err = readClass(currentUser, students, classes[0])
	// if err == nil {
	// 	t.Error("Deleted class still exists")
	// }
}

func TestCreateAnnouncement(t *testing.T) {
	err := createAnnouncement(announcement, mockRequest())
	if err != nil {
		t.Fatalf("Error creating announcement: %v", err)
	}

	// Read the announcement
	readAnnouncement, err := readAnnouncement(announcement.AnnouncementID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading announcement: %v", err)
	}

	// Assert that the created and read announcement are equal
	if !reflect.DeepEqual(announcement, readAnnouncement) {
		t.Errorf("Created and read announcements are not equal. Expected: %+v, Got: %+v", announcement, readAnnouncement)
	}
}

func TestReadAnnouncement(t *testing.T) {
	announcement, err := readAnnouncement(announcement.AnnouncementID, mockRequest())
	if err != nil {
		t.Fatalf("Failed to read announcement: %v", err)
	}

	if announcement.Subject != "Test Announcement" {
		t.Fatalf("Expected Title %v, got %v", "Test Announcement", announcement.Subject)
	}
}

func TestUpdateAnnouncement(t *testing.T) {
	// Update the announcement content
	updates := map[string]interface{}{
		"content": "This is an updated announcement.",
	}

	err := updateAnnouncement(announcement, updates, mockRequest())
	if err != nil {
		t.Fatalf("Error updating announcement: %v", err)
	}

	// Read the updated announcement
	readAnnouncement, err := readAnnouncement(announcement.AnnouncementID, mockRequest())
	if err != nil {
		t.Fatalf("Error reading announcement: %v", err)
	}

	// Assert that the updated announcement's content is correct
	if readAnnouncement.Content != updates["content"] {
		t.Errorf("Updated announcement's content is incorrect. Expected: %v, Got: %v", updates["content"], readAnnouncement.Content)
	}
}

func TestDeleteAnnouncement(t *testing.T) {
	// Delete the announcement
	err := deleteAnnouncement(announcement, mockRequest())
	if err != nil {
		t.Fatalf("Error deleting announcement: %v", err)
	}

	// Try to read the deleted announcement
	// _, err = readAnnouncement(currentUser, announcement)
	// if err == nil {
	// 	t.Error("Deleted announcement still exists")
	// }
}
