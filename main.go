package main

import (
	"log"
	"net/http"
)

// func init() {
// 	// Example: Create a new student
// 	student := NewStudent("John Doe", 11, "Digital Navigators", "Scott Smith")
// 	err := fireDB.Create("/students/"+student.ID.String(), student)
// 	if err != nil {
// 		log.Fatalf("Failed to create student: %v", err)
// 	}
// 	fmt.Println("Student created successfully")

// 	// Example: Read the student data
// 	var readStudent Student
// 	err = fireDB.Read("/students/"+student.ID.String(), &readStudent)
// 	if err != nil {
// 		log.Fatalf("Failed to read student: %v", err)
// 	}
// 	fmt.Printf("Student data: %+v\n", readStudent)

// 	// Example: Update the student data
// 	updates := map[string]interface{}{
// 		"Age": 12,
// 	}
// 	err = fireDB.Update("/students/"+student.ID.String(), updates)
// 	if err != nil {
// 		log.Fatalf("Failed to update student: %v", err)
// 	}
// 	fmt.Println("Student updated successfully")

// 	// Example: Delete the student data
// 	err = fireDB.Delete("/students/" + student.ID.String())
// 	if err != nil {
// 		log.Fatalf("Failed to delete student: %v", err)
// 	}
// 	fmt.Println("Student deleted successfully")
// }

func main() {
	// http.HandleFunc("/1", serverhome)
	// http.HandleFunc("/2", setCookieHandler)
	// http.ListenAndServe(":8080", handler())
	// http.ListenAndServeTLS("192.168.1.129:8080", "server.crt", "server.key", handler())

	// err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", handler())
	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", handler())
	if err != nil {
		log.Fatal(err)
	}
}
