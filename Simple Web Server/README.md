# Subject: Simple Web Server in Go

Goal
- Build a tiny HTTP server using only Go’s standard library.

You’ll Practice
- Defining handlers (GET/POST)
- Routing by path
- Serving static files
- Reading form data
- Setting headers and status codes

Run
1) cd "Simple Web Server"
2) go run main.go
3) Open http://localhost:8080

Try
- Static files: http://localhost:8080
- Hello:        http://localhost:8080/hello
- Form (POST):  curl -X POST -d "name=Alice&age=30" http://localhost:8080/form

Project Files
- main.go
- static/index.html
- static/form.html

Mini‑Tasks
1) Add a new GET /time that returns the current time.
2) Make /hello return JSON when the Accept header includes application/json.
3) Write a test for formHandler that checks 405 on non‑POST.

Gotchas
- Run from the “Simple Web Server” folder so ./static resolves.
- For method mismatches, return 405 and set the Allow header.
- Call r.ParseForm() before reading form values.

Next Steps
- Use html/template to render responses
- Add /healthz
- Read PORT from an env var (default 8080)