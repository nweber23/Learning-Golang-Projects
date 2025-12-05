# Subject: Simple Web Server in Go

Goal
- Serve static files and handle a few HTTP routes using Go’s net/http.

You’ll Practice
- Defining handlers (GET/POST)
- Routing by path with the default mux
- Serving static files with http.FileServer
- Reading form data
- Setting headers and status codes

Run
1) cd "Simple Web Server"
2) go run main.go
3) Open http://localhost:8080

Routes
- GET  /            → serves ./static (e.g., index.html, form.html)
- GET  /hello       → plain-text “Hello!”
- POST /form        → echoes name and age from form data
- GET  /form.html   → static HTML form page

Behavior
- /hello: only GET allowed; non-GET returns 405 with Allow: GET
- /form: only POST allowed; non-POST returns 405 with Allow: POST
- /hello and /form respond with Content-Type: text/plain; charset=utf-8
- Static files are served from ./static
- Listen on :8080
- Run from the project folder so ./static resolves correctly

Try
- Static files: http://localhost:8080
- Hello:        curl -i http://localhost:8080/hello
- Form (POST):  curl -i -X POST -d "name=Alice&age=30" http://localhost:8080/form

Project Files
- main.go
- static/index.html
- static/form.html