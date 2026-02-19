# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a collection of Go learning projects, each focusing on different language features and patterns. Each project is self-contained in its own directory with its own README, go.mod, and executable.

### Project Types

- **CLI Tools**: Standalone command-line applications that read input and process data
  - Github User Activity: REST API consumption and JSON parsing
  - Email Verifier: DNS lookups with net package
  - Number Guessing Game: Game loop, user input, and state management

- **Web Servers**: HTTP handlers using Go's net/http
  - Simple Web Server: Static file serving, form handling, HTTP routing
  - Personal Blog: Full web application with server-side templating, session-based auth, CRUD operations

- **APIs**: JSON-based REST services
  - Basic CRUD API: JSON encoding/decoding with gorilla/mux routing

## Development Commands

### Running a Project

Each project has its own directory with a README. The typical pattern is:

```bash
cd "<Project Name>"
go run main.go          # For single-file projects
go run .                # For multi-file projects with go.mod
```

### Building for Release

```bash
go build -o output-name        # Creates executable in current directory
```

### Testing

Run tests for a specific project:

```bash
go test ./...           # Run all tests in project
go test -v ./...        # With verbose output
go test -run TestName   # Run specific test
```

### Dependency Management

Projects use Go modules (go.mod). When adding dependencies:

```bash
go get github.com/package/name     # Add dependency
go mod tidy                        # Remove unused dependencies
go mod download                    # Download all dependencies locally
```

Current Go version: **1.21** (as specified in go.mod files)

## Code Architecture Patterns

### Project Structure

- **Single-file projects**: CLI tools are typically just `main.go` + go.mod
- **Multi-file projects**: Web applications use handler packages
  - `main.go`: Initializes routes and starts server
  - `handlers/`: HTTP handler functions grouped by functionality
  - `middleware/`: Authentication, logging, CORS (e.g., session auth wrapper)
  - `static/`: CSS, images, client-side assets
  - `templates/`: HTML templates for server-side rendering
  - `articles/`: Data storage (JSON files)

### Common Patterns

**CLI Entry Points:**
- Use `fmt.Scan()`, `bufio.Scanner`, or `flag` package for input
- Handle errors with logging and graceful exit codes
- No external dependencies unless necessary

**HTTP Handlers:**
- Typed as `func(http.ResponseWriter, *http.Request) error` or `func(http.ResponseWriter, *http.Request)`
- Set `Content-Type` header for JSON responses
- Use standard `net/http` mux or `gorilla/mux` for routing
- Middleware wraps handlers for auth checks

**Session Management (Personal Blog):**
- Uses `net/http.Cookie` for session tokens
- Middleware checks for valid session before allowing access to protected routes
- Credentials are hardcoded for learning (admin/admin in Personal Blog)

**File-based Storage:**
- Articles stored as JSON files in `articles/` directory
- Read/write with `os` and `encoding/json` packages

## Running and Testing Examples

### CLI Projects

```bash
cd "Github User Activity"
go build -o github-activity
./github-activity octocat

cd "Email Verifier"
echo "example.com" | go run main.go

cd "Number Guessing Game"
go run main.go              # Interactive game
```

### Web Projects

```bash
cd "Simple Web Server"
go run main.go              # Serves on :8080
# Test: curl http://localhost:8080/hello

cd "Basic CRUD API"
go run .                    # API on :8080
# Test: curl http://localhost:8080/api/books

cd "Personal Blog"
go run main.go              # On :8080
# Guest: http://localhost:8080
# Admin: http://localhost:8080/admin/login (admin/admin)
```

## Key Dependencies

- **gorilla/mux**: Used in Basic CRUD API and Personal Blog for advanced routing
- **Standard library**: Preferred for most projects (net/http, encoding/json, net for DNS, math/rand, bufio)

## Important Notes

- Each project directory should be run from its own root (so relative paths like `./static` resolve correctly)
- Projects in this repo are learning exercises—use them as references for patterns, not production templates
- Some projects use hardcoded data (books in CRUD API, articles in Personal Blog) rather than databases
- The "Image Process Service" directory exists but is not yet fully implemented
