package main

// Student struct represents the structure of student data
type Student struct {
	GoogleID   string `json:"googleID"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Age        int    `json:"age"`
	Class      string `json:"class"`
	Instructor string `json:"instructor"`
	ParentName string `json:"parentName"`
	Role       string `json:"role"`
}
