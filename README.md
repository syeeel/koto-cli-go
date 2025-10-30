# koto - ToDo Management CLI

**koto** (meaning "thing" or "matter" in Japanese) is an interactive ToDo management CLI tool developed in Go.
It provides a comfortable task management experience with a beautiful terminal UI using the [Bubbletea](https://github.com/charmbracelet/bubbletea) framework.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)

## ✨ Features

- 🎨 **Rich TUI** - Beautiful terminal interface with Bubbletea/Lipgloss
- ⚡ **Lightweight & Fast** - Pure Go (no CGO required) with fast startup
- 📊 **Priority Management** - 3-level priority system (🔴High 🟡Medium 🟢Low)
- 📅 **Due Date Management** - Set due dates with overdue warnings
- 💾 **SQLite Storage** - Reliable local database for data persistence
- 📤 **Export/Import** - Backup and migration in JSON format
- ⌨️ **Vim-like Keybindings** - Comfortable navigation with j/k
- 🔍 **Status Filtering** - Filter by pending/completed status
- 🍅 **Pomodoro Timer** - 25-minute timer to support focused work with automatic time tracking

## 📦 Installation

### Installation Script (Recommended for macOS/Linux)

The easiest method:

```bash
curl -sSfL https://raw.githubusercontent.com/syeeel/koto-cli-go/main/install.sh | sh
```

This script will:
- Auto-detect the latest version
- Download for your OS/architecture
- Install to `~/.local/bin`
- Guide you through PATH setup

### Homebrew (macOS - Coming Soon 🚧)

**Currently under setup. Will be available from the next release (v1.0.1 or later).**

Once ready, you'll be able to install with:

```bash
brew tap syeeel/tap
brew install koto
```

**Until then, please use the installation script (above).**

### Go install

If you have a Go environment:

```bash
go install github.com/syeeel/koto-cli-go/cmd/koto@latest
```

### Pre-built Binaries

You can download binaries for your platform from the [Releases](https://github.com/syeeel/koto-cli-go/releases/latest) page.

Supported platforms:
- **macOS**: darwin_amd64 (Intel), darwin_arm64 (Apple Silicon)
- **Linux**: linux_amd64, linux_arm64
- **Windows**: windows_amd64

After downloading, extract and place the binary in a directory in your PATH.

### Build from Source

```bash
# Clone the repository
git clone https://github.com/syeeel/koto-cli-go.git
cd koto-cli-go

# Download dependencies
go mod download

# Build
go build -o bin/koto ./cmd/koto

# Run
./bin/koto
```

## 🚀 Usage

### Starting the Application

```bash
koto
```

Once started, an interactive TUI will be displayed.

### Basic Commands

#### Adding a ToDo

```bash
/add Go shopping
/add Write report --desc="Summarize Chapter 5" --priority=high --due=2025-10-25
```

**Options**:
- `--desc="description"` - Detailed description of the ToDo
- `--priority=low|medium|high` - Priority level (default: medium)
- `--due=YYYY-MM-DD` - Due date

#### Listing ToDos

```bash
/list                      # Show all ToDos
/list --status=pending     # Pending only
/list --status=completed   # Completed only
```

#### Completing a ToDo

```bash
/done 1    # Mark ToDo with ID 1 as completed
```

#### Editing a ToDo

```bash
/edit 1 --title="New title"
/edit 1 --priority=high
/edit 1 --desc="New description"
/edit 1 --due=2025-12-31
```

#### Deleting a ToDo

```bash
/delete 1    # Delete ToDo with ID 1
```

#### Export/Import

```bash
/export ~/my-todos.json     # Export to JSON file
/import ~/todos-backup.json # Import from JSON file
```

#### Help

```bash
/help    # Show help screen
```

#### Pomodoro Timer

```bash
/pomo              # Start a 25-minute timer
/pomo 1            # Start a 25-minute timer linked to ToDo ID 1 (automatically records work time)
```

**How to use the Pomodoro Timer**:
- A dedicated screen is displayed during the timer
- An alarm sounds after 25 minutes
- If a task ID is specified, work time is automatically recorded
- Press `Esc` to cancel the timer and return to the main screen

### ⌨️ Keyboard Shortcuts

| Key | Action |
|------|------|
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `Enter` | Execute command |
| `Esc` | Clear input field |
| `?` | Show/hide help screen |
| `Ctrl+C` | Exit application |

### 📺 Screen Layout

```
📝 koto - ToDo Manager

  ⬜ 🔴 [1] Prepare for important meeting
> ✅ 🟡 [2] Shopping list ⚠ OVERDUE
  ⬜ 🟢 [3] Reply to emails

> /add New task

Commands: /add, /list, /done, /delete, /edit, /help | Navigate: ↑/↓ or j/k | Help: ? | Quit: Ctrl+C
```

**Display Explanation**:
- `>` - Currently selected ToDo (cursor position)
- `⬜` - Pending
- `✅` - Completed
- `🔴🟡🟢` - Priority (High, Medium, Low)
- `[number]` - ToDo ID
- `⚠ OVERDUE` - Overdue warning
- `🍅 XXXm` - Cumulative work time (recorded by Pomodoro timer)

## 📁 Data Storage Location

All ToDos are stored in the following SQLite database:

```
~/.koto/koto.db
```

To back up, copy this file or use the `/export` command.

## 🏗️ Architecture

koto adopts a layered architecture based on clean architecture principles:

```
┌─────────────────┐
│   TUI Layer     │  Bubbletea UI (command input, display)
├─────────────────┤
│ Service Layer   │  Business logic, validation
├─────────────────┤
│Repository Layer │  Data access (SQLite)
├─────────────────┤
│  Model Layer    │  Data structure definitions
└─────────────────┘
```

### Directory Structure

```
koto-cli-go/
├── cmd/
│   └── koto/              # Main entry point
│       └── main.go
├── internal/
│   ├── model/             # Data models (Todo, Status, Priority)
│   │   ├── todo.go
│   │   └── todo_test.go
│   ├── repository/        # Data access layer
│   │   ├── repository.go
│   │   ├── sqlite.go
│   │   └── sqlite_test.go
│   ├── service/           # Business logic layer
│   │   ├── todo_service.go
│   │   └── todo_service_test.go
│   ├── tui/               # Terminal UI
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── views.go
│   │   ├── styles.go
│   │   └── commands.go
│   └── config/            # Configuration management
│       └── config.go
├── migrations/            # Database schema
│   └── 001_init.sql
├── docs/                  # Documentation
│   ├── design/            # Design documents
│   └── implementation/    # Implementation management
├── go.mod
├── go.sum
└── README.md
```

## 🛠️ Development Environment

### Requirements

- Go 1.21 or later
- SQLite 3 (not required as we use Pure Go implementation)

### Dependencies

- [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [github.com/charmbracelet/bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) - Pure Go SQLite

### Development Commands

```bash
# Download dependencies
go mod download

# Run tests
go test ./...
go test -v ./internal/model/...     # Model layer only
go test -v ./internal/repository/... # Repository layer only
go test -v ./internal/service/...    # Service layer only

# Lint
go vet ./...
golangci-lint run  # If golangci-lint is installed

# Build
go build -o bin/koto ./cmd/koto

# Cross-compile
GOOS=darwin GOARCH=amd64 go build -o bin/koto-darwin-amd64 ./cmd/koto
GOOS=linux GOARCH=amd64 go build -o bin/koto-linux-amd64 ./cmd/koto
GOOS=windows GOARCH=amd64 go build -o bin/koto-windows-amd64.exe ./cmd/koto
```

### DevContainer Environment

This project includes a DevContainer environment for VS Code / Cursor.

```bash
# Open in VS Code / Cursor
# Simply select "Reopen in Container" to set up your development environment
```

## 🧪 Testing

This project implements comprehensive tests for each layer.

```bash
# Run all tests
go test ./...

# Coverage report
go test -cover ./...

# Detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Test Statistics**:
- Model layer: 3 test functions, 7 subtests
- Repository layer: 9 test functions (using in-memory DB)
- Service layer: 13 test functions (using mock Repository)

## 📝 License

This project is released under the [MIT License](LICENSE.md).

## 🤝 Contributing

Pull requests are welcome! For bug reports and feature requests, please use [Issues](https://github.com/syeeel/koto-cli-go/issues).

### Development Guidelines

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Add and run tests (`go test ./...`)
4. Commit (`git commit -m 'feat: Add amazing feature'`)
5. Push (`git push origin feature/amazing-feature`)
6. Create a pull request

For more details, see the development guide in [.claude/CLAUDE.md](.claude/CLAUDE.md).

## 🔗 References

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) - Pure Go SQLite


## 💡 FAQ

### Q: Where is the data stored?

A: Data is stored as an SQLite database at `~/.koto/koto.db`.

### Q: Can I sync ToDos across multiple machines?

A: Currently, there is no sync feature. However, you can export to a JSON file using `/export` and import it on another machine with `/import`.

### Q: Is there a search feature for ToDos?

A: The search feature is not implemented in the current version. It is planned for a future release.

### Q: Does it work on Windows?

A: Yes, thanks to the Pure Go implementation, it works on Windows, macOS, and Linux.

## 📮 Contact

For bug reports and questions, please use [GitHub Issues](https://github.com/syeeel/koto-cli-go/issues).

---

**koto** - Manage your tasks comfortably! 🎵
