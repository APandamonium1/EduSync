package main

import (
	"time"
	// "github.com/google/uuid"
)

// Student struct for storing student information
type Student struct {
	// ID            uuid.UUID `json:"id"`
	GoogleID      string    `json:"google_id"`
	Name          string    `json:"name"`
	Age           int       `json:"age"`
	LessonCredits float32   `json:"lesson_credits"`
	Email         string    `json:"email"`
	ContactNumber string    `json:"contact_number"`
	Class         string    `json:"class"`
	Instructor    string    `json:"instructor"`
	ParentName    string    `json:"parent_name"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Instructor struct for storing instructor information
type Instructor struct {
	// ID               uuid.UUID `json:"id"`
	GoogleID         string    `json:"google_id"`
	Name             string    `json:"name"`
	ContactNumber    string    `json:"contact_number"`
	Email            string    `json:"email"`
	BasePay          float64   `json:"base_pay"`
	NumberOfStudents int       `json:"number_of_students"`
	Role             string    `json:"role"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Parent struct for storing parent information
type Parent struct {
	// ID        uuid.UUID `json:"id"`
	GoogleID      string    `json:"google_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	ContactNumber string    `json:"contact_no"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewStudent creates a new Student instance
// func NewStudent(name string, age int, lessonCredits float32, email string, contactNumber string, class string, instructor string, parentName string) Student {
func NewStudent(googleID string, name string, age int, lessonCredits float32, email, contactNumber, class, instructor, parentName, role string) Student {
	return Student{
		// ID:            uuid.New(),
		GoogleID:      googleID,
		Name:          name,
		Age:           age,
		LessonCredits: lessonCredits,
		Email:         email,
		ContactNumber: contactNumber,
		Class:         class,
		Instructor:    instructor,
		ParentName:    parentName,
		Role:          role,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// NewInstructor creates a new Instructor instance
// func NewInstructor(name string, contactNumber string, email string, basePay float64, numberOfStudents int) Instructor {
func NewInstructor(googleID, name, contactNumber, email, role string, basePay float64, numberOfStudents int) Instructor {
	return Instructor{
		// ID:               uuid.New(),
		GoogleID:         googleID,
		Name:             name,
		ContactNumber:    contactNumber,
		Email:            email,
		BasePay:          basePay,
		NumberOfStudents: numberOfStudents,
		Role:             role,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// func NewParent(name string, email string, contactNo string) Parent {
func NewParent(googleID, name, email, contactNumber, role string) Parent {
	return Parent{
		// ID:        uuid.New(),
		GoogleID:      googleID,
		Name:          name,
		Email:         email,
		ContactNumber: contactNumber,
		Role:          role,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
