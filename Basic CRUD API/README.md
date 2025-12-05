# Subject: Basic CRUD API in Go

Goal
- Build a tiny RESTful JSON API using Go and gorilla/mux.

You’ll Practice
- Defining CRUD handlers (GET/POST/PUT/DELETE)
- Routing with path params (gorilla/mux)
- JSON encoding/decoding with net/http
- Setting response headers and status codes

Run
1) cd "Basic CRUD API"
2) go run .
3) API at http://localhost:8080

Endpoints
- GET    /api/books           → list all books
- GET    /api/books/{id}      → get a book by ID
- POST   /api/books           → create a book (auto‑generates string ID)
- PUT    /api/books/{id}      → replace a book by ID
- DELETE /api/books/{id}      → delete a book by ID

Data Model
- Book: { id: string, title: string, author: string }
- Storage: in‑memory slice seeded with:
  - 1: "1984" by George Orwell
  - 2: "The Great Gatsby" by F. Scott Fitzgerald

Behavior
- Content-Type is application/json for all responses.
- GET /api/books/{id}: 404 when not found.
- PUT /api/books/{id}: 404 when not found; replaces entire book.
- DELETE /api/books/{id}: returns the remaining list.

Try
- List:         curl http://localhost:8080/api/books
- Get by ID:    curl http://localhost:8080/api/books/1
- Create:       curl -X POST -H "Content-Type: application/json" -d '{"title":"Dune","author":"Frank Herbert"}' http://localhost:8080/api/books
- Update:       curl -X PUT  -H "Content-Type: application/json" -d '{"title":"Dune (Updated)","author":"Frank Herbert"}' http://localhost:8080/api/books/1
- Delete:       curl -X DELETE http://localhost:8080/api/books/1

Project Files
- main.go
- go.mod