# goapi - API Testing CLI

A command-line tool for managing and testing REST APIs. Create projects, define routes, and execute tests with a single command.

## Features

- **Project Management**: Create and organize API projects with a base URL
- **Route Management**: Define API routes (GET, POST, PUT, DELETE, PATCH) with human-readable names
- **Test Execution**: Execute routes individually or in batch and view results in a formatted table
- **Error Handling**: Graceful error reporting with timeouts and connection error support
- **Response Metrics**: Track HTTP status codes and response duration for each request

## Installation

### Prerequisites

- Go 1.26.0 or later

### Build

```bash
go build -o goapi ./cmd
```

### Run Directly (without building)

```bash
go run ./cmd [command] [flags]
```

## Usage

### Project Commands

#### Create a Project

```bash
goapi project create --name "MyAPI" --url "https://api.example.com" [--description "My API"]
```

**Flags:**
- `--name` (required): Project name
- `--url` (required): Base URL for all routes in this project
- `--description` (optional): Project description

**Example:**
```bash
goapi project create --name "GitHub API" --url "https://api.github.com" --description "GitHub v3 REST API"
```

#### List Projects

```bash
goapi project list
```

Shows all projects in a table format with name, base URL, description, and creation date.

#### Delete a Project

```bash
goapi project delete --name "MyAPI"
```

**Flags:**
- `--name` (required): Project name to delete

---

### Route Commands

#### Add a Route

```bash
goapi route add --project "MyAPI" --method GET --path "/users" [--name "Get Users"] [--description "Fetch all users"]
```

**Flags:**
- `--project` (required): Project name
- `--method` (required): HTTP method (GET, POST, PUT, DELETE, PATCH)
- `--path` (required): Route path (e.g., `/users`, `/users/{id}`)
- `--name` (optional): Human-readable name. Auto-generated as "METHOD path" if not provided
- `--description` (optional): Route description

**Example:**
```bash
goapi route add --project "MyAPI" --method POST --path "/users" --name "Create User"
goapi route add --project "MyAPI" --method GET --path "/users/{id}" --name "Get User by ID"
```

#### List Routes

```bash
goapi route list --project "MyAPI"
```

Shows all routes in a project in a table format with ID, name, method, and path.

#### Update a Route

```bash
goapi route update --project "MyAPI" --route "Get Users" [--method POST] [--path "/v2/users"] [--rename "New Name"]
```

**Flags:**
- `--project` (required): Project name
- `--route` (required): Route name (the route to update)
- `--method` (optional): New HTTP method
- `--path` (optional): New path
- `--rename` (optional): New name for the route
- `--description` (optional): New description

**Example:**
```bash
goapi route update --project "MyAPI" --route "Get Users" --rename "List All Users"
```

#### Delete a Route

```bash
goapi route delete --project "MyAPI" --route "Get Users"
```

**Flags:**
- `--project` (required): Project name
- `--route` (required): Route name to delete

---

### Test Commands

#### Test All Routes in a Project

```bash
goapi test --project "MyAPI" [--timeout 5s]
```

Executes all routes in the project and displays results in a table.

**Flags:**
- `--project` (required): Project name
- `--route` (optional): Test a specific route by name (if omitted, tests all routes)
- `--timeout` (optional): Request timeout (default: 5s)

**Example:**
```bash
goapi test --project "MyAPI"
```

#### Test a Single Route

```bash
goapi test --project "MyAPI" --route "Get Users"
```

Tests only the specified route.

**Example:**
```bash
goapi test --project "MyAPI" --route "Create User" --timeout 10s
```

#### Test Output

The test command displays results in a table with the following columns:

| Column | Description |
|--------|-------------|
| ID | Route ID |
| Name | Route name |
| Method | HTTP method |
| Path | Route path |
| Status | HTTP status code (or "Error" if request failed) |
| Duration | Response time |

