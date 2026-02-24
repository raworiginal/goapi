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

### Completed (Sessions 1-3)

**Session 1:**
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

**Session 2:**
1. ✅ Created `cmd/project.go` with fully functional Cobra subcommands:
   - `goapi project create --name "..." --url "..." [--description "..."]`
     - Validates URLs using `net/url.Parse()`
     - Returns user-friendly errors for duplicate names
     - Prints success confirmation
   - `goapi project list`
     - Displays all projects in format: `Name - BaseURL (DateCreated)`
     - Shows "No projects found" when list is empty
   - `goapi project delete --name "..."`
     - Removes projects from database
     - Proper error handling

2. ✅ Integrated storage layer with CLI:
   - Updated `cmd/main.go` with `PersistentPreRunE` that calls `storage.InitDB()`
   - All commands properly handle database errors
   - Cobra's RunE pattern for automatic error handling

3. ✅ Tested the project command workflow:
   - Create, list, and delete operations work end-to-end
   - Error handling verified (missing flags, duplicate names, database errors)
   - Help messages display correctly for all commands

**Session 3:**
1. ✅ Refactored storage layer for scalability:
   - Split `internal/storage/storage.go` into multiple files:
     - `storage.go` — Database initialization and global DB instance
     - `project.go` — Project CRUD operations
     - `route.go` — Route CRUD operations (see below)

2. ✅ Implemented route data model in `internal/route/route.go`:
   - **HTTPMethod custom type** — GET, POST, PUT, DELETE, PATCH constants for type safety
   - **Route struct:** `ID`, `ProjectID` (foreign key), `Method` (HTTPMethod), `Path`, `Description`, `DateCreated`
   - **Relationships:** Route belongs to Project (enforced at database level with GORM)

3. ✅ Implemented complete route storage CRUD in `internal/storage/route.go`:
   - `CreateRoute(r *route.Route) error` — Add routes to project
   - `ListRoutesByProject(projectID uint) ([]*route.Route, error)` — Fetch all routes for a project
   - `GetRoute(id uint) (*route.Route, error)` — Fetch single route by ID
   - `UpdateRoute(id uint, updates *route.UpdateRouteInput) error` — Update route fields (type-safe)
   - `DeleteRoute(id uint) error` — Remove route with "not found" error handling

4. ✅ Implemented type-safe updates:
   - **UpdateRouteInput struct** — Uses pointers (`*HTTPMethod`, `*string`) to distinguish "not provided" from "empty"
   - GORM `Updates()` only modifies non-nil fields, preventing accidental overwrites

5. ✅ Added route migration to `InitDB()`:
   - `DB.AutoMigrate(&route.Route{})` creates/updates route table

**Session 4:**
1. ✅ Created `cmd/route.go` with all Cobra subcommands:
   - `goapi route add --project "..." --method GET --path "/users" [--description "..."]`
     - Validates project exists before adding
     - Normalizes HTTP method (case-insensitive input)
     - Full error handling for invalid methods
   - `goapi route list --project "..."`
     - Displays routes in formatted table (ID, Method, Path)
     - Handles empty list case gracefully
     - Error handling for missing projects
   - `goapi route update --project "..." --id <route-id> [--method/--path/--description]`
     - All fields optional, only updates provided fields
     - Uses UpdateRouteInput with pointer-based nil checking
     - Validates project exists before updating
   - `goapi route delete --project "..." --id <route-id>`
     - Validates project exists before deleting
     - Proper error handling for missing routes

2. ✅ Implemented helper function `ParseHTTPMethod()` in `internal/route/route.go`:
   - Extracted switch logic into reusable function
   - Converts and validates HTTP method strings
   - Eliminates code duplication across commands

3. ✅ Tested all route commands end-to-end:
   - Create, list, update, delete operations fully functional
   - Error handling verified (invalid methods, missing projects, non-existent routes)
   - Flag validation working (required vs optional)
   - Table formatting with `text/tabwriter` works cleanly

**Session 5:**
1. ✅ Enhanced Route model with human-readable names:
   - **Data Model:** Added `Name` field to Route struct with composite unique index `(ProjectID, Name)`
   - **Database:** GORM composite unique constraint prevents duplicate names per project
   - **Storage Layer:** Added `GetRouteByName(projectID uint, name string)` helper function
   - **Name Validation:** CreateRoute validates names are not empty before persisting

