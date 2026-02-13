# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Learning Mode

This is a **learning and practice project**. Claude operates in strict mentoring mode:

- **NEVER** write code or run commands on behalf of the user unless explicitly asked
- **ALWAYS** explain what needs to be done and why, then let the user implement it
- Provide the commands or code snippets the user should run/write, but let them execute
- When the user makes a mistake, explain what went wrong and guide them to fix it themselves
- Offer educational insights about Go, Cobra, BubbleTea, and architectural patterns
- Ask the user to make architectural and design decisions — present trade-offs, not answers
- Review code the user writes and suggest improvements with explanations

## Project Overview

This is a Go-based API testing CLI/TUI application (potential names: `goapi`, `gorest`). It allows users to:
- Create projects with a base URL
- Add and manage API routes
- Test routes individually or in batch
- Support JWT authorization
- Use variables in route URLs (e.g., resource IDs)
- Output results to terminal or JSON file
- Operate via command-line flags or interactive TUI

**Tech Stack:**
- Go (language)
- Cobra (CLI framework)
- BubbleTea (TUI framework)
- SQLite (data persistence, if needed)

**Reference:** See Go-Api-Cli.md for detailed user stories and specification.

## Project Structure (As Development Progresses)

Once initialized, the project should follow this structure:

```
go-api-cli/
├── cmd/                      # CLI commands (Cobra)
│   └── main.go              # Entry point
├── internal/
│   ├── api/                 # API client logic
│   ├── project/             # Project management
│   ├── route/               # Route management
│   ├── auth/                # JWT and auth handling
│   ├── output/              # Result formatting (terminal, JSON)
│   ├── tui/                 # BubbleTea TUI components
│   └── storage/             # Data persistence (SQLite)
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── Makefile                 # Build and development commands
└── CLAUDE.md                # This file
```

## Common Development Commands

Once the project is initialized with `go mod init`:

```bash
# Build the binary
go build -o goapi ./cmd

# Run directly
go run ./cmd [arguments]

# Run tests
go test ./...
go test -v ./...         # Verbose
go test -run TestName ./path  # Single test

# Code quality
go fmt ./...             # Format code
go vet ./...             # Vet code for issues
goimports -w ./...       # Organize imports (if golangci-lint installed)

# Dependency management
go mod tidy              # Clean up dependencies
go mod download          # Download dependencies
```

## Architectural Guidance

### Command Structure (Cobra)
Organize Cobra commands by top-level operations:
- `project` - Create, list, delete projects
- `route` - Add, update, delete, list routes
- `auth` - Set/update JWT credentials
- `test` - Execute routes/projects
- `output` - Configure output format

Each command should have both CLI flag and TUI entry points.

### TUI Architecture (BubbleTea)
Use BubbleTea's Model/Update/View pattern:
- **Models:** Separate models for each screen (ProjectList, RouteEditor, TestResults)
- **Navigation:** Implement a state machine for screen transitions
- **Async Operations:** Use Cmd for API calls and long-running operations to keep UI responsive

### Data Storage
Decide early whether SQLite is necessary:
- **With SQLite:** Use for project/route persistence across sessions
- **Without SQLite:** Store in JSON files in a config directory (simpler, less dependency)

Consider: `~/.config/goapi/` or `~/.goapi/` for user data.

### API Client
- Abstract HTTP client behind an interface for easy mocking in tests
- Support common HTTP methods: GET, POST, PUT, DELETE, PATCH
- Handle status codes, timeouts, and connection errors gracefully

### JWT Handling
- Store auth tokens securely (avoid plaintext in config files)
- Support token refresh if the API provides it
- Allow per-project or global JWT configuration

## Key Decision Points

1. **Data Persistence:** SQLite vs JSON files - impacts database code and queries
2. **Configuration Location:** Where to store projects and routes
3. **Variable Syntax:** How to define variables in route URLs (e.g., `{id}`, `$id`, `{{id}}`)
4. **Output Formats:** JSON structure for test results, terminal formatting

## Testing Strategy

- Unit tests for API client, project/route logic, JWT handling
- Use mock HTTP server for API client tests
- Integration tests for full workflows
- TUI testing is limited by BubbleTea; focus on Model logic unit tests

## Current Status

### Completed (Session 1)

1. ✅ Initialized Go module: `github.com/raworiginal/goapi`
2. ✅ Set up basic project structure:
   - `cmd/main.go` — Root Cobra command (`goapi`)
   - `internal/project/project.go` — Project data model
   - `internal/storage/storage.go` — GORM + SQLite database layer

3. ✅ Implemented data models and storage:
   - **Project struct:** `ID` (primary key), `Name` (unique), `BaseURL`, `DateCreated`, `Description`
   - **Database location:** `~/.config/goapi/goapi.db` (SQLite)
   - **Database decisions:** Using GORM ORM + numeric ID as primary key
   - **CRUD operations:** `CreateProject()`, `GetProject()`, `ListProjects()`, `DeleteProject()`

4. ✅ Dependencies added:
   - Cobra (CLI framework)
   - GORM (ORM)
   - GORM SQLite driver

### Next Steps

**Session 2 priorities:**
1. Create `cmd/project.go` with Cobra subcommands:
   - `goapi project create --name "..." --url "..." [--description "..."]`
   - `goapi project list`
   - `goapi project delete --name "..."`
2. Integrate storage layer with project commands (use CreateProject, GetProject, ListProjects, DeleteProject)
3. Test the project command workflow
4. Create `cmd/route.go` with similar subcommand structure
5. Implement route data model and storage CRUD

### Architectural Decisions Made

- **Database:** SQLite with GORM (decided over JSON for scalability and learning value)
- **Primary Key:** Numeric ID with unique constraint on name (allows future renaming, better for relationships)
- **Data Location:** `~/.config/goapi/` (standard XDG-like path)
- **ORM:** GORM (selected for learning and maintainability)
