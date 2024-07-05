package main

import (
	"time"
	// "github.com/google/uuid"
)

// User struct to represent the base user with common fields
type User struct {
	// ID            uuid.UUID `json:"id"`
	GoogleID      string    `json:"google_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	ContactNumber string    `json:"contact_number"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Student struct for storing student information
type Student struct {
	User
	Age           int     `json:"age"`
	LessonCredits float32 `json:"lesson_credits"`
	ClassID       string  `json:"class_id"`
	// Instructor    string  `json:"instructor"`
	ParentName string `json:"parent_name"`
}

// Instructor struct for storing instructor information
type Instructor struct {
	User
	BasePay          float64 `json:"base_pay"`
	NumberOfStudents int     `json:"number_of_students"`
}

type Admin struct {
	User
	BasePay   float64 `json:"base_pay"`
	Incentive float64 `json:"incentive"`
}

// Parent struct for storing parent information
type Parent struct {
	User
}

type Class struct {
	ClassID    string  `json:"class_id"`
	Name       string  `json:"class_name"`
	Instructor string  `json:"instructor"`
	Hours      float64 `json:"hours"`
}

// NewStudent creates a new Student instance
// func NewStudent(googleID, name, email, contactNumber, class, instructor, parentName, role string, age int, lessonCredits float32) Student {
func NewStudent(googleID, name, email, contactNumber, classID, parentName, role string, age int, lessonCredits float32) Student {
	return Student{
		User: User{
			GoogleID:      googleID,
			Name:          name,
			Email:         email,
			ContactNumber: contactNumber,
			Role:          role,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		Age:           age,
		LessonCredits: lessonCredits,
		ClassID:       classID,
		// Instructor:    instructor,
		ParentName: parentName,
	}
}

// NewInstructor creates a new Instructor instance
func NewInstructor(googleID, name, email, contactNumber, role string, basePay float64, numberOfStudents int) Instructor {
	return Instructor{
		User: User{
			GoogleID:      googleID,
			Name:          name,
			Email:         email,
			ContactNumber: contactNumber,
			Role:          role,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		BasePay:          basePay,
		NumberOfStudents: numberOfStudents,
	}
}

// NewAdmin creates a new Admin instance
func NewAdmin(googleID, name, email, contactNumber, role string, basePay, incentive float64) Admin {
	return Admin{
		User: User{
			GoogleID:      googleID,
			Name:          name,
			Email:         email,
			ContactNumber: contactNumber,
			Role:          role,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		BasePay:   basePay,
		Incentive: incentive,
	}
}

// NewParent creates a new Parent instance
func NewParent(googleID, name, email, contactNumber, role string) Parent {
	return Parent{
		User: User{
			GoogleID:      googleID,
			Name:          name,
			Email:         email,
			ContactNumber: contactNumber,
			Role:          role,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}
}

func NewClass(classID, name, instructor string, hours float64) Class {
	return Class{
		ClassID:    classID,
		Name:       name,
		Instructor: instructor,
		Hours:      hours,
	}
}
