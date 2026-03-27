package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Student represents a student in the directory.
// Struct tags control the JSON field names.
type Student struct {
	Name  string `json:"name"`
	Major string `json:"major"`
	Year  int    `json:"year"`
}

// In-memory store: name -> Student
var students = map[string]Student{
	"alice": {Name: "alice", Major: "Computer Science", Year: 3},
	"bob":   {Name: "bob", Major: "Mathematics", Year: 2},
}

// handleGetStudent handles GET /students/{name}
func handleGetStudent(w http.ResponseWriter, r *http.Request) {
	// TODO 1: Check that the HTTP method is GET.
	//         If not: http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) and return.

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// TODO 2: Extract the student name from the URL path.
	//         strings.TrimPrefix(r.URL.Path, "/students/") returns the name.
	//         If empty: http.Error(w, "...", http.StatusBadRequest) and return.

	name := strings.TrimPrefix(r.URL.Path, "/students/")
	if name == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// TODO 3: Look up the student in the "students" map.
	//         Use the comma-ok idiom: student, ok := students[name]
	//         If not found: http.Error(w, "...", http.StatusNotFound) and return.

	student, ok := students[name]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// TODO 4: Set the Content-Type header to "application/json".
	//         w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Content-Type", "application/json")

	// TODO 5: Encode the student as JSON and write it into the response.
	//         json.NewEncoder(w).Encode(student)

	if err := json.NewEncoder(w).Encode(student); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAddStudent handles POST /students
func handleAddStudent(w http.ResponseWriter, r *http.Request) {
	// TODO 1: Check that the HTTP method is POST.

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// TODO 2: Decode the JSON body into a Student struct.
	//         var student Student
	//         err := json.NewDecoder(r.Body).Decode(&student)
	//         Don't forget to check the error!

	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO 3: Save the student in the map: students[student.Name] = student

	students[student.Name] = student

	// TODO 4: Set the status code to 201 Created: w.WriteHeader(http.StatusCreated)
	//         Return the new student as JSON.

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(student); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleListStudents(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	studentList := make([]Student, 0, len(students))
	for _, s := range students {
		studentList = append(studentList, s)
	}
	if err := json.NewEncoder(w).Encode(studentList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleDeleteStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, "/students/")
	if name == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if _, ok := students[name]; !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	delete(students, name)
	w.WriteHeader(http.StatusNoContent)
}

func handleRouter(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/students":
		handleAddStudent(w, r)
	case r.Method == http.MethodGet && (r.URL.Path == "/students" || r.URL.Path == "/students/"):
		handleListStudents(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/students/"):
		handleGetStudent(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/students/"):
		handleDeleteStudent(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/students", handleRouter)
	http.HandleFunc("/students/", handleRouter)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
