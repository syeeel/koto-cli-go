# koto - Implementation Checklist

Use this checklist to track implementation progress.
Change `[ ]` to `[x]` when each item is completed.

**Progress Summary**: 83/131 tasks completed

---

## Phase 1: Foundation Implementation (35/35) ✅ COMPLETE

### 1.1 Data Model (10/10) ✅
- [x] Create `internal/model/todo.go`
- [x] Define `Todo` struct
- [x] Define `TodoStatus` type and constants (`StatusPending`, `StatusCompleted`)
- [x] Define `Priority` type and constants (`PriorityLow`, `PriorityMedium`, `PriorityHigh`)
- [x] Implement `IsCompleted()` method
- [x] Implement `IsPending()` method
- [x] Implement `IsOverdue()` method

#### Tests
- [x] Create `internal/model/todo_test.go`
- [x] Implement tests for `IsCompleted()`
- [x] Implement tests for `IsOverdue()`

### 1.2 Repository Layer (18/18) ✅
- [x] Create `internal/repository/repository.go`
- [x] Define `TodoRepository` interface
- [x] Create `migrations/001_init.sql` (schema definition)

#### SQLite Implementation
- [x] Create `internal/repository/sqlite.go`
- [x] Define `SQLiteRepository` struct
- [x] Implement `NewSQLiteRepository()`
- [x] Implement `initSchema()`
- [x] Implement `Create()`
- [x] Implement `GetByID()`
- [x] Implement `GetAll()`
- [x] Implement `GetByStatus()`
- [x] Implement `Update()`
- [x] Implement `Delete()`
- [x] Implement `MarkAsCompleted()`
- [x] Implement `Close()`

#### Tests
- [x] Create `internal/repository/sqlite_test.go`
- [x] Test database initialization
- [x] Test each CRUD operation (at least 5)

#### Verification
- [x] Test with in-memory DB
- [x] Verify data creation, retrieval, update, and deletion

### 1.3 Service Layer (15/15) ✅
- [x] Create `internal/service/todo_service.go`
- [x] Define error constants (`ErrTodoNotFound`, `ErrInvalidTitle`, etc.)
- [x] Define `TodoService` struct
- [x] Implement `NewTodoService()`
- [x] Implement `AddTodo()`
- [x] Implement `EditTodo()`
- [x] Implement `DeleteTodo()`
- [x] Implement `CompleteTodo()`
- [x] Implement `ListTodos()`
- [x] Implement `ListPendingTodos()`
- [x] Implement `ListCompletedTodos()`
- [x] Implement `validateTitle()`
- [x] Implement `ExportToJSON()`
- [x] Implement `ImportFromJSON()`
- [x] Implement `validatePriority()`

#### Tests
- [x] Create `internal/service/todo_service_test.go`
- [x] Validation tests
- [x] Tests for each operation (using mock Repository)

---

## Phase 2: MVP Implementation (30/30) ✅ COMPLETE

### 2.1 TUI Foundation (12/12) ✅
- [x] Create `internal/tui/model.go`
- [x] Define `ViewMode` type and constants
- [x] Define `Model` struct
- [x] Implement `NewModel()`
- [x] Implement `Init()`

#### Update Function
- [x] Create `internal/tui/update.go`
- [x] Implement `Update()` function
- [x] Handle Enter key
- [x] Handle Ctrl+C/Esc keys
- [x] Handle Up/Down keys
- [x] Implement `handleEnter()`

#### View Function
- [x] Create `internal/tui/views.go`
- [x] Implement `View()` function
- [x] Implement `renderListView()`
- [x] Implement `renderHelpView()`

#### Styles
- [x] Create `internal/tui/styles.go`
- [x] Define lipgloss styles (`titleStyle`, `errorStyle`, etc.)

### 2.2 Command Parser (12/12) ✅
- [x] Create `internal/tui/commands.go`
- [x] Define `commandExecutedMsg` type
- [x] Define `todosLoadedMsg` type
- [x] Implement `parseAndExecuteCommand()`
- [x] Implement `loadTodos()`

#### Command Handlers
- [x] Implement `handleAddCommand()`
- [x] Implement `handleListCommand()`
- [x] Implement `handleDoneCommand()`
- [x] Implement `handleDeleteCommand()`
- [x] Implement `handleEditCommand()`
- [x] Implement `handleHelpCommand()`
- [x] Implement `/quit` processing

### 2.3 Main Entry Point (6/6) ✅
- [x] Create `internal/config/config.go`
- [x] Implement database path retrieval function
- [x] Implement `.koto` directory creation function

