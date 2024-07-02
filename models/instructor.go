package main

// Instructor struct represents the structure of an instructor
type Instructor struct {
	GoogleID         string  `json:"googleID"`
	Name             string  `json:"name"`
	ContactNumber    string  `json:"contactNumber"`
	Email            string  `json:"email"`
	BasePay          float64 `json:"basePay"`
	NumberOfStudents int     `json:"numberOfStudents"`
	Role             string  `json:"role"`
}
