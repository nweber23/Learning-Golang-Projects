# Guide: Student Directory

This guide walks you through the exercise step by step.

---

## Step 1 — Set up the project

Open two terminals. Both start from inside the `go-skeleton/` directory.

```bash
# Terminal 1: Server
cd server
go mod init student-server
go run main.go

# Terminal 2: Client
cd client
go mod init student-client
go run main.go
```

> Tip: Always start the server first, then the client.

---

## Step 2 — Understand the struct definition

Both server and client need the same data structure:

```go
type Student struct {
    Name  string `json:"name"`
    Major string `json:"major"`
    Year  int    `json:"year"`
}
```

The **struct tags** (`json:"name"`) control how the fields appear in JSON.
Without tags, Go would use the field names directly (`Name`, `Major`, `Year`).

---

## Step 3 — Implement the server

### 3.1 — Handler function for GET

Every handler always has this signature:

```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    // w = what you write the response into
    // r = the incoming request
}
```

**Steps for the GET handler:**
1. Check the HTTP method: `r.Method != http.MethodGet`
2. Extract the name from the URL path: `strings.TrimPrefix(r.URL.Path, "/students/")`
3. Look up the student in the `students` map: `student, ok := students[name]`
4. On success: set the Content-Type header and encode the struct as JSON

```go
// Tell the client: this response contains JSON
w.Header().Set("Content-Type", "application/json")

// Write the struct as JSON into the response
json.NewEncoder(w).Encode(student)
```

### 3.2 — Handler function for POST

**Steps for the POST handler:**
1. Check the HTTP method: `r.Method != http.MethodPost`
2. Read and decode the JSON body into a struct:

```go
var student Student
if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
    http.Error(w, "Invalid JSON", http.StatusBadRequest)
    return
}
```

3. Save the student in the map: `students[student.Name] = student`
4. Return HTTP 201: `w.WriteHeader(http.StatusCreated)`

### 3.3 — Register routes and start the server

```go
func main() {
    http.HandleFunc("/students/", handleGetStudent)
    http.HandleFunc("/students",  handleAddStudent)

    fmt.Println("Server running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

---

## Step 4 — Implement the client

### 4.1 — Send a GET request

```go
resp, err := http.Get("http://localhost:8080/students/" + name)
if err != nil {
    return Student{}, err
}
defer resp.Body.Close() // ALWAYS close the body!
```

Why `defer`? The response body is an open network connection.
`defer` ensures it is closed at the end of the function — regardless of whether an error occurs or not.

### 4.2 — Read the JSON response

```go
var student Student
err = json.NewDecoder(resp.Body).Decode(&student)
```

### 4.3 — Send a POST request with a JSON body

```go
// Convert the struct to JSON bytes
data, err := json.Marshal(newStudent)

// POST with JSON body
resp, err := http.Post(
    "http://localhost:8080/students",
    "application/json",
    bytes.NewBuffer(data),
)
```

---

## Step 5 — Test it

You can also test the server directly with `curl`:

```bash
# Fetch a student
curl "http://localhost:8080/students/alice"

# Add a new student
curl -X POST http://localhost:8080/students \
     -H "Content-Type: application/json" \
     -d '{"name":"dave","major":"Physics","year":1}'
```

---

## Concepts Cheat Sheet

| Concept                | Go syntax                              |
|------------------------|----------------------------------------|
| Define a struct        | `type Student struct { ... }`          |
| Create a struct        | `s := Student{Name: "Alice", ...}`     |
| Create a map           | `m := map[string]Student{}`            |
| Read map with ok-check | `v, ok := m["key"]`                    |
| Encode JSON            | `json.NewEncoder(w).Encode(v)`         |
| Decode JSON            | `json.NewDecoder(r.Body).Decode(&v)`   |
| Start HTTP server      | `http.ListenAndServe(":8080", nil)`    |
| Register handler       | `http.HandleFunc("/route", handler)`   |
| Send GET request       | `http.Get(url)`                        |
| Send POST request      | `http.Post(url, contentType, body)`    |
| Check error            | `if err != nil { return ..., err }`    |
| Defer cleanup          | `defer resp.Body.Close()`              |
