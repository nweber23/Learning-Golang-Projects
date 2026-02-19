# Learning Go

A collection of Go projects demonstrating language fundamentals, web development patterns, and CLI tooling.

## Projects

### CLI Tools

**Github User Activity**
Fetch and display recent GitHub user activity from the command line. Demonstrates REST API consumption, JSON parsing, and command-line argument handling.

**Email Verifier**
Check domain email readiness by verifying MX, SPF, and DMARC records using DNS lookups. Shows net package usage for domain verification.

**Number Guessing Game**
Interactive game where the player guesses a random number with difficulty levels and high score tracking. Demonstrates game loops, user input handling, and persistent state.

**Slack Bot for Age Calculation**
Slack bot that calculates age from a given date. Integrates with Slack API for interactive bot development.

### Web Servers & APIs

**Simple Web Server**
Basic HTTP server serving static files and handling form submissions. Shows routing, static file serving, and form data processing with Go's net/http.

**Basic CRUD API**
RESTful JSON API for book management with create, read, update, and delete operations. Uses gorilla/mux for routing and demonstrates JSON encoding/decoding.

**Personal Blog**
Full-featured blog application with public guest section and password-protected admin dashboard. Demonstrates server-side templating, session-based authentication, file-based storage, and CRUD operations.

### In Progress

**Image Process Service**
Service for image processing operations (not yet fully implemented).

## Getting Started

Each project is self-contained with its own directory and README. To work on a project:

```bash
cd "<Project Name>"
cat README.md
go run main.go
```

For projects with multiple files:

```bash
cd "<Project Name>"
go run .
```

To build an executable:

```bash
cd "<Project Name>"
go build -o output-name
```

## Requirements

- Go 1.21 or later
- For web projects: curl or a browser for testing

## Repository Structure

```
Learning-Go/
├── Basic CRUD API/           # REST API example
├── Email Verifier/           # DNS lookup CLI
├── Github User Activity/     # API consumption CLI
├── Image Process Service/    # Image processing (WIP)
├── Number Guessing Game/     # Interactive CLI game
├── Personal Blog/            # Full web application
├── Simple Web Server/        # Static file server
├── Slack Bot for Age Calculation/  # Slack integration
├── CLAUDE.md                 # Development guidance for Claude Code
└── README.md                 # This file
```

## Learning Path

If new to Go, a suggested order:

1. **Number Guessing Game** - Basic syntax, loops, user input
2. **Email Verifier** - Standard library (net, bufio)
3. **Github User Activity** - HTTP requests, JSON, error handling
4. **Simple Web Server** - HTTP routing, static files
5. **Basic CRUD API** - REST patterns, third-party packages
6. **Personal Blog** - Full application with auth, templates, persistence

## Development Notes

Each project focuses on specific Go concepts and best practices for that problem domain. Projects use:
- Standard library where possible
- Third-party packages (gorilla/mux) when appropriate for the learning goal
- File-based storage instead of databases (learning focus, not production)
- Hardcoded data and credentials where specified in individual project READMEs

For detailed development guidance, see CLAUDE.md.