**Example Output:**
```
ID  Name          Method  Path   Status  Duration
1   Get Users     GET     /users 200     125.456789ms
2   Create User   POST    /users 201     234.567891ms
3   Delete User   DELETE  /users/1 404   89.123456ms
4   Slow Endpoint GET     /delay/5 Error 1.000234567s
```

---

## Complete Workflow Example

```bash
# 1. Create a project
goapi project create --name "JSONPlaceholder" --url "https://jsonplaceholder.typicode.com"

# 2. Add some routes
goapi route add --project "JSONPlaceholder" --method GET --path "/posts" --name "List Posts"
goapi route add --project "JSONPlaceholder" --method GET --path "/posts/1" --name "Get Post 1"
goapi route add --project "JSONPlaceholder" --method POST --path "/posts" --name "Create Post"

# 3. List your routes
goapi route list --project "JSONPlaceholder"

# 4. Test all routes
goapi test --project "JSONPlaceholder"

# 5. Test a single route with custom timeout
goapi test --project "JSONPlaceholder" --route "Get Post 1" --timeout 3s

# 6. Update a route
goapi route update --project "JSONPlaceholder" --route "List Posts" --rename "Fetch All Posts"

# 7. Delete a route
goapi route delete --project "JSONPlaceholder" --route "Create Post"
```

---

## Data Storage

Projects and routes are stored in a SQLite database located at:

```
~/.config/goapi/goapi.db
```

This directory is created automatically on first run.

---

## Development Commands

```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run tests
go test ./...
go test -v ./...         # Verbose output
go test -run TestName ./path  # Run specific test

# Clean up dependencies
go mod tidy
```

---

## Future Features

### 1. JWT/Bearer Token Authentication
- Store and manage API tokens securely
- Automatically inject authentication headers in requests
- Support per-project or global token configuration
- `goapi auth set --project "MyAPI" --token "YOUR_JWT_TOKEN"`

### 2. Request Customization
- Add custom headers to routes (e.g., `Content-Type`, `Accept`)
- Define request bodies for POST/PUT/PATCH requests
- Support for route variables (e.g., `{id}`, `{user_id}`) with value substitution
- Dynamic request templates with environment variables

### 3. Output Formatting & Export
- Export test results to JSON file for CI/CD integration
- Pretty-print results in terminal with color coding
- Generate HTML reports with test summaries
- Batch testing with aggregated pass/fail statistics
- `goapi test --project "MyAPI" --output results.json`

### 4. Interactive TUI Mode
- Browse projects and routes with a full-featured terminal UI
- Execute tests with real-time result updates
- Edit routes and projects directly in the UI
- Keyboard shortcuts for common operations
- `goapi tui` - launch interactive mode

### 5. Advanced Testing
- Test scheduling and automation
- Response validation (assert status codes, headers, body content)
- Load testing and performance profiling
- Retry logic with exponential backoff for flaky endpoints
- Request/response logging and debugging

### 6. Configuration Management
- Import/export projects and routes to YAML/JSON
- Share configurations across team members
- Version control friendly project definitions
- Environment-specific configurations

---

## Architecture

### Project Structure

```
go-api-cli/
├── cmd/                      # CLI commands
│   ├── main.go              # Entry point
│   ├── project.go           # Project commands
│   ├── route.go             # Route commands
│   └── test.go              # Test command
├── internal/
│   ├── api/                 # HTTP client abstraction
│   │   └── client.go        # API request execution
│   ├── project/             # Project data model
│   ├── route/               # Route data model
│   └── storage/             # Database operations (SQLite + GORM)
├── go.mod                   # Go module definition
└── README.md                # This file
```

### Key Design Decisions

- **CLI Framework**: Cobra for command structure
- **Database**: SQLite with GORM ORM
- **HTTP Client**: Standard library `net/http` with custom abstraction
- **Output Formatting**: `text/tabwriter` for aligned table output

---

## License

MIT

## Contributing

Contributions are welcome! This is a learning project designed to practice Go, Cobra, and software architecture patterns.