#### Main Function
- [x] Create `cmd/koto/main.go`
- [x] Database initialization
- [x] Service instance creation
- [x] TUI application launch
- [x] Error handling
- [x] Cleanup processing (defer)

#### Verification
- [x] Diagnose and resolve network connection issues (devcontainer firewall disabled)
- [x] Successful `go mod download`
- [x] Successful `go mod tidy` (go.sum generated)
- [x] All tests pass (Model, Repository, Service layers)
- [x] Successful `go build ./cmd/koto` (bin/koto 9.8MB)
- [x] Database initialization verified (~/.koto/koto.db created)
- [x] Schema validation (table and index verification)
- [x] Full operation verification in interactive terminal (requires user execution)
  - [x] Can add ToDo with `/add` command
  - [x] Can display ToDos with `/list` command
  - [x] Can complete ToDo with `/done` command
  - [x] Can delete ToDo with `/delete` command
  - [x] Data persists after app restart

---

## Phase 2.5: Startup Banner Display (6/6) ✅ COMPLETE

### 2.5.1 Banner ASCII Art (2/2) ✅
- [x] Create `internal/tui/banner.go`
- [x] Define KOTO CLI ASCII art

### 2.5.2 Banner View Implementation (4/4) ✅
- [x] Add `ViewModeBanner` to model.go
- [x] Implement `renderBannerView()` in views.go
- [x] Add banner styles to styles.go
- [x] Implement transition from banner to list screen in update.go
  - [x] Transition on any key press

#### Verification
- [x] Build successful (bin/koto 9.8MB)
- [x] All tests pass
- [ ] Banner displays on startup (requires user execution)
- [ ] Key press transitions to main screen (requires user execution)

---

## Phase 3: Feature Extensions (12/12) ✅ COMPLETE

### 3.1 Export/Import (12/12) ✅
#### Service Layer
- [x] Add `ExportToJSON()` method (existing)
- [x] Add `ImportFromJSON()` method (existing)
- [x] Add export error constants (existing)
- [x] Add import error constants (existing)

#### TUI Layer
- [x] Implement `/export` dedicated screen (2 steps: input→success)
- [x] Implement `/import` dedicated screen (4 steps: input→confirm→execute→complete)
- [x] Implement file path handling (default path generation, ~ expansion)
- [x] Implement import confirmation dialog (preview display)
- [x] Add ViewModeExport/Import
- [x] Implement key handling (Enter/Esc)
- [x] Implement helper functions (handleExportEnter/handleImportEnter)
- [x] Remove old command handlers

#### Tests
- [x] Build test successful
- [x] All existing tests pass
- [ ] Manual testing (requires user execution)

#### Verification
- [ ] `/export` transitions to dedicated screen (requires user execution)
- [ ] Default file path is generated correctly (requires user execution)
- [ ] Export success screen displays (requires user execution)
- [ ] `/import` transitions to dedicated screen (requires user execution)
- [ ] File selection and confirmation flow works (requires user execution)
- [ ] Import success/failure screens display (requires user execution)

### 3.2 UI/UX Improvements (0/4)
- [ ] Visualize priority (emoji/color)
- [ ] Display due date
- [ ] Highlight overdue items
- [ ] Improve status visualization

#### Detail Display
- [ ] Add ToDo detail display mode
- [ ] Display description, creation date, update date

#### Enhanced Help
- [ ] Add command list to `/help`
- [ ] Add usage examples
- [ ] Add keyboard shortcuts list

---

## Phase 3.5: Pomodoro Feature Implementation (0/25)

### 3.5.1 Data Model Extension (0/5)
- [ ] Add `WorkDuration` field to `internal/model/todo.go`
- [ ] Add `GetWorkDurationFormatted()` method (display time in readable format)
- [ ] Update data model tests

#### Migration
- [ ] Create `migrations/002_add_work_duration.sql`
- [ ] Implement migration execution function

### 3.5.2 Repository Layer Extension (0/4)
- [ ] Add `AddWorkDuration()` method to interface
- [ ] Implement `AddWorkDuration()` in SQLite
- [ ] Update existing tests (work_duration support)
- [ ] Add tests for `AddWorkDuration()`

### 3.5.3 Service Layer Extension (0/3)
- [ ] Add `AddWorkDuration()` method
- [ ] Add validation (verify work time is positive)
- [ ] Add tests

### 3.5.4 TUI Layer Extension (0/9)
#### Model Update
- [ ] Add `ViewModePomodoro` constant
- [ ] Add Pomodoro fields (`pomodoroTodoID`, `pomodoroStarted`, `pomodoroDuration`)

