package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


// Student must match the server's data structure exactly.
type Student struct {
	Name  string `json:"name"`
	Major string `json:"major"`
	Year  int    `json:"year"`
}

// getStudent fetches a student by name from the server.
func getStudent(name string) (Student, error) {
	resp, err := http.Get("http://localhost:8080/students/" + name)
	if err != nil {
		return Student{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Student{}, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	var student Student
	if err := json.NewDecoder(resp.Body).Decode(&student); err != nil {
		return Student{}, err
	}
	return student, nil
}

// addStudent sends a new student to the server.
func addStudent(s Student) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8080/students", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func listStudents() ([]Student, error) {
	resp, err := http.Get("http://localhost:8080/students")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	var students []Student
	if err := json.NewDecoder(resp.Body).Decode(&students); err != nil {
		return nil, err
	}
	return students, nil
}

func deleteStudent(name string) error {
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/students/"+name, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func main() {
	student, err := getStudent("alice")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Fetched:     Name=%-10s  Major=%-20s  Year=%d\n", student.Name, student.Major, student.Year)

	newStudent := Student{
		Name:  "charlie",
		Major: "Physics",
		Year:  1,
	}
	if err := addStudent(newStudent); err != nil {
		fmt.Println("Error:", err)
		return
	}

	student, err = getStudent("charlie")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("New student: Name=%-10s  Major=%-20s  Year=%d\n", student.Name, student.Major, student.Year)

	studentList, err := listStudents()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("\nAll students:")
	for _, s := range studentList {
		fmt.Printf("  - Name=%-10s  Major=%-20s  Year=%d\n", s.Name, s.Major, s.Year)
	}

	if err := deleteStudent("charlie"); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("\nDeleted charlie")

	studentList, err = listStudents()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Remaining students:")
	for _, s := range studentList {
		fmt.Printf("  - Name=%-10s  Major=%-20s  Year=%d\n", s.Name, s.Major, s.Year)
	}
}
