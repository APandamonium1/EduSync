package main

import (
	"time"
	// "github.com/google/uuid"
)

// User represent the base user with common fields
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

// Student represents struct for storing student information
type Student struct {
	User
	Age           int     `json:"age"`
	LessonCredits float32 `json:"lesson_credits"`
	ClassID       string  `json:"class_id"`
	ParentID      string  `json:"parent_id"`
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
	FolderID string `json:"folder_id"`
}

type Class struct {
	ClassID    string    `json:"class_id"`
	FolderID   string    `json:"folder_id"`
	Name       string    `json:"class_name"`
	Instructor string    `json:"instructor"`
	Duration   float64   `json:"duration"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Announcement struct {
	AnnouncementID string    `json:"announcement_id"`
	Subject        string    `json:"subject"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// NewStudent creates a new Student instance
func NewStudent(googleID, name, email, contactNumber, classID, parentID, role string, age int, lessonCredits float32) Student {
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
		ParentID:      parentID,
		// Instructor:    instructor,
		// ParentName: parentName,
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

func NewClass(classID, name, instructor string, duration float64, createdAt, updatedAt time.Time) Class {
	return Class{
		ClassID:    classID,
		Name:       name,
		Instructor: instructor,
		Duration:   duration,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func NewAnnouncement(announcementID, subject, content string, createdAt, updatedAt time.Time) Announcement {
	return Announcement{
		AnnouncementID: announcementID,
		Subject:        subject,
		Content:        content,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
