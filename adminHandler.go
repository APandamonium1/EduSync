package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// AdminHandler handles all admin-related routes.
//
// This handler is responsible for setting up routes for admin-related tasks,
// such as searching for students, parents, instructors, and announcements,
// as well as editing and updating their information.
//
// Example usage:
//   router := mux.NewRouter()
//   AdminHandler(router)

func AdminHandler(router *mux.Router) {
	router.HandleFunc("/admin", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/index.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/admin/profile", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/profile.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/admin/profile", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/profile.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	//searve the search student page
	router.HandleFunc("/admin/search_student", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/search_student.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	// Serve the create student page
	router.HandleFunc("/admin/create_student", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/create_student.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	// Search for students by name and class
	router.HandleFunc("/admin/api/search_student", func(res http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		class := req.URL.Query().Get("class")
		students, err := searchStudents(name, class)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(students)
	}).Methods("GET")

	// Example usage:
	//   GET /admin/api/search_student?name=John&class=10th
	//   Response: JSON list of students matching the search criteria

	// Edit student information
	router.HandleFunc("/admin/student/{googleID}/edit", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]

		student, err := readStudent(googleID, req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("templates/admin/edit_student.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, student)
	}).Methods("GET")

	// Example usage:
	//   GET /admin/student/1234567890/edit
	//   Response: HTML form to edit student information

	// Update student information
	router.HandleFunc("/admin/student/{googleID}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]

		switch req.Method {
		case http.MethodGet:
			student, err := readStudent(googleID, req)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(student)
		case http.MethodPut:
			var updates map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			if err := updateStudent(googleID, updates, req); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT")

	// Example usage:
	//   PUT /admin/student/1234567890
	//   Request Body: JSON object with updated student information
	//   Response: HTTP Status No Content (204)

	// Create a new student
	router.HandleFunc("/admin/student/", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			var student Student // Replace 'Student' with your actual student struct type
			if err := json.NewDecoder(req.Body).Decode(&student); err != nil {
				http.Error(res, fmt.Sprintf(`{"error": "Invalid request payload: %v"}`, err), http.StatusBadRequest)
				return
			}
			student.GoogleID = uuid.New().String()
			student.Role = "Student"
			student.CreatedAt = time.Now()
			student.UpdatedAt = time.Now()
			if err := createStudent(student, req); err != nil { // Implement createStudent to handle the logic
				http.Error(res, fmt.Sprintf(`{"error": "Failed to create student: %v"}`, err), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusCreated)
			json.NewEncoder(res).Encode(student)
		default:
			http.Error(res, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		}
	}).Methods("POST")

	// Serve the search parent page
	router.HandleFunc("/admin/search_parent", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/search_parent.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	// Serve the create parent page
	router.HandleFunc("/admin/create_parent", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/create_parent.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	router.HandleFunc("/admin/api/search_parent", func(res http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		parents, err := searchParents(name)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(parents)
	}).Methods("GET")

	router.HandleFunc("/admin/parent/{googleID}/edit", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]

		parent, err := readParent(googleID, req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("templates/admin/edit_parent.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, parent)
	}).Methods("GET")

	router.HandleFunc("/admin/parent/{googleID}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]

		switch req.Method {
		case http.MethodGet:
			parent, err := readParent(googleID, req)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(parent)
		case http.MethodPut:
			var updates map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			if err := updateParent(googleID, updates, req); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT")

	// Create a new parent
	router.HandleFunc("/admin/parent/", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			var parent Parent
			if err := json.NewDecoder(req.Body).Decode(&parent); err != nil {
				http.Error(res, fmt.Sprintf(`{"error": "Invalid request payload: %v"}`, err), http.StatusBadRequest)
				return
			}
			parent.GoogleID = uuid.New().String()
			parent.Role = "Parent"
			parent.CreatedAt = time.Now()
			parent.UpdatedAt = time.Now()
			if err := createParent(parent, req); err != nil {
				http.Error(res, fmt.Sprintf(`{"error": "Failed to create parent: %v"}`, err), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusCreated)
			json.NewEncoder(res).Encode(parent)
		default:
			http.Error(res, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		}
	}).Methods("POST")

	// Serve the search instructor page
	router.HandleFunc("/admin/search_instructor", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/search_instructor.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	// Serve the create instructor page
	router.HandleFunc("/admin/create_instructor", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/create_instructor.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	// Serve the search instructor API endpoint
	router.HandleFunc("/admin/api/search_instructor", func(res http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		instructors, err := searchInstructors(name)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(instructors)
	}).Methods("GET")

	router.HandleFunc("/admin/instructor/{googleID}/edit", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]

		instructor, err := readInstructor(googleID, req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("templates/admin/edit_instructor.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, instructor)
	}).Methods("GET")

	router.HandleFunc("/admin/instructor/{googleID}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		googleID := vars["googleID"]

		switch req.Method {
		case http.MethodGet:
			instructor, err := readInstructor(googleID, req)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(instructor)
		case http.MethodPut:
			var updates map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			if err := updateInstructor(googleID, updates, req); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT")

	// Create a new instructor
	router.HandleFunc("/admin/instructor/", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			var instructor Instructor
			if err := json.NewDecoder(req.Body).Decode(&instructor); err != nil {
				http.Error(res, fmt.Sprintf(`{"error": "Invalid request payload: %v"}`, err), http.StatusBadRequest)
				return
			}
			instructor.GoogleID = uuid.New().String()
			instructor.Role = "Instructor"
			instructor.CreatedAt = time.Now()
			instructor.UpdatedAt = time.Now()
			if err := createInstructor(instructor, req); err != nil {
				http.Error(res, fmt.Sprintf(`{"error": "Failed to create instructor: %v"}`, err), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusCreated)
			json.NewEncoder(res).Encode(instructor)
		default:
			http.Error(res, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		}
	}).Methods("POST")

	//Serve the search announcement page
	router.HandleFunc("/admin/search_announcement", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/search_announcement.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	// Serve the create announcement page
	router.HandleFunc("/admin/create_announcement", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/create_announcement.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	// Search for an announcement
	router.HandleFunc("/admin/api/search_announcement", func(res http.ResponseWriter, req *http.Request) {
		subject := req.URL.Query().Get("subject")
		// content := req.URL.Query().Get("content")
		announcements, err := searchAnnouncements(subject)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(announcements)
	}).Methods("GET")

	router.HandleFunc("/admin/announcement/{announcementID}/edit", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		announcementID := vars["announcementID"]

		announcement, err := readAnnouncement(announcementID, req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("templates/admin/edit_announcement.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, announcement)
	}).Methods("GET")

	router.HandleFunc("/admin/announcement/{announcementID}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		announcementID := vars["announcementID"]

		switch req.Method {
		case http.MethodGet:
			announcement, err := readAnnouncement(announcementID, req)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(announcement)
		case http.MethodPut:
			var updates map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			if err := updateAnnouncement(announcementID, updates, req); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT")

	// Create a new announcement
	router.HandleFunc("/admin/announcement", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			var announcement Announcement
			if err := json.NewDecoder(req.Body).Decode(&announcement); err != nil {
				http.Error(res, fmt.Sprintf(`{"error": "Invalid request payload: %v"}`, err), http.StatusBadRequest)
				return
			}
			announcement.AnnouncementID = uuid.New().String()
			announcement.CreatedAt = time.Now()
			announcement.UpdatedAt = time.Now()
			if err := createAnnouncement(announcement, req); err != nil {
				http.Error(res, fmt.Sprintf(`{"error": "Failed to create announcement: %v"}`, err), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusCreated)
			json.NewEncoder(res).Encode(announcement)
		default:
			http.Error(res, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		}
	}).Methods("POST")

	// Example usage:
	//   POST /admin/announcement
	//   Request Body: JSON object with announcement details
	//   Response: HTTP Status Created (201)

	// Serve the search class page
	router.HandleFunc("/admin/search_class", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/search_class.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	// Serve the create class page
	router.HandleFunc("/admin/create_class", func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFiles("templates/admin/create_class.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, nil)
	}).Methods("GET")

	// Search for classes by classID
	router.HandleFunc("/admin/api/search_class", func(res http.ResponseWriter, req *http.Request) {
		classID := req.URL.Query().Get("classID")
		classes, err := searchClasses(classID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(classes)
	}).Methods("GET")

	// Edit class information page
	router.HandleFunc("/admin/class/{classID}/edit", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		classID := vars["classID"]

		// Retrieve students and class information
		students, err := GetStudentsByClassID(classID)
		if err != nil {
			http.Error(res, fmt.Sprintf("error retrieving students: %v", err), http.StatusInternalServerError)
			return
		}

		class, err := readClass(students, classID, req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("templates/admin/edit_class.html")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(res, class)
	}).Methods("GET")

	// Update class information
	router.HandleFunc("/admin/class/{classID}", func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		classID := vars["classID"]

		// Retrieve students
		students, err := GetStudentsByClassID(classID)
		if err != nil {
			http.Error(res, fmt.Sprintf("error retrieving students: %v", err), http.StatusInternalServerError)
			return
		}

		switch req.Method {
		case http.MethodGet:
			class, err := readClass(students, classID, req)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(class)
		case http.MethodPut:
			var updates map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			class := Class{ClassID: classID}
			if err := updateClass(class, updates, req); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT")

	// Create a new class
	router.HandleFunc("/admin/class/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			var class Class
			if err := json.NewDecoder(req.Body).Decode(&class); err != nil {
				http.Error(res, fmt.Sprintf(`{"error": "Invalid request payload: %v"}`, err), http.StatusBadRequest)
				return
			}
			// class.ClassID = uuid.New().String()
			// class.CreatedAt = time.Now()
			// class.UpdatedAt = time.Now()
			if err := createClass(class, req); err != nil {
				http.Error(res, fmt.Sprintf(`{"error": "Failed to create class: %v"}`, err), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusCreated)
			json.NewEncoder(res).Encode(class)
		} else {
			http.Error(res, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		}
	}).Methods("POST")

	router.HandleFunc("/admin/api/profile", func(res http.ResponseWriter, req *http.Request) {
		currentUser, err := GetCurrentUser(req)
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}
		switch req.Method {
		case http.MethodGet:
			admin, err := readAdmin(currentUser.GoogleID, req)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(admin)
		case http.MethodPut:
			var updates map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			if err := updateAdmin(currentUser.GoogleID, updates, req); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusNoContent)
		}
	}).Methods("GET", "PUT")
}
