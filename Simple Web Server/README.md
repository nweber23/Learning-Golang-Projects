Simple Web Server (Go)

Overview
- Minimal HTTP server using Go’s net/http.
- Serves static files from ./static and exposes /hello and /form endpoints.

Requirements
- Go installed (1.20+ recommended).

Run (Linux)
1) cd "Simple Web Server"
2) go run main.go
3) Open http://localhost:8080

Endpoints
- GET /            → Serves ./static (index.html by default).
- GET /hello       → Returns "Hello!" (text/plain).
- POST /form       → Expects form fields "name" and "age"; echoes them (text/plain).

Project Structure
- main.go
- static/index.html    (links to form.html)
- static/form.html     (POSTs to /form)

Quick Checks
- Hello: curl http://localhost:8080/hello
- Form:  curl -X POST -d "name=Alice&age=30" http://localhost:8080/form

Port
- 8080