# koto - ToDo Management CLI Tool Basic Design Document

## 1. Project Overview

**koto** is an interactive ToDo list management CLI tool developed in Go.
It provides a rich and intuitive terminal UI using the bubbletea framework.

## 2. Purpose and Goals

### Purpose
- Provide a ToDo list management tool that can be used seamlessly in a terminal environment
- Improve productivity through an intuitive and easy-to-use interface
- Lightweight and fast operation

### Goals
- Achieve the same level of operability as dedicated screens like Claude Code
- Intuitive operation through slash commands
- Data persistence to local storage
- Cross-platform support (Linux, macOS, Windows)

## 3. Main Features

### 3.1 Command List

| Command | Function | Description |
|---------|----------|-------------|
| `koto` | Launch app | Enter dedicated screen and start interactive mode |
| `/add` | Add ToDo | Add a new ToDo item |
| `/edit` | Edit ToDo | Edit an existing ToDo item |
| `/delete` | Delete ToDo | Delete a ToDo item |
| `/done` | Complete ToDo | Mark a ToDo item as completed |
| `/list` | Display list | Display all ToDo items |
| `/export` | Export | Export ToDo data to JSON file |
| `/import` | Import | Import ToDo data from JSON file |
| `/pomo` | Pomodoro Timer | Start a 25-minute timer (optionally specify task ID) |
| `/exit` | Exit | Exit the application |

### 3.2 Feature Details

#### ToDo Management Features
- Create, edit, delete, and complete ToDo items
- Manage ToDo item status (incomplete/complete)
- Set ToDo item priority (optional)
- Set ToDo item deadline (optional)

#### UI/UX Features
- Real-time command input and feedback
- Keyboard navigation (arrow keys, Enter, Esc, etc.)
- Colorful and highly visible display
- Responsive layout

#### Data Management Features
- Export to JSON file (ensures data portability)
- Import from JSON file (migration from other environments, restore from backup)
- Automatic database initialization (on first launch)

#### Pomodoro Timer Feature
- 25-minute countdown timer
- Alarm notification when timer ends
- Task-specific work time recording (automatic calculation of cumulative time)
- Display dedicated timer screen
- Timer cancellation feature

## 4. Technology Stack

### Programming Language
- **Go 1.21+**

### Main Libraries
- **bubbletea**: TUI framework
- **bubbles**: UI components for bubbletea
- **lipgloss**: Styling library
- **modernc.org/sqlite**: Data persistence (Pure Go implementation, no CGO required)

### Development Tools
- Go modules: Dependency management
- golangci-lint: Code quality checks
- go test: Testing framework

## 5. System Architecture Overview

```
┌─────────────────────────────────────────┐
│           koto CLI Application          │
├─────────────────────────────────────────┤
│  Presentation Layer (Bubbletea)        │
│  ┌────────────────────────────────┐    │
│  │  Main Model (TUI State)        │    │
│  │  ├─ Command Parser             │    │
│  │  ├─ View Renderer              │    │
│  │  └─ Event Handler              │    │
│  └────────────────────────────────┘    │
├─────────────────────────────────────────┤
│  Business Logic Layer                  │
│  ┌────────────────────────────────┐    │
│  │  Todo Service                  │    │
│  │  ├─ Add                        │    │
│  │  ├─ Edit                       │    │
│  │  ├─ Delete                     │    │
│  │  ├─ Done                       │    │
│  │  └─ List                       │    │
│  └────────────────────────────────┘    │
├─────────────────────────────────────────┤
│  Data Access Layer                     │
│  ┌────────────────────────────────┐    │
│  │  Repository Interface          │    │
│  │  └─ SQLite Implementation     │    │
│  └────────────────────────────────┘    │
├─────────────────────────────────────────┤
│  Data Storage                          │
│  ┌────────────────────────────────┐    │
│  │  SQLite Database               │    │
│  │  (~/.koto/koto.db)             │    │
│  └────────────────────────────────┘    │
└─────────────────────────────────────────┘
```

## 6. Data Model Overview