2. ✅ Updated all route CRUD operations:
   - **CreateRoute()** — Now validates Name is not empty
   - **GetRouteByName()** — NEW helper function for name-based lookups (supports name-based CLI)
   - **UpdateRoute()** — Supports renaming via UpdateRouteInput.Name field
   - **ListRoutesByProject()** — No changes (already returns all routes for project)
   - **DeleteRoute()** — No changes (still uses ID internally)

3. ✅ Refactored route CLI commands to use names:
   - **route add** — Added `--name` flag (optional; auto-generates "METHOD path" if not provided)
   - **route list** — Displays ID, Name, Method, Path in aligned table with tabwriter
   - **route update** — Changed from `--id` to `--route` (name-based); added `--rename` for renaming
   - **route delete** — Changed from `--id` to `--route` (name-based identifier)
   - **Pattern:** CLI uses names, storage layer uses IDs internally (decouples UX from implementation)

4. ✅ Improved error handling consistency:
   - Wrapped all error returns with user-friendly context using `fmt.Errorf(..., %w, err)`
   - **Consistent patterns:**
     - `ParseHTTPMethod` errors: `"invalid HTTP method 'X': <error>"`
     - `GetRouteByName` errors: `"route 'X' not found in project 'Y': <error>"`
     - Storage operation errors: `"failed to <operation> route 'X': <error>"`
   - Users see context about what failed, what input caused it, and the underlying error

5. ✅ Improved project listing UI:
   - **Before:** Simple line format `Name - BaseURL (DateCreated)`
   - **After:** Professional table with 4 columns (Name, Base URL, Description, Date Created)
   - **Formatting:** Used text/tabwriter for clean alignment; dates formatted as `YYYY-MM-DD HH:MM`
   - **Consistency:** Matches route list table styling for unified CLI experience

6. ✅ Project maintenance:
   - Added `.claude/memory/` and `.claude/settings.local.json` to `.gitignore`
   - Kept `.claude/settings.json` in version control (team configuration with plugins and permissions)

7. ✅ Tested all changes end-to-end:
   - Route CRUD with names (add, list, update/rename, delete) all working
   - Error messages tested (nonexistent routes, invalid methods, unique constraint violations)
   - Project listing table rendering correctly with variable-length data
   - All commands compile without errors

### Next Steps

**Session 6 priorities:**
1. **API Client Design** (`internal/api/client.go`):
   - Create HTTP client abstraction with interface for easy mocking
   - Support GET, POST, PUT, DELETE, PATCH methods
   - Handle request/response bodies, headers, and status codes
   - Implement timeout and connection error handling
   - Plan JWT authorization support for future auth command

2. **Test Command Implementation** (`cmd/test.go`):
   - Execute a single route or all routes in a project
   - Flags: `--project`, `--route` (optional, run all if not specified)
   - Display results: status code, response time, response body
   - Plan output formatting (terminal pretty-print, JSON file export)

3. **Output Formatting** (`internal/output/`):
   - Design result struct (route, method, status, duration, body)
   - Implement terminal formatter (aligned table or pretty JSON)
   - Implement JSON file exporter for batch test results
   - Consider result aggregation for project-level test runs

4. **Consider for later:**
   - JWT auth command (`cmd/auth.go`) — Store tokens securely
   - Route variables/templating (e.g., `{id}`, `{user_id}`) in paths
   - Interactive TUI mode with BubbleTea
   - Batch testing with result summaries

### Architectural Decisions Made

- **Database:** SQLite with GORM (decided over JSON for scalability and learning value)
- **Primary Key:** Numeric ID with unique constraint on name (allows future renaming, better for relationships)
- **Data Location:** `~/.config/goapi/` (standard XDG-like path)
- **ORM:** GORM (selected for learning and maintainability)
- **Route Identification:** CLI uses human-readable names, storage uses numeric IDs (clean separation of concerns)
- **Error Handling:** Always wrap errors with context at each layer using `%w` verb (idiomatic Go)