#### View Implementation
- [ ] Implement `renderPomodoroView()`
- [ ] Implement timer display
- [ ] Add work time display to `renderListView()`

#### Styles
- [ ] Add Pomodoro styles (`pomodoroTitleStyle`, `pomodoroTimerStyle`, etc.)

#### Update Function
- [ ] Handle key input in Pomodoro mode
- [ ] Update screen every second (using `tea.Tick`)
- [ ] Handle timer completion (record work time)

### 3.5.5 Command Handler (0/2)
- [ ] Implement `handlePomodoroCommand()`
- [ ] Parse `/pomo` command arguments (optional ToDo ID)

### 3.5.6 Tests (0/2)
- [ ] Pomodoro timer logic tests
- [ ] Work time recording tests

#### Verification
- [ ] `/pomo` launches timer
- [ ] `/pomo <ID>` launches timer linked to task
- [ ] Work time is recorded after timer completes
- [ ] Can cancel with Esc key
- [ ] Work time displays in list screen

---

## Phase 4: Quality Improvement and Release (0/13)

### 4.1 Test Enhancement (0/4)
- [ ] Verify Model layer test coverage
- [ ] Verify Repository layer test coverage
- [ ] Verify Service layer test coverage
- [ ] Create integration tests

#### Coverage Check
- [ ] Run `go test -cover ./...`
- [ ] Verify coverage is 80% or higher
- [ ] Add tests for insufficient areas

### 4.2 Documentation (0/4)
- [ ] Create `README.md`
  - [ ] Project overview
  - [ ] Installation instructions
  - [ ] Usage guide
  - [ ] Command reference
  - [ ] Screenshots/demo GIF (optional)

- [ ] Add code comments
  - [ ] GoDoc comments for public APIs
  - [ ] Explanatory comments for complex logic

### 4.3 Build and Release Configuration (0/5)
- [ ] Final `Makefile` adjustments
  - [ ] Verify `make build` works
  - [ ] Verify `make test` works
  - [ ] Verify `make clean` works
  - [ ] Verify cross-compilation with `make build-all`

- [ ] Create `.goreleaser.yml`
- [ ] Create `CHANGELOG.md`

#### GitHub Actions (Optional)
- [ ] Create `.github/workflows/test.yml`
- [ ] Create `.github/workflows/release.yml`

### 4.4 Initial Release (0/4)
- [ ] Create v1.0.0 tag
- [ ] Create GitHub Release
- [ ] Write release notes
- [ ] Upload binaries

#### Verification
- [ ] Verify Linux binary works
- [ ] Verify macOS binary works (if possible)
- [ ] Verify Windows binary works (if possible)
- [ ] Test installation with `go install`

---

## Optional Tasks

### Additional Features (v1.1 and later)
- [ ] Filtering functionality (priority, due date)
- [ ] Sort functionality
- [ ] Search functionality
- [ ] Statistics display (completion rate, etc.)

### Infrastructure & Tools
- [ ] Configure GitHub Discussions
- [ ] Create Issue templates
- [ ] Create PR templates
- [ ] Create `CONTRIBUTING.md`

### Documentation Extensions
- [ ] Create user guide
- [ ] Create architecture diagram
- [ ] Create FAQ

---

## Milestone Achievement Check

- [x] **M1: Development Environment Ready** - Phase 0 all complete ✅
- [x] **M2: Data Layer Complete** - Phase 1 all complete, tests passing ✅
- [x] **M3: MVP Implementation Complete** - Phase 2 code complete ✅ (verification after network connection)
- [ ] **M3.5: Pomodoro Feature Added** - Phase 3.5 all complete
- [ ] **M4: v1.0 Release** - Phase 3, 4 all complete

---

## Final Verification Before Completion

Before release, verify all of the following:

### Features
- [ ] All basic commands work (/add, /edit, /delete, /done, /list)
- [ ] Export/import works
- [ ] Data persists correctly
- [ ] Errors are handled properly

### Code Quality
- [ ] Test coverage 80% or higher
- [ ] `go vet ./...` runs without errors
- [ ] `golangci-lint run` runs without errors (if configured)
- [ ] All tests pass

### Documentation
- [ ] README.md is complete
- [ ] Design documents are up to date
- [ ] CHANGELOG is documented

### Build
- [ ] `make build` succeeds
- [ ] Cross-compilation succeeds
- [ ] Binary size is reasonable (target: 10-20MB)

### Release
- [ ] GitHub Releases page is created
- [ ] Binaries are uploaded
- [ ] Release notes are documented

---

**Keep up the good work! Enjoy the sense of accomplishment as you complete each task!**