### ToDo Item
- **ID**: Unique identifier (auto-generated)
- **Title**: Title (required)
- **Description**: Description (optional)
- **Status**: Status (incomplete/complete)
- **Priority**: Priority (Low/Medium/High) (optional)
- **DueDate**: Deadline date/time (optional)
- **WorkDuration**: Cumulative work time (minutes) (automatically recorded by Pomodoro timer)
- **CreatedAt**: Creation date/time (auto-generated)
- **UpdatedAt**: Update date/time (auto-updated)

## 7. UI Flow Overview

```
[App Launch]
    ↓
[Main Screen]
    ├─ ToDo list display
    │   └─ Focus on task with ↑/↓
    │       └─ Empty input + Enter → [Task Detail Screen] → [Main Screen]
    └─ Command input area
        ↓
[Command Input]
    ├─ /add → [Add ToDo Screen] → [Main Screen]
    ├─ /edit → [Select ToDo] → [Edit Screen] → [Main Screen]
    ├─ /delete → [Select ToDo] → [Delete Confirmation] → [Main Screen]
    ├─ /done → [Select ToDo] → [Complete Process] → [Main Screen]
    ├─ /list → [Update List Display] → [Main Screen]
    ├─ /pomo → [Pomodoro Timer Screen] → [Complete/Cancel] → [Main Screen]
    └─ /quit → [Exit App]

[Task Detail Screen]
    ├─ Display all task information (title, description, status, priority, work time, deadline, creation date, update date)
    ├─ Esc → [Main Screen]
    ├─ e → [Edit Screen] → [Main Screen]
    └─ d → Complete task → [Main Screen]
```

## 8. Non-Functional Requirements

### Performance
- App startup time: Within 1 second
- Command execution time: Within 100ms
- Smooth operation even with 10,000 ToDo items

### Security
- Database file is stored in the user's home directory
- File permissions: 600 (read/write only by owner)

### Maintainability
- Modular design
- Clear separation of concerns
- Comprehensive test coverage (80% or higher)
- Detailed documentation

### Usability
- Intuitive command system
- Rich help messages
- Clear error messages
- Keyboard shortcuts provided

## 9. Distribution Method

### 9.1 Single Binary Distribution
By adopting modernc.org/sqlite (Pure Go implementation), we have the following advantages:
- **No dependencies**: SQLite library is included in the binary
- **Single file**: Users only need to download one executable file
- **Cross-compilation support**: Easy to build for all platforms: Windows, macOS, Linux

### 9.2 Distribution Channels
- **GitHub Releases**: Provide binaries for each platform
  - `koto-linux-amd64`
  - `koto-darwin-amd64` / `koto-darwin-arm64`
  - `koto-windows-amd64.exe`
- **Go install**: `go install github.com/yourusername/koto@latest`
- **Homebrew**: Planned for future support

### 9.3 Build Method
```bash
# Pure Go, so cross-compilation is easy without CGO
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o koto-linux-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o koto-darwin-arm64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o koto-windows-amd64.exe
```

### 9.4 User Experience
```bash
# Download and run immediately (no additional installation required)
wget https://github.com/user/koto/releases/download/v1.0.0/koto-linux-amd64
chmod +x koto-linux-amd64
./koto-linux-amd64
```

## 10. Advantages of Adopting SQLite

### 10.1 Data Persistence
- **Reliability**: Data consistency through transaction support
- **Performance**: Fast search through indexing
- **Scalability**: No performance degradation even with large amounts of data (tens of thousands of records)

### 10.2 Ensuring Data Portability
SQLite is in binary format, but data portability is ensured through the following methods:
- Export to JSON format via `/export` command
- Import from JSON format via `/import` command
- Easy backup and migration

### 10.3 Easy Future Extensibility
SQLite's powerful query capabilities make it easy to implement the following features:
- Complex search and filtering (full-text search of title, description, tags)
- Statistics and analysis (completion rate, deadline compliance rate, priority-based aggregation)
- Advanced sorting (multiple conditions, custom sort order)

## 11. Future Extensibility

- Tag feature
- Advanced filtering and search features (full-text search)
- Category classification
- Recurring tasks
- CSV export
- Cloud sync feature
- Team sharing feature
- Pomodoro timer extensions
  - Customizable timer duration
  - Break time timers (5 min/15 min)
  - Record number of Pomodoros
  - Work time statistics and graph display
  - Work time goal setting
