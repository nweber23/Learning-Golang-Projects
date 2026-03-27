# Exercise: Student Directory — Server & Client | Part of a Introduction into GO by Schwarz IT

## Goal

Build a simple **server-client system in Go** that exchanges student data over HTTP using JSON.

The **server** keeps a directory of students in memory and exposes two endpoints.
The **client** communicates with the server: it fetches students and creates new ones.

---

## What to implement

### Server (`server/main.go`)

The server should run on **port 8080** and provide two routes:

| Method   | Route              | Description                              |
|----------|--------------------|------------------------------------------|
| `GET`    | `/students/{name}` | Returns a single student as JSON.        |
| `POST`   | `/students`        | Accepts a new student as JSON.           |

**GET /students/{name}**
- Name is part of the URL path, e.g. `/students/alice`
- Returns the student as JSON on success (HTTP 200)
- Returns HTTP 404 if the student is not found

**POST /students**
- Body: JSON with `name`, `major`, `year`
- Stores the student in the in-memory map
- Returns HTTP 201

### Client (`client/main.go`)

The client should:
1. Fetch an existing student via `GET` and print the result
2. Create a new student via `POST`
3. Fetch the newly created student again and print the result

---

## Data Structure

```go
type Student struct {
    Name  string `json:"name"`
    Major string `json:"major"`
    Year  int    `json:"year"`
}
```

---

## Requirements

- [ ] The server must be started with `http.ListenAndServe`
- [ ] JSON is encoded/decoded using `encoding/json`
- [ ] Errors are handled correctly (`if err != nil`)
- [ ] `defer resp.Body.Close()` is used in the client
- [ ] HTTP status codes are set and checked correctly

---

## Bonus (optional)

- Return all stored students via `GET /students` as a JSON array
- Add a `DELETE /students/{name}` endpoint
- Print a nicely formatted output in the client using `fmt.Printf`
